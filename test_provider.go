package zcloud

import (
	
)

func NewTestProvider() TestProvider {
	return TestProvider{
		BucketsField: []Bucket{},
	}
}

type TestProvider struct {
	BucketsField []Bucket
}

func (p TestProvider) Buckets () ([]Bucket, error) {
	return p.BucketsField, nil
}

func (p TestProvider) Bucket (name string) Bucket {
	b := NewTestBucket(name, &p)
	p.BucketsField = append(p.BucketsField, b)
	return b
}
