package storage

import (
	"bytes"
	"context"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	BUCKET_NAME = "uploads"
)

type MinioService struct {
	client *minio.Client
}

func NewMinioClient(endpoint, accessKey, secretKey string) (*minio.Client, error) {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	return minioClient, nil
}

func NewMinioService(client *minio.Client) *MinioService {
	return &MinioService{client: client}
}

func (m *MinioService) UploadFile(body []byte, hash string) (bool, error) {
	ctx := context.Background()

	_, err := createBucketIfNotExist(ctx, m.client)
	if err != nil {
		return false, err
	}

	_, err = m.client.PutObject(context.Background(), BUCKET_NAME, hash, bytes.NewReader(body), int64(len(body)), minio.PutObjectOptions{})
	if err != nil {
		return false, err
	}
	log.Printf("- %s - The object with id %s has been uploaded into the bucket named %s", m.client.EndpointURL(), hash, BUCKET_NAME)
	return true, nil
}

func createBucketIfNotExist(ctx context.Context, client *minio.Client) (bool, error) {
	if ok, err := client.BucketExists(ctx, BUCKET_NAME); !ok {
		if err != nil {
			return false, err
		}
		client.MakeBucket(ctx, BUCKET_NAME, minio.MakeBucketOptions{})
		log.Printf("The Bucket %s has been created", BUCKET_NAME)
	}
	return true, nil
}
