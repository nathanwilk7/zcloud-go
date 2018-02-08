package zcloud

import (
	"os"
	"testing"
)

func TestBuckets (t *testing.T) {
	p := NewAwsProvider(os.Getenv("ZCLOUD_AWS_KEY_ID"), os.Getenv("ZCLOUD_AWS_SECRET_KEY"), os.Getenv("ZCLOUD_AWS_REGION"))
}
