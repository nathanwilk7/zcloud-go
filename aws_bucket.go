package zcloud

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func NewAwsBucket (name string) AwsBucket {
	return AwsBucket{
		name: name,
	}
}

type AwsBucket struct {
	name string
}
func (b AwsBucket) Create () error {
	svc := s3.New(session.New())
	input := &s3.CreateBucketInput{
		Bucket: aws.String(b.Name()),
	}
	_, err := svc.CreateBucket(input)
	if err != nil {
		return err
	}
	return nil
}

func (b AwsBucket) Name () string {
	return b.name
}

func (b AwsBucket) Object (key string) Object {
	return NewAwsObject(b.Name(), key)
}

func (b AwsBucket) Objects () []Object {
	return nil
}
