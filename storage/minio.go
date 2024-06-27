package storage

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/dlion/spacelift-challenge/docker"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	INSTANCE_NAME = "amazin-object-storage"
	BUCKET_NAME   = "uploads"
)

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

type Minio struct {
	client *minio.Client
}

func NewMinioService(client *minio.Client) *Minio {
	return &Minio{client: client}
}

func (m *Minio) UploadFile(body []byte, hash string) (bool, error) {
	ctx := context.Background()

	_, err := createBucketIfNotExist(ctx, m.client)
	if err != nil {
		return false, err
	}

	_, err = m.client.PutObject(context.Background(), BUCKET_NAME, hash, bytes.NewReader(body), int64(len(body)), minio.PutObjectOptions{})
	if err != nil {
		return false, err
	}
	return true, nil
}

func createBucketIfNotExist(ctx context.Context, client *minio.Client) (bool, error) {
	if ok, err := client.BucketExists(ctx, BUCKET_NAME); !ok {
		if err != nil {
			return false, err
		}
		client.MakeBucket(ctx, BUCKET_NAME, minio.MakeBucketOptions{})
	}
	return true, nil
}

type MinioInstance struct {
	Name   string
	URL    string
	Access string
	Secret string
}

func GetMinioInstancesFromDocker(client *client.Client) ([]MinioInstance, error) {
	containers, err := docker.GetContainerList(client, context.Background())
	if err != nil {
		return nil, err
	}

	var instances []MinioInstance
	for _, container := range containers {
		if minioInstanceExist(container) {
			containerInspect, err := docker.InspectContainerByID(client, container.ID)
			if err != nil {
				return nil, err
			}

			mapKey := docker.GetAccessSecretKeyFromTheContainer(containerInspect)
			ip := docker.GetIPAddressFromTheContainer(containerInspect)
			port := docker.GetPortFromTheContainer(containerInspect)

			instances = append(instances, MinioInstance{
				Name:   container.Names[0],
				URL:    "http://" + ip + ":" + port,
				Access: mapKey["MINIO_ROOT_USER"],
				Secret: mapKey["MINIO_ROOT_PASSWORD"],
			})
		}
	}

	return instances, nil
}

func minioInstanceExist(container types.Container) bool {
	name, ok := container.Labels["com.docker.compose.service"]
	return ok && strings.Contains(name, fmt.Sprintf("%s-node", INSTANCE_NAME))
}
