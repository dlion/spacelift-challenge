package container

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

const (
	INSTANCE_NAME = "amazin-object-storage"
)

type MinioInstance struct {
	Name   string
	URL    string
	Access string
	Secret string
}

func GetMinioInstancesFromDocker(client *client.Client) ([]MinioInstance, error) {
	containers, err := GetContainerList(client, context.Background())
	if err != nil {
		return nil, err
	}

	var instances []MinioInstance
	for _, container := range containers {
		if minioInstanceExist(container) {
			containerInspect, err := InspectContainerByID(client, container.ID)
			if err != nil {
				return nil, err
			}

			mapKey := GetAccessSecretKeyFromTheContainer(containerInspect)
			ip := GetIPAddressFromTheContainer(containerInspect)
			port := GetPortFromTheContainer(containerInspect)

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
