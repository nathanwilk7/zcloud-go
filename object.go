package zcloud

import (
	"fmt"
	"io"
	"time"
)

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

type ErrObjectDoesNotExist struct {
	bucket string
	key string
}

func (e ErrObjectDoesNotExist) Error () string {
	return fmt.Sprintf("Object does not exist. Bucket: %s, Key: %s", e.bucket, e.key)
}
