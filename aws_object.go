package zcloud

import (
	"bytes"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func NewAwsObject (bucket, key string, b *AwsBucket) AwsObject {
	return AwsObject{
		bucket: bucket,
		key: key,
		b: b,
	}
}

type AwsObject struct {
	bucket, key string
	lastModified time.Time
	size int
	reader AwsObjectReader
	writer *AwsObjectWriter
	b *AwsBucket
}

func (o AwsObject) Delete () error {
	s3svc := s3.New(o.b.p.session)
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(o.bucket),
		Key: aws.String(o.Key()),
	}
	_, err := s3svc.DeleteObject(input)
	return err
}

func (o AwsObject) Key () string {
	return o.key
}

func (o AwsObject) LastModified () (time.Time, error) {
	if err := o.get(); err != nil {
		return time.Time{}, err
	}
	return o.lastModified, nil
}

func (o AwsObject) Size () (int, error) {
	if err := o.get(); err != nil {
		return 0, err
	}
	return o.size, nil
}

func NewAwsObjectReader () AwsObjectReader {
	return AwsObjectReader{}
}

type AwsObjectReader struct {
	rc io.ReadCloser
}

func (o AwsObject) Reader () (io.ReadCloser, error) {
	if err := o.get(); err != nil {
		return nil, err
	}
	return o.reader, nil
}

func (w AwsObjectReader) Read (b []byte) (int, error) {
	return w.rc.Read(b)
}

func (w AwsObjectReader) Close () error {
	return w.rc.Close()
}

type AwsObjectWriter struct {
	o *AwsObject
	ui s3manager.UploadInput
	b []byte
}

func (o AwsObject) Writer () (io.WriteCloser, error) {
	w := AwsObjectWriter{}
	o.writer = &w
	o.writer.ui = s3manager.UploadInput{
		Bucket: aws.String(o.bucket),
		Key: aws.String(o.Key()),
	}
	o.writer.o = &o
	return o.writer, nil
}

func (w *AwsObjectWriter) Write (b []byte) (int, error) {
	w.b = append(w.b, b...)
	return len(b), nil
}

func (w *AwsObjectWriter) Close () error {
	w.ui.Body = bytes.NewReader(w.b)
	s3svc := s3.New(w.o.b.p.session)
	uploader := s3manager.NewUploaderWithClient(s3svc)
	_, err := uploader.Upload(&w.ui)
	if err != nil {
		return err
	}
	return nil
}

func (o *AwsObject) get () error {
	s3svc := s3.New(o.b.p.session)
	input := &s3.GetObjectInput{
		Bucket: aws.String(o.bucket),
		Key: aws.String(o.Key()),
	}
	object, err := s3svc.GetObject(input)
	if err != nil {
		return err
	}
	o.lastModified = *object.LastModified
	o.size = int(aws.Int64Value(object.ContentLength))
	o.reader.rc = object.Body
	return nil
}
