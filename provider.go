package zcloud

import (
	"fmt"
	"io"
	"time"
)

func AwsProviderParams (name, keyID, secretKey, region string) ProviderParams {
	return ProviderParams{
		Name: name,
		AwsId: keyID,
		AwsSecret: secretKey,
		AwsRegion: region,
	}
}

func GCloudProviderParams (name, projectID string) ProviderParams {
	return ProviderParams{
		Name: name,
		GCloudProjectID: projectID,
	}
}

type ProviderParams struct {
	Name string
	AwsId, AwsSecret, AwsRegion string
	GCloudProjectID string
}

func NewProvider (params ProviderParams) (Provider, error) {
	switch params.Name {
	case "AWS":
		return newAwsProvider(params.AwsId, params.AwsSecret, params.AwsRegion), nil
	case "GCLOUD":
		return newGCloudProvider(params.GCloudProjectID), nil
	case "TEST":
		return NewTestProvider(), nil
	}
	return nil, fmt.Errorf("%s is not a valid provider name", params.Name)
}

type Provider interface {
	StorageProvider
}

func NewStorageProvider (params ProviderParams) (StorageProvider, error) {
	switch params.Name {
		// StorageProvider's go here...
	}
	if p, err := NewProvider(params); err == nil {
		return p, nil
	}
	return nil, fmt.Errorf("%s is not a valid provider name", params.Name)
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
	ObjectsQuery (query *ObjectsQueryParams) ([]Object, error)
}

type ObjectsQueryParams struct {
	Prefix string
}

type Object interface {
	CopyTo (Object) error
	Delete () error
	Info () (ObjectInfo, error)
	Key () string
	Reader () (io.ReadCloser, error)
	Writer () (io.WriteCloser, error)
}

type ObjectInfo interface {
	LastModified () time.Time
	Size () int
}

type ObjectCopier interface {
	Copy () error
}
