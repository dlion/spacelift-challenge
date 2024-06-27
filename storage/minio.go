package storage

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
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
	containers, err := client.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		return nil, err
	}

	var instances []MinioInstance
	for _, container := range containers {
		if name, ok := container.Labels["com.docker.compose.service"]; ok && strings.Contains(name, fmt.Sprintf("%s-node", INSTANCE_NAME)) {
			containerInspect, err := client.ContainerInspect(context.Background(), container.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to inspect container %s: %v", container.ID, err)
			}

			accessKey, secretKey := getAccessSecretKeysFromTheContainer(containerInspect)
			ip := getIPAddressFromTheContainer(containerInspect)
			port := getPortFromTheContainer(containerInspect)

			instances = append(instances, MinioInstance{
				Name:   container.Names[0],
				URL:    "http://" + ip + ":" + port,
				Access: accessKey,
				Secret: secretKey,
			})
		}
	}

	return instances, nil
}

func getAccessSecretKeysFromTheContainer(ci types.ContainerJSON) (string, string) {
	envMap := make(map[string]string)

	for _, env := range ci.Config.Env {
		pair := strings.SplitN(env, "=", 2)
		if len(pair) == 2 {
			envMap[pair[0]] = pair[1]
		}
	}

	return envMap["MINIO_ACCESS_KEY"], envMap["MINIO_SECRET_KEY"]
}

func getIPAddressFromTheContainer(ci types.ContainerJSON) string {
	if ci.NetworkSettings != nil {
		for _, network := range ci.NetworkSettings.Networks {
			return network.IPAddress
		}
	}
	return "0.0.0.0"
}

func getPortFromTheContainer(ci types.ContainerJSON) string {
	for portMapping, bindings := range ci.NetworkSettings.Ports {
		if len(bindings) > 0 {
			return bindings[0].HostPort
		}
		return portMapping.Port()
	}
	return ""
}
