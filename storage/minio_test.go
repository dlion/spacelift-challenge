package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"testing"

	"github.com/docker/docker/client"
	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	minioTest "github.com/testcontainers/testcontainers-go/modules/minio"
)

func TestMinio(t *testing.T) {

	t.Run("It should get information from the minio instance", func(t *testing.T) {
		ctx := context.Background()

		minioContainer, err := minioTest.RunContainer(ctx,
			testcontainers.WithImage("minio/minio:RELEASE.2024-01-16T16-07-38Z"),
			testcontainers.WithEnv(map[string]string{
				"MINIO_ACCESS_KEY": "test",
				"MINIO_SECRET_KEY": "test",
			}),
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

		cli, err := client.NewClientWithOpts(client.FromEnv)
		if err != nil {
			t.Fatalf("Failed to create Docker client: %v", err)
		}

		instances, err := GetMinioInstancesFromDocker(cli)
		assert.NoError(t, err)

		assert.NotEmpty(t, instances)
		assert.Equal(t, 1, len(instances))
		assert.Contains(t, instances[0].URL, "http://")
		assert.Equal(t, instances[0].Access, "test")
		assert.Equal(t, instances[0].Secret, "test")
	})

	t.Run("It should upload a file to a minio instance", func(t *testing.T) {
		ctx := context.Background()

		minioContainer, err := minioTest.RunContainer(ctx,
			testcontainers.WithImage("minio/minio:RELEASE.2024-01-16T16-07-38Z"),
			testcontainers.WithEnv(map[string]string{
				"MINIO_ACCESS_KEY": "test",
				"MINIO_SECRET_KEY": "test",
			}),
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
		ok, err := minioService.UploadFile(body, "123")
		assert.NoError(t, err)
		assert.Equal(t, ok, true)

		ok, err = minioClient.BucketExists(ctx, BUCKET_NAME)
		assert.NoError(t, err)
		assert.True(t, ok)

		obj, err := minioClient.GetObject(ctx, BUCKET_NAME, "123", minio.GetObjectOptions{})
		assert.NoError(t, err)

		buffer, err := io.ReadAll(obj)
		assert.NoError(t, err)
		assert.Equal(t, body, buffer)
	})

	t.Run("It should return an error if the file already exist", func(t *testing.T) {
		ctx := context.Background()
		minioContainer, err := minioTest.RunContainer(ctx,
			testcontainers.WithImage("minio/minio:RELEASE.2024-01-16T16-07-38Z"),
			testcontainers.WithEnv(map[string]string{
				"MINIO_ACCESS_KEY": "test",
				"MINIO_SECRET_KEY": "test",
			}),
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
		ok, err := minioService.UploadFile(anotherBody, "123")
		assert.NoError(t, err)
		assert.Equal(t, ok, true)

		obj, err := minioClient.GetObject(ctx, BUCKET_NAME, "123", minio.GetObjectOptions{})
		assert.NoError(t, err)
		buffer, err := io.ReadAll(obj)
		assert.NoError(t, err)
		assert.Equal(t, anotherBody, buffer)
	})
}
