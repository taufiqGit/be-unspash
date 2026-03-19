package db

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// InitS3 menginisialisasi S3 client untuk BiznetGio NEO Object Storage.
//
// Env vars yang dibutuhkan:
//   - S3_ENDPOINT   : URL endpoint BiznetGio, contoh: https://nos.wjv-1.neo.id
//   - S3_ACCESS_KEY : Access key dari dashboard NEO Object Storage
//   - S3_SECRET_KEY : Secret key dari dashboard NEO Object Storage
//   - S3_REGION     : Region BiznetGio, contoh: wjv-1
//   - S3_BUCKET     : Nama bucket yang digunakan
func InitS3() (*s3.Client, error) {
	endpoint := os.Getenv("S3_ENDPOINT")
	accessKey := os.Getenv("S3_ACCESS_KEY")
	secretKey := os.Getenv("S3_SECRET_KEY")
	region := os.Getenv("S3_REGION")

	if endpoint == "" {
		return nil, fmt.Errorf("S3_ENDPOINT is required")
	}
	if accessKey == "" {
		return nil, fmt.Errorf("S3_ACCESS_KEY is required")
	}
	if secretKey == "" {
		return nil, fmt.Errorf("S3_SECRET_KEY is required")
	}
	if region == "" {
		return nil, fmt.Errorf("S3_REGION is required")
	}

	client := s3.New(s3.Options{
		// Arahkan ke endpoint BiznetGio NEO Object Storage
		BaseEndpoint: aws.String(endpoint),

		// Region sesuai zona BiznetGio (mis. wjv-1)
		Region: region,

		// Kredensial statis dari dashboard NEO Object Storage
		Credentials: aws.NewCredentialsCache(
			credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		),

		// BiznetGio mendukung path-style: https://<endpoint>/<bucket>/<key>
		UsePathStyle: true,
	})

	// Verifikasi koneksi dengan ping ke bucket
	bucket := os.Getenv("S3_BUCKET")
	if bucket != "" {
		_, err := client.HeadBucket(context.Background(), &s3.HeadBucketInput{
			Bucket: aws.String(bucket),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to connect to S3 bucket %q: %w", bucket, err)
		}
	}

	return client, nil
}
