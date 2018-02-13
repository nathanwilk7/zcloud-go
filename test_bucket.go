package zcloud

import (
	"fmt"
)

func NewTestBucket (name string, p *TestProvider) TestBucket {
	return TestBucket{
		NameField: name,
		P: p,
	}
}

type TestBucket struct {
	NameField string
	P *TestProvider
}

func (b TestBucket) Create () error {
	return nil
}

func (b TestBucket) Delete () error {
	index := -1
	for i := range b.P.BucketsField {
		if b.P.BucketsField[i].Name() == b.Name() {
			index = i
		}
	}
	if index == -1 {
		return fmt.Errorf("Bucket %v does not exist in %v", b.Name, b.P.BucketsField)
	}
	copy(b.P.BucketsField[index:], b.P.BucketsField[index + 1:])
	b.P.BucketsField[len(b.P.BucketsField) - 1] = TestBucket{}
	b.P.BucketsField = b.P.BucketsField[:len(b.P.BucketsField) - 1]
	return nil
}

func (b TestBucket) Name () string {
	return b.NameField
}

func (b TestBucket) Object (key string) Object {
	return NewTestObject(b.Name(), key, &b)
}

func (b TestBucket) Objects () ([]Object, error) {
	return b.ObjectsQuery(nil)
}

func (b TestBucket) ObjectsQuery (q *ObjectsQueryParams) ([]Object, error) {
	o, _ := b.Objects()
	return o, nil
}
