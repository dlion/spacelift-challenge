package storage

import (
	"context"
	"log"
	"testing"

	"github.com/docker/docker/client"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/minio"
)

func TestMinio(t *testing.T) {

	ctx := context.Background()

	minioContainer, err := minio.RunContainer(ctx,
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
		log.Fatalf("failed tso start container: %s", err)
	}
	defer func() {
		if err := minioContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	}()

	// Create Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		t.Fatalf("Failed to create Docker client: %v", err)
	}

	instances, err := GetMinioInstancesFromDocker(cli)
	assert.NoError(t, err)

	// Assertions
	assert.NotEmpty(t, instances)
	assert.Equal(t, 1, len(instances))
	assert.Contains(t, instances[0].URL, "http://")
	assert.Equal(t, instances[0].Access, "test")
	assert.Equal(t, instances[0].Secret, "test")
}
