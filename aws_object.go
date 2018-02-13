package zcloud

import (
	"bytes"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func newAwsObject (bucket, key string, b *awsBucket) awsObject {
	return awsObject{
		bucket: bucket,
		key: key,
		b: b,
	}
}

type awsObject struct {
	bucket, key string
	lastModified time.Time
	size int
	reader awsObjectReader
	writer *awsObjectWriter
	b *awsBucket
}

func (o awsObject) Delete () error {
	s3svc := s3.New(o.b.p.session)
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(o.bucket),
		Key: aws.String(o.Key()),
	}
	_, err := s3svc.DeleteObject(input)
	return err
}

func (o awsObject) Key () string {
	return o.key
}

func (o awsObject) LastModified () (time.Time, error) {
	if err := o.get(); err != nil {
		return time.Time{}, err
	}
	return o.lastModified, nil
}

func (o awsObject) Size () (int, error) {
	if err := o.get(); err != nil {
		return 0, err
	}
	return o.size, nil
}

func newAwsObjectReader () awsObjectReader {
	return awsObjectReader{}
}

type awsObjectReader struct {
	rc io.ReadCloser
}

func (o awsObject) Reader () (io.ReadCloser, error) {
	if err := o.get(); err != nil {
		return nil, err
	}
	return o.reader, nil
}

func (w awsObjectReader) Read (b []byte) (int, error) {
	return w.rc.Read(b)
}

func (w awsObjectReader) Close () error {
	return w.rc.Close()
}

type awsObjectWriter struct {
	o *awsObject
	ui s3manager.UploadInput
	b []byte
}

func (o awsObject) Writer () (io.WriteCloser, error) {
	w := awsObjectWriter{}
	o.writer = &w
	o.writer.ui = s3manager.UploadInput{
		Bucket: aws.String(o.bucket),
		Key: aws.String(o.Key()),
	}
	o.writer.o = &o
	return o.writer, nil
}

func (w *awsObjectWriter) Write (b []byte) (int, error) {
	w.b = append(w.b, b...)
	return len(b), nil
}

func (w *awsObjectWriter) Close () error {
	w.ui.Body = bytes.NewReader(w.b)
	s3svc := s3.New(w.o.b.p.session)
	uploader := s3manager.NewUploaderWithClient(s3svc)
	_, err := uploader.Upload(&w.ui)
	if err != nil {
		return err
	}
	return nil
}

func (o *awsObject) get () error {
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
