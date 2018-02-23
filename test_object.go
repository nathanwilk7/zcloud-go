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
	ReaderField TestObjectReader
	WriterField *TestObjectWriter
	B *TestBucket
}

func (src TestObject) CopyTo (dest Object) error {
	return nil
}

func (o TestObject) Delete () error {
	return nil
}

func (o TestObject) Key () string {
	return o.KeyField
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

func (o TestObject) Info () (ObjectInfo, error) {
	return TestObjectInfo{}, nil
}

type TestObjectInfo struct {
	LastModifiedField time.Time
	SizeField int
}

func (i TestObjectInfo) LastModified () time.Time {
	return i.LastModifiedField
}

func (i TestObjectInfo) Size () int {
	return i.SizeField
}
