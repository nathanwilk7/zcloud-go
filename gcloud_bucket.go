package zcloud

import (
	"google.golang.org/api/iterator"

	gs "cloud.google.com/go/storage"
)

func NewGCloudBucket (name string, p *GCloudProvider) GCloudBucket {
	return GCloudBucket{
		name: name,
		p: p,
	}
}

type GCloudBucket struct {
	name string
	p *GCloudProvider
}

func (b GCloudBucket) Create () error {
	if err := b.p.client.Bucket(b.Name()).Create(b.p.context, b.p.project, nil); err != nil {
		return err
	}
	return nil
}

func (b GCloudBucket) Delete () error {
	if err := b.p.client.Bucket(b.Name()).Delete(b.p.context); err != nil {
		return err
	}
	return nil
}

func (b GCloudBucket) Name () string {
	return b.name
}

func (b GCloudBucket) Object (key string) Object {
	return NewGCloudObject(b.Name(), key, &b)
}

func (b GCloudBucket) Objects () ([]Object, error) {
	return b.ObjectsQuery(nil)
}

func (b GCloudBucket) ObjectsQuery (q *ObjectsQueryParams) ([]Object, error) {
	var gsq *gs.Query
	if q != nil {
		gsq = &gs.Query{
			Prefix: q.Prefix,
		}
	}
	it := b.p.client.Bucket(b.Name()).Objects(b.p.context, gsq)
	os := []Object{}
	for {
		o, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return os, err
		}
		os = append(os, NewGCloudObject(b.Name(), o.Name, &b))
	}
	return os, nil
}
