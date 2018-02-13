package zcloud

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func newAwsProvider(id, secret, region string) awsProvider {
	s, err := getSession(id, secret, region)
	if err != nil {
		panic(err)
	}
	return awsProvider{
		session: s,
	}
}

type awsProvider struct {
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

func (p awsProvider) Buckets () ([]Bucket, error) {
	svc := s3.New(p.session)
	input := &s3.ListBucketsInput{}
	awsBuckets, err := svc.ListBuckets(input)
	if err != nil {
		return nil, err
	}
	buckets := make([]Bucket, len(awsBuckets.Buckets))
	for i, awsBucket := range awsBuckets.Buckets {
		buckets[i] = newAwsBucket(*awsBucket.Name, &p)
	}
	return buckets, nil
}

func (p awsProvider) Bucket (name string) Bucket {
	return newAwsBucket(name, &p)
}
