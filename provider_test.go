package zcloud

import (
	"io"
	"io/ioutil"
	"os"
	"testing"
	// "time"
)

func TestAws (t *testing.T) {
	params := AwsProviderParams("AWS", os.Getenv("ZCLOUD_AWS_KEY_ID"), os.Getenv("ZCLOUD_AWS_SECRET_KEY"), os.Getenv("ZCLOUD_AWS_REGION"))
	p, err := NewProvider(params)
	if err != nil {
		t.Fatal(err)
	}
	t.Run("Provider", func (t *testing.T) {
		testProvider(t, p)
	})
}

func TestGCloud (t *testing.T) {
	params := GCloudProviderParams("GCLOUD", os.Getenv("ZCLOUD_GCLOUD_PROJECT"))
	p, err := NewProvider(params)
	if err != nil {
		t.Fatal(err)
	}
	t.Run("Provider", func (t *testing.T) {
		testProvider(t, p)
	})
}

const testBucketName = "zcloud-testing-go" 

func testProvider (t *testing.T, p Provider) {
	b := p.Bucket(testBucketName)
	err := b.Create()
	if err != nil {
		t.Fatalf("Bucket create %v", err)
	}
	bs, err := p.Buckets()
	if err != nil {
		t.Fatal("Buckets", err)
	}
	hasBucket := false
	for _, b = range bs {
		if b.Name() == testBucketName {
			hasBucket = true
		}
	}
	if !hasBucket {
		t.Fatalf("Bucket %v isn't in Buckets %v", testBucket, bs)
	}
	t.Run("Bucket", func (t *testing.T) {
		testBucket(t, b)
	})
	err = b.Delete()
	if err != nil {
		t.Fatalf("Bucket delete %v", err)
	}
	bs, err = p.Buckets()
	if err != nil {
		t.Fatalf("Buckets %v", err)
	}
	for _, b := range bs {
		if b.Name() == testBucketName {
			t.Fatalf("%v was not deleted", testBucketName)
		}
	}
}

const testObjectKey = "test.txt"
var testObjectDataConst = []byte{'n', 'a', 't', 'e', 'e'}

func testBucket (t *testing.T, b Bucket) {
	o := b.Object(testObjectKey)
	w, err := o.Writer()
	if err != nil {
		t.Fatalf("Writer %v", err)
	}
	n, err := w.Write(testObjectDataConst)
	if n != len(testObjectDataConst) {
		t.Fatalf("Only wrote %v of %v bytes: %v", n, len(testObjectDataConst), testObjectDataConst)
	}
	if err != nil {
		t.Fatalf("Error when writing object %v", err)
	}
	err = w.Close()
	if err != nil {
		t.Fatalf("Error when closing object after writing %v", err)
	}
	os, err := b.Objects()
	if err != nil {
		t.Fatalf("Objects %v", err)
	}
	hasObject := false
	for _, o := range os {
		if o.Key() == testObjectKey {
			hasObject = true
		}
	}
	if !hasObject {
		t.Fatalf("Objects %v did not contain %v", os, testObjectKey)
	}
	o2 := b.Object(testObjectKey)
	t.Run("Object", func (t *testing.T) {
		testObject(t, o, o2)
	})
	err = o.Delete()
	if err != nil {
		t.Fatal("Object Delete %v", err)
	}
	os, err = b.Objects()
	for _, o := range os {
		if o.Key() == testObjectKey {
			t.Fatalf("%v was not deleted", testObjectKey)
		}
	}
}

func testObject (t *testing.T, o Object, o2 Object) {
	// prevTime := time.Now()
	w, err := o.Writer()
	if err != nil {
		t.Fatalf("Writer %v", err)
	}
	n, err := w.Write(testObjectDataConst)
	if n != len(testObjectDataConst) {
		t.Fatalf("Only wrote %v bytes instead of %v", n, len(testObjectDataConst))
	}
	if err != nil {
		t.Fatalf("Object Write %v", err)
	}
	err = w.Close()
	if err != nil {
		t.Fatalf("Writer Close %v", err)
	}
	// l, err := o.LastModified()
	// if err != nil {
	// 	t.Fatalf("Last Modified %v", err)
	// }
	// postTime := time.Now()
	// if !l.After(prevTime) || !l.Before(postTime) {
	// 	t.Fatalf("Last Modified %v, Previous %v, Post %v", l, prevTime, postTime)
	// }
	s, err := o.Size()
	if err != nil {
		t.Fatalf("Size %v", err)
	}
	if s != len(testObjectDataConst) {
		t.Fatalf("Size is %v but should be %v", s, len(testObjectDataConst))
	}
	r, err := o2.Reader()
	if err != nil {
		t.Fatalf("Reader %v", err)
	}
	b, err := ioutil.ReadAll(r)
	if len(b) != len(testObjectDataConst) {
		t.Fatalf("Only read %v of %v bytes", n, len(testObjectDataConst))
	}
	if err != nil && err != io.EOF {
		t.Fatalf("Read %v", err)
	}
	err = r.Close()
	if err != nil {
		t.Fatalf("Read Close %v", err)
	}
	for i := range testObjectDataConst {
		if b[i] != testObjectDataConst[i] {
			t.Fatalf("Read %v, but should be %v", b, testObjectDataConst)
		}
	}
}
