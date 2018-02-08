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
}

type Bucket interface {
	Create () error
	Delete () error
	Name () string
	Object (key string) Object
	Objects () ([]Object, error)
}

type Object interface {
	Key () string
	LastModified () (time.Time, error)
	Reader () (io.ReadCloser, error)
	Writer () (io.WriteCloser, error)
	Size () (int, error)
}
