package zcloud

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
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
	writer *awsObjectWriter
	b *awsBucket
}

func (src awsObject) CopyTo (dest Object) error {
	// Added this type check to make it easy to do fast copying
	d, ok := dest.(awsObject)
	if !ok {
		return fmt.Errorf("AWS CopyTo currently only works for objects of the same provider. src: %v, dest: %v", src, dest)
	}
	coi := &s3.CopyObjectInput{
		Bucket: aws.String(d.b.Name()),
		CopySource: aws.String(fmt.Sprintf("%s/%s", src.b.Name(), src.Key())),
		Key: aws.String(d.Key()),
	}
	s3svc := s3.New(src.b.p.session)
	_, err := s3svc.CopyObject(coi)
	return err
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

type awsObjectReader struct {
	rc io.ReadCloser
}

func (o awsObject) Reader () (io.ReadCloser, error) {
	s3svc := s3.New(o.b.p.session)
	input := &s3.GetObjectInput{
		Bucket: aws.String(o.bucket),
		Key: aws.String(o.Key()),
	}
	object, err := s3svc.GetObject(input)
	if err != nil {
		return nil, err
	}
	return object.Body, nil
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

func (o awsObject) Info () (ObjectInfo, error) {
	aoi := awsObjectInfo{}
	s3svc := s3.New(o.b.p.session)
	input := &s3.GetObjectInput{
		Bucket: aws.String(o.bucket),
		Key: aws.String(o.Key()),
	}
	object, err := s3svc.GetObject(input)
	if err != nil {
		aerr, ok := err.(awserr.Error)
		if ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchKey:
				err = ErrObjectDoesNotExist{
					bucket: o.b.Name(),
					key: o.Key(),
				}
			}
		}
		return aoi, err
	}
	aoi.lastModified = *object.LastModified
	aoi.size = int(aws.Int64Value(object.ContentLength))
	return aoi, nil
}

func (i awsObjectInfo) LastModified () time.Time {
	return i.lastModified
}

func (i awsObjectInfo) Size () int {
	return i.size
}

type awsObjectInfo struct {
	lastModified time.Time
	size int
}

type awsObjectCopier struct {
	coi *s3.CopyObjectInput
}
