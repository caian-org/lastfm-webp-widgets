//go:build save_s3

package storage

import (
	// standard
	"bytes"
	"context"
	"fmt"
	"os"
	"time"

	// 3rd-party
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	// internal
	. "github.com/caian-org/lastfm-webp-widgets/internal/logger"
)

const (
	S3_BUCKET_ENV_NAME = "S3_BUCKET_NAME"
)

var (
	s3c    *s3.S3
	bucket *string
)

func storageInit() {
	if b, ok := os.LookupEnv(S3_BUCKET_ENV_NAME); ok {
		bucket = &b
	} else {
		panic(fmt.Sprintf("unset env %s", S3_BUCKET_ENV_NAME))
	}

	config := aws.NewConfig().WithRegion("us-east-1")
	session := session.Must(session.NewSession(config))
	s3 := s3.New(session)

	s3c = s3
}

func storageSave(key string, buf bytes.Buffer) {
	data := bytes.NewReader(buf.Bytes())

	ctx := context.Background()
	ctx, cancelFunc := context.WithTimeout(ctx, 10*time.Second)
	defer cancelFunc()

	params := s3.PutObjectInput{
		Body:         data,
		Bucket:       bucket,
		Key:          aws.String(key),
		ACL:          aws.String("public-read"),
		ContentType:  aws.String("image/webp"),
		CacheControl: aws.String("max-age=0"),
	}

	if _, err := s3c.PutObjectWithContext(ctx, &params); err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == request.CanceledErrorCode {
			fmt.Fprintf(os.Stderr, "upload canceled due to timeout, %v\n", err)
		} else {
			fmt.Fprintf(os.Stderr, "failed to upload object, %v\n", err)
		}

		panic(err)
	}

	Log.Info("object uploaded", "bucket", *bucket, "key", key)
}
