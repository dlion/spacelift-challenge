package storage

import (
	"context"
	"io"
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

func (m *MinioService) UploadFile(body io.Reader, hash string) error {
	ctx := context.Background()

	_, err := createBucketIfNotExist(ctx, m.client)
	if err != nil {
		return err
	}

	_, err = m.client.PutObject(context.Background(), BUCKET_NAME, hash, body, -1, minio.PutObjectOptions{})
	if err != nil {
		log.Printf("- %s - Couldn't upload the object with id %s into the bucket named: %s\n", m.client.EndpointURL(), hash, BUCKET_NAME)
		return err
	}
	log.Printf("- %s - The object with id %s has been uploaded into the bucket named: %s\n", m.client.EndpointURL(), hash, BUCKET_NAME)
	return nil
}

func (m *MinioService) GetFile(hash string) (io.Reader, error) {
	file, err := m.client.GetObject(context.Background(), BUCKET_NAME, hash, minio.GetObjectOptions{})
	if err != nil {
		log.Printf("- %s - Couldn't get the file %s from the bucket %s", m.client.EndpointURL(), hash, BUCKET_NAME)
		return nil, err
	}

	log.Printf("- %s - Got the file name %s from the bucket: %s", m.client.EndpointURL(), hash, BUCKET_NAME)
	return file, nil
}

func createBucketIfNotExist(ctx context.Context, client *minio.Client) (bool, error) {
	if ok, err := client.BucketExists(ctx, BUCKET_NAME); !ok {
		if err != nil {
			log.Printf("- %s - Couldn't create the bucket named %s ", client.EndpointURL(), BUCKET_NAME)
			return false, err
		}
		client.MakeBucket(ctx, BUCKET_NAME, minio.MakeBucketOptions{})
		log.Printf("- %s - The Bucket %s has been created", client.EndpointURL(), BUCKET_NAME)
	}
	return true, nil
}
