package zcloud

import (
	"google.golang.org/api/iterator"

	gs "cloud.google.com/go/storage"
)

func newGCloudBucket (name string, p *gCloudProvider) gCloudBucket {
	return gCloudBucket{
		name: name,
		p: p,
	}
}

type gCloudBucket struct {
	name string
	p *gCloudProvider
}

func (b gCloudBucket) Create () error {
	if err := b.p.client.Bucket(b.Name()).Create(b.p.context, b.p.project, nil); err != nil {
		return err
	}
	return nil
}

func (b gCloudBucket) Delete () error {
	if err := b.p.client.Bucket(b.Name()).Delete(b.p.context); err != nil {
		return err
	}
	return nil
}

func (b gCloudBucket) Name () string {
	return b.name
}

func (b gCloudBucket) Object (key string) Object {
	return newGCloudObject(b.Name(), key, &b)
}

func (b gCloudBucket) Objects () ([]Object, error) {
	return b.ObjectsQuery(nil)
}

func (b gCloudBucket) ObjectsQuery (q *ObjectsQueryParams) ([]Object, error) {
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
		os = append(os, newGCloudObject(b.Name(), o.Name, &b))
	}
	return os, nil
}
