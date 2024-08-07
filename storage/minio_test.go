package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"testing"

	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	minioTest "github.com/testcontainers/testcontainers-go/modules/minio"
)

func TestMinio(t *testing.T) {
	t.Run("It should upload a file to a minio instance", func(t *testing.T) {
		ctx := context.Background()

		minioContainer, err := minioTest.RunContainer(ctx,
			testcontainers.WithImage("minio/minio:RELEASE.2024-01-16T16-07-38Z"),
			testcontainers.CustomizeRequest(
				testcontainers.GenericContainerRequest{
					ContainerRequest: testcontainers.ContainerRequest{
						Labels: map[string]string{
							"com.docker.compose.service": "amazin-object-storage-node-1",
						},
					},
				},
			),
		)
		if err != nil {
			log.Fatalf("failed to start container: %s", err)
		}
		defer func() {
			if err := minioContainer.Terminate(ctx); err != nil {
				log.Fatalf("failed to terminate container: %s", err)
			}
		}()

		ip, err := minioContainer.ContainerIP(ctx)
		assert.NoError(t, err)

		minioClient, err := NewMinioClient(fmt.Sprintf("%s:9000", ip), "minioadmin", "minioadmin")
		assert.NoError(t, err)

		minioService := NewMinioService(minioClient)

		body := []byte("Hello World")
		err = minioService.UploadFile(bytes.NewReader(body), "123")
		assert.NoError(t, err)

		ok, err := minioClient.BucketExists(ctx, BUCKET_NAME)
		assert.NoError(t, err)
		assert.True(t, ok)

		obj, err := minioClient.GetObject(ctx, BUCKET_NAME, "123", minio.GetObjectOptions{})
		assert.NoError(t, err)

		buffer, err := io.ReadAll(obj)
		assert.NoError(t, err)
		assert.Equal(t, body, buffer)
	})

	t.Run("It should return override the file if the file already exist", func(t *testing.T) {
		ctx := context.Background()
		minioContainer, err := minioTest.RunContainer(ctx,
			testcontainers.WithImage("minio/minio:RELEASE.2024-01-16T16-07-38Z"),
			testcontainers.CustomizeRequest(
				testcontainers.GenericContainerRequest{
					ContainerRequest: testcontainers.ContainerRequest{
						Labels: map[string]string{
							"com.docker.compose.service": "amazin-object-storage-node-1",
						},
					},
				},
			),
		)
		if err != nil {
			log.Fatalf("failed to start container: %s", err)
		}
		defer func() {
			if err := minioContainer.Terminate(ctx); err != nil {
				log.Fatalf("failed to terminate container: %s", err)
			}
		}()
		ip, err := minioContainer.ContainerIP(ctx)
		assert.NoError(t, err)
		minioClient, err := NewMinioClient(fmt.Sprintf("%s:9000", ip), "minioadmin", "minioadmin")
		assert.NoError(t, err)
		err = minioClient.MakeBucket(ctx, BUCKET_NAME, minio.MakeBucketOptions{})
		assert.NoError(t, err)
		body := []byte("Hello World")
		minioClient.PutObject(ctx, BUCKET_NAME, "123", bytes.NewReader(body), int64(len(body)), minio.PutObjectOptions{})

		minioService := NewMinioService(minioClient)
		anotherBody := []byte("Another World")
		err = minioService.UploadFile(bytes.NewReader(anotherBody), "123")
		assert.NoError(t, err)

		obj, err := minioClient.GetObject(ctx, BUCKET_NAME, "123", minio.GetObjectOptions{})
		assert.NoError(t, err)
		buffer, err := io.ReadAll(obj)
		assert.NoError(t, err)
		assert.Equal(t, anotherBody, buffer)
	})

	t.Run("It should get a file from a minio instance", func(t *testing.T) {
		ctx := context.Background()

		minioContainer, err := minioTest.RunContainer(ctx,
			testcontainers.WithImage("minio/minio:RELEASE.2024-01-16T16-07-38Z"),
			testcontainers.CustomizeRequest(
				testcontainers.GenericContainerRequest{
					ContainerRequest: testcontainers.ContainerRequest{
						Labels: map[string]string{
							"com.docker.compose.service": "amazin-object-storage-node-1",
						},
					},
				},
			),
		)
		if err != nil {
			log.Fatalf("failed to start container: %s", err)
		}
		defer func() {
			if err := minioContainer.Terminate(ctx); err != nil {
				log.Fatalf("failed to terminate container: %s", err)
			}
		}()

		ip, err := minioContainer.ContainerIP(ctx)
		assert.NoError(t, err)

		minioClient, err := NewMinioClient(fmt.Sprintf("%s:9000", ip), "minioadmin", "minioadmin")
		assert.NoError(t, err)
		err = minioClient.MakeBucket(ctx, BUCKET_NAME, minio.MakeBucketOptions{})
		assert.NoError(t, err)
		body := []byte("Hello World")
		_, err = minioClient.PutObject(ctx, BUCKET_NAME, "1916298011", bytes.NewReader(body), int64(len(body)), minio.PutObjectOptions{})
		assert.NoError(t, err)

		minioService := NewMinioService(minioClient)
		file, err := minioService.GetFile("1916298011")
		assert.NoError(t, err)

		buffer, err := io.ReadAll(file)
		assert.NoError(t, err)

		assert.Equal(t, []byte("Hello World"), buffer)
	})

	t.Run("It should return an error if the file we are trying to get doesn't exist", func(t *testing.T) {
		ctx := context.Background()

		minioContainer, err := minioTest.RunContainer(ctx,
			testcontainers.WithImage("minio/minio:RELEASE.2024-01-16T16-07-38Z"),
			testcontainers.CustomizeRequest(
				testcontainers.GenericContainerRequest{
					ContainerRequest: testcontainers.ContainerRequest{
						Labels: map[string]string{
							"com.docker.compose.service": "amazin-object-storage-node-1",
						},
					},
				},
			),
		)
		if err != nil {
			log.Fatalf("failed to start container: %s", err)
		}
		defer func() {
			if err := minioContainer.Terminate(ctx); err != nil {
				log.Fatalf("failed to terminate container: %s", err)
			}
		}()

		ip, err := minioContainer.ContainerIP(ctx)
		assert.NoError(t, err)

		minioClient, err := NewMinioClient(fmt.Sprintf("%s:9000", ip), "minioadmin", "minioadmin")
		assert.NoError(t, err)
		err = minioClient.MakeBucket(ctx, BUCKET_NAME, minio.MakeBucketOptions{})
		assert.NoError(t, err)

		minioService := NewMinioService(minioClient)
		_, err = minioService.GetFile("1916298011")
		assert.Error(t, err)
	})
}
