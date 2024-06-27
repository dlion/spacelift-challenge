package docker

import (
	"context"
	"log"
	"testing"

	"github.com/docker/docker/client"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/minio"
)

func TestDocker(t *testing.T) {
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

	t.Run("It should get the port from a minio instance", func(t *testing.T) {
		cli, err := client.NewClientWithOpts(client.FromEnv)
		if err != nil {
			t.Fatalf("Failed to create Docker client: %v", err)
		}

		containerInspect, err := cli.ContainerInspect(ctx, minioContainer.GetContainerID())
		assert.NoError(t, err)

		instanceIP := GetIPAddressFromTheContainer(containerInspect)

		ip, err := minioContainer.ContainerIP(ctx)
		assert.NoError(t, err)
		assert.Equal(t, ip, instanceIP)
	})

	t.Run("It should get the ip from a minio instance", func(t *testing.T) {
		cli, err := client.NewClientWithOpts(client.FromEnv)
		if err != nil {
			t.Fatalf("Failed to create Docker client: %v", err)
		}

		containerInspect, err := cli.ContainerInspect(ctx, minioContainer.GetContainerID())
		assert.NoError(t, err)

		instancePort := GetPortFromTheContainer(containerInspect)

		mappedPort, err := minioContainer.MappedPort(ctx, "9000")
		assert.NoError(t, err)

		assert.Equal(t, mappedPort.Port(), instancePort)
	})

	t.Run("It should get the access and secret key from a minio instance", func(t *testing.T) {
		cli, err := client.NewClientWithOpts(client.FromEnv)
		if err != nil {
			t.Fatalf("Failed to create Docker client: %v", err)
		}

		containerInspect, err := cli.ContainerInspect(ctx, minioContainer.GetContainerID())
		assert.NoError(t, err)

		keyMap := GetAccessSecretKeyFromTheContainer(containerInspect)

		assert.Contains(t, keyMap, "MINIO_ACCESS_KEY")
		assert.Contains(t, keyMap, "MINIO_SECRET_KEY")
		assert.Equal(t, keyMap["MINIO_ACCESS_KEY"], "test")
		assert.Equal(t, keyMap["MINIO_SECRET_KEY"], "test")
	})
}