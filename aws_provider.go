package zcloud

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func NewAwsProvider(id, secret, region string) AwsProvider {
	s, err := getSession(id, secret, region)
	if err != nil {
		panic(err)
	}
	return AwsProvider{
		//Id: id,
		//Secret: secret,
		//Region: region,
		session: s,
	}
}

type AwsProvider struct {
	//Id, Secret, Region string
	session *session.Session
}

const defaultToken = ""

func getSession (id, secret, region string) (*session.Session, error) {
	sess, err := session.NewSessionWithOptions(
		session.Options{
			Config: aws.Config{
				Region: aws.String(region),
				Credentials: credentials.NewStaticCredentials(id, secret, defaultToken),
			},
		},
	)
	if err != nil {
		return nil, err
	}
	return sess, nil
}

func (p AwsProvider) Buckets () ([]Bucket, error) {
	svc := s3.New(p.session)
	input := &s3.ListBucketsInput{}
	awsBuckets, err := svc.ListBuckets(input)
	if err != nil {
		return nil, err
	}
	buckets := make([]Bucket, len(awsBuckets.Buckets))
	for i, awsBucket := range awsBuckets.Buckets {
		buckets[i] = NewAwsBucket(*awsBucket.Name, &p)
	}
	return buckets, nil
}

func (p AwsProvider) Bucket (name string) Bucket {
	return NewAwsBucket(name, &p)
}
