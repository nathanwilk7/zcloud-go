package zcloud

import (
	"io"
	"time"
)

func NewTestObject (bucket, key string, b *TestBucket) TestObject {
	return TestObject{
		BucketField: bucket,
		KeyField: key,
		B: b,
	}
}

type TestObject struct {
	BucketField, KeyField string
	LastModifiedField time.Time
	SizeField int
	ReaderField TestObjectReader
	WriterField *TestObjectWriter
	B *TestBucket
}

func (o TestObject) Delete () error {
	return nil
}

func (o TestObject) Key () string {
	return o.KeyField
}

func (o TestObject) LastModified () (time.Time, error) {
	return time.Time{}, nil
}

func (o TestObject) Size () (int, error) {
	return o.SizeField, nil
}

type TestObjectReader struct {}

func (o TestObject) Reader () (io.ReadCloser, error) {
	return nil, nil
}

func (w TestObjectReader) Read (b []byte) (int, error) {
	return 0, nil
}

func (w TestObjectReader) Close () error {
	return nil
}

type TestObjectWriter struct {
	O *TestObject
	B []byte
}

func (o TestObject) Writer () (io.WriteCloser, error) {
	return nil, nil
}

func (w *TestObjectWriter) Write (b []byte) (int, error) {
	return 0, nil
}

func (w *TestObjectWriter) Close () error {
	return nil
}
