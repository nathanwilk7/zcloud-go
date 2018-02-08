package zcloud

import (
	"fmt"
	"io"
	"time"
)

type NewProviderParams struct {
	Name string
	AwsId, AwsSecret, AwsRegion string
}

func NewProvider (params NewProviderParams) (Provider, error) {
	switch params.Name {
	case "AWS":
		return NewAwsProvider(params.AwsId, params.AwsSecret, params.AwsRegion), nil
	}
	return nil, fmt.Errorf("%s is not a valid provider name", params.Name)
}

type Provider interface {
	StorageProvider
}

type StorageProvider interface {
	Buckets () ([]Bucket, error)
	Bucket (name string) Bucket
	RemoveBucket (name string) error
}

type Bucket interface {
	Create () error
	Name () string
	Object (key string) Object
	Objects () []Object
}

type Object interface {
	Key () string
	LastModified () time.Time
	Reader () io.ReadCloser
	Writer () io.WriteCloser
	Size () int
}
