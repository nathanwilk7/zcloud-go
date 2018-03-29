package zcloud

import (
	"fmt"
	"io"
	"time"

	gs "cloud.google.com/go/storage"
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
	gco *gs.ObjectHandle
}

func (src gCloudObject) CopyTo (dest Object) error {
	src.gco = src.getGCloudObject()
	// Added this type check to make it easy to do fast copying
	d, ok := dest.(gCloudObject)
	if !ok {
		return fmt.Errorf("gcloud CopyTo currently only works for objects of the same provider. src: %v, dest: %v", src, dest)
	}
	d.gco = d.getGCloudObject()
	_, err := d.gco.CopierFrom(src.gco).Run(src.b.p.context)
	return err
}

func (o gCloudObject) Delete () error {
	err := o.getGCloudObject().Delete(o.b.p.context)
	return err
}

func (o gCloudObject) Key () string {
	return o.key
}

type gCloudObjectReader struct {
	rc io.ReadCloser
}

func (o gCloudObject) Reader () (io.ReadCloser, error) {
	r, err := o.getGCloudObject().NewReader(o.b.p.context)
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
	w := o.getGCloudObject().NewWriter(o.b.p.context)
	o.writer.wc = w
	return o.writer, nil
}

func (w gCloudObjectWriter) Write (b []byte) (int, error) {
	return w.wc.Write(b)
}

func (w gCloudObjectWriter) Close () error {
	return w.wc.Close()
}

func (o *gCloudObject) getGCloudObject () *gs.ObjectHandle {
	if o.gco == nil {
		o.gco = o.b.getGCloudBucket().Object(o.Key())
	}
	return o.gco
}

func (o gCloudObject) Info () (ObjectInfo, error) {
	gcoi := gCloudObjectInfo{}
	objAttrs, err := o.getGCloudObject().Attrs(o.b.p.context)
	if err != nil {
		switch err {
		case gs.ErrObjectNotExist:
			err = ErrObjectDoesNotExist{
				bucket: o.b.Name(),
				key: o.Key(),
			}
		}
		return gcoi, err
	}
	gcoi.size = int(objAttrs.Size)
	gcoi.lastModified = objAttrs.Updated
	return gcoi, nil
}

type gCloudObjectInfo struct {
	lastModified time.Time
	size int
}

func (i gCloudObjectInfo) LastModified () time.Time {
	return i.lastModified
}

func (i gCloudObjectInfo) Size () int {
	return i.size
}

type GCloudObjectDoesNotExist struct {
	bucket string
	key string
}
