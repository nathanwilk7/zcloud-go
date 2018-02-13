package zcloud

import (
	"io"
	"time"
)

func newGCloudObject (bucket, key string, b *gCloudBucket) gCloudObject {
	return gCloudObject{
		b: b,
		bucket: bucket,
		key: key,
	}
}

type gCloudObject struct {
	b *gCloudBucket
	bucket, key string
	lastModified time.Time
	size int
	reader gCloudObjectReader
	writer gCloudObjectWriter
}
func (o gCloudObject) Delete () error {
	err := o.b.p.client.Bucket(o.b.Name()).Object(o.Key()).Delete(o.b.p.context)
	return err
}

func (o gCloudObject) Key () string {
	return o.key
}

func (o gCloudObject) LastModified () (time.Time, error) {
	if err := o.get(); err != nil {
		return time.Time{}, err
	}
	return o.lastModified, nil
}

func (o gCloudObject) Size () (int, error) {
	if err := o.get(); err != nil {
		return 0, err
	}
	return o.size, nil
}

type gCloudObjectReader struct {
	rc io.ReadCloser
}

func (o gCloudObject) Reader () (io.ReadCloser, error) {
	r, err := o.b.p.client.Bucket(o.b.Name()).Object(o.Key()).NewReader(o.b.p.context)
	if err != nil {
		return nil, err
	}
	o.reader.rc = r
	return o.reader, nil
}

func (w gCloudObjectReader) Read (b []byte) (int, error) {
	return w.rc.Read(b)
}

func (w gCloudObjectReader) Close () error {
	return w.rc.Close()
}

type gCloudObjectWriter struct {
	wc io.WriteCloser
}

func (o gCloudObject) Writer () (io.WriteCloser, error) {
	w := o.b.p.client.Bucket(o.b.Name()).Object(o.Key()).NewWriter(o.b.p.context)
	o.writer.wc = w
	return o.writer, nil
}

func (w gCloudObjectWriter) Write (b []byte) (int, error) {
	return w.wc.Write(b)
}

func (w gCloudObjectWriter) Close () error {
	return w.wc.Close()
}

func (o *gCloudObject) get () error {
	objAttrs, err := o.b.p.client.Bucket(o.b.Name()).Object(o.Key()).Attrs(o.b.p.context)
	if err != nil {
		return err
	}
	o.size = int(objAttrs.Size)
	o.lastModified = objAttrs.Updated
	return nil
}
