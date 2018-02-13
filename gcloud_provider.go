package zcloud

import (
	"golang.org/x/net/context"
	
	gs "cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

func newGCloudProvider (projectID string) gCloudProvider {
	c := context.Background()
	client, err := gs.NewClient(c)
	if err != nil {
		panic(err)
	}
	return gCloudProvider{
		context: c,
		client: client,
		project: projectID,
	}
}

type gCloudProvider struct {
	context context.Context
	client *gs.Client
	project string
}

func (p gCloudProvider) Buckets () ([]Bucket, error) {
	it := p.client.Buckets(p.context, p.project)
	bs := []Bucket{}
	for {
		b, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		bs = append(bs, newGCloudBucket(b.Name, &p))
	}
	return bs, nil
}

func (p gCloudProvider) Bucket (name string) Bucket {
	return newGCloudBucket(name, &p)
}
