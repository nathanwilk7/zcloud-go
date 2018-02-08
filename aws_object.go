package zcloud

import (
	"bytes"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func NewAwsObject (bucket, key string) AwsObject {
	return AwsObject{
		bucket: bucket,
		key: key,
	}
}

type AwsObject struct {
	bucket, key string
	lastModified time.Time
	size int
	reader AwsObjectReader
	writer AwsObjectWriter
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
	ui s3manager.UploadInput
	b []byte
}

func (o AwsObject) Writer () (io.WriteCloser, error) {
	o.writer.ui = s3manager.UploadInput{
		Bucket: aws.String(o.bucket),
		Key: aws.String(o.Key()),
	}
	return o.writer, nil
}

func (w AwsObjectWriter) Write (b []byte) (int, error) {
	w.b = append(w.b, b...)
	return len(b), nil
}

func (w AwsObjectWriter) Close () error {
	w.ui.Body = bytes.NewReader(w.b)
	s3svc := s3.New(session.New())
	uploader := s3manager.NewUploaderWithClient(s3svc)
	_, err := uploader.Upload(&w.ui)
	if err != nil {
		return err
	}
	return nil
}

func (o AwsObject) get () error {
	s3svc := s3.New(session.New())
	input := &s3.GetObjectInput{
		Bucket: aws.String(o.bucket),
		Key: aws.String(o.Key()),
	}
	object, err := s3svc.GetObject(input)
	if err != nil {
		return err
	}
	o.lastModified = *object.LastModified
	o.size = int(*object.ContentLength)
	o.reader.rc = object.Body
	return nil
}
