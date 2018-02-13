package zcloud

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func newAwsBucket (name string, p *awsProvider) awsBucket {
	return awsBucket{
		name: name,
		p: p,
	}
}

type awsBucket struct {
	name string
	p *awsProvider
}

func (b awsBucket) Create () error {
	s3svc := s3.New(b.p.session)
	input := &s3.CreateBucketInput{
		Bucket: aws.String(b.Name()),
	}
	_, err := s3svc.CreateBucket(input)
	if err != nil {
		return err
	}
	return nil
}

func (b awsBucket) Delete () error {
	s3svc := s3.New(b.p.session)
	input := &s3.DeleteBucketInput{
		Bucket: aws.String(b.Name()),
	}
	_, err := s3svc.DeleteBucket(input)
	if err != nil {
		return err
	}
	return nil
}

func (b awsBucket) Name () string {
	return b.name
}

func (b awsBucket) Object (key string) Object {
	return newAwsObject(b.Name(), key, &b)
}

func (b awsBucket) Objects () ([]Object, error) {
	return b.ObjectsQuery(nil)
}

func (b awsBucket) ObjectsQuery (q *ObjectsQueryParams) ([]Object, error) {
	s3svc := s3.New(b.p.session)
	os := []Object{}
	n := b.Name()
	var prefix string
	if q != nil {
		prefix = q.Prefix
	}
	params := &s3.ListObjectsV2Input{
		Bucket: &n,
		Prefix: aws.String(prefix),
	}
	err := s3svc.ListObjectsV2Pages(params,
		func(page *s3.ListObjectsV2Output, lastPage bool) bool {
			for _, c := range page.Contents {
				os = append(os, newAwsObject(b.Name(), *c.Key, &b))
			}
			return true
		})
	if err != nil {
		return os, err
	}
	return os, nil
}
