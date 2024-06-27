package storage

import (
	"context"
	"fmt"
	"strings"

	"github.com/dlion/spacelift-challenge/docker"
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
				Access: mapKey["MINIO_ACCESS_KEY"],
				Secret: mapKey["MINIO_SECRET_KEY"],
			})
		}
	}

	return instances, nil
}

func minioInstanceExist(container types.Container) bool {
	name, ok := container.Labels["com.docker.compose.service"]
	return ok && strings.Contains(name, fmt.Sprintf("%s-node", INSTANCE_NAME))
}
