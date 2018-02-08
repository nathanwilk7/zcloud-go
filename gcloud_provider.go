package zcloud

import (
	"golang.org/x/net/context"
	
	gs "cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

func NewGCloudProvider(projectID string) GCloudProvider {
	c := context.Background()
	client, err := gs.NewClient(c)
	if err != nil {
		panic(err)
	}
	return GCloudProvider{
		context: c,
		client: client,
		project: projectID,
	}
}

type GCloudProvider struct {
	context context.Context
	client *gs.Client
	project string
}

func (p GCloudProvider) Buckets () ([]Bucket, error) {
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
		bs = append(bs, NewGCloudBucket(b.Name, &p))
	}
	return bs, nil
}

func (p GCloudProvider) Bucket (name string) Bucket {
	return NewGCloudBucket(name, &p)
}
