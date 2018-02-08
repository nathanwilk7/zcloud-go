package zcloud

import (
	"io"
	"time"
)

func NewGCloudObject (bucket, key string, b *GCloudBucket) GCloudObject {
	return GCloudObject{
		b: b,
		bucket: bucket,
		key: key,
	}
}

type GCloudObject struct {
	b *GCloudBucket
	bucket, key string
	lastModified time.Time
	size int
	reader GCloudObjectReader
	writer GCloudObjectWriter
}

func (o GCloudObject) Key () string {
	return o.key
}

func (o GCloudObject) LastModified () (time.Time, error) {
	if err := o.get(); err != nil {
		return time.Time{}, err
	}
	return o.lastModified, nil
}

func (o GCloudObject) Size () (int, error) {
	if err := o.get(); err != nil {
		return 0, err
	}
	return o.size, nil
}

func NewGCloudObjectReader () GCloudObjectReader {
	return GCloudObjectReader{}
}

type GCloudObjectReader struct {
	rc io.ReadCloser
}

func (o GCloudObject) Reader () (io.ReadCloser, error) {
	r, err := o.b.p.client.Bucket(o.b.Name()).Object(o.Key()).NewReader(o.b.p.context)
	if err != nil {
		return nil, err
	}
	o.reader.rc = r
	return o.reader, nil
}

func (w GCloudObjectReader) Read (b []byte) (int, error) {
	return w.rc.Read(b)
}

func (w GCloudObjectReader) Close () error {
	return w.rc.Close()
}

type GCloudObjectWriter struct {
	wc io.WriteCloser
}

func (o GCloudObject) Writer () (io.WriteCloser, error) {
	w := o.b.p.client.Bucket(o.b.Name()).Object(o.Key()).NewWriter(o.b.p.context)
	o.writer.wc = w
	return o.writer, nil
}

func (w GCloudObjectWriter) Write (b []byte) (int, error) {
	return w.wc.Write(b)
}

func (w GCloudObjectWriter) Close () error {
	return w.wc.Close()
}

func (o GCloudObject) get () error {
	objAttrs, err := o.b.p.client.Bucket(o.b.Name()).Object(o.Key()).Attrs(o.b.p.context)
	if err != nil {
		return err
	}
	o.size = int(objAttrs.Size)
	o.lastModified = objAttrs.Updated
	return nil
}
