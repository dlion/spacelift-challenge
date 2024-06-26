package storage

import (
	"context"
	"fmt"
	"strconv"
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
			ip := getContainerIP(container.NetworkSettings)
			port := getContainerPort(container.Ports)

			containerInspect, err := client.ContainerInspect(context.Background(), container.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to inspect container %s: %v", container.ID, err)
			}

			accessKey, secretKey := getAccessSecretKeysFromTheContainer(containerInspect)
			instance := MinioInstance{
				Name:   container.Names[0],
				URL:    "http://" + ip + ":" + port,
				Access: accessKey,
				Secret: secretKey,
			}
			instances = append(instances, instance)
		}
	}

	return instances, nil
}

func getContainerPort(ports []types.Port) string {
	for _, port := range ports {
		if port.Type == "tcp" && port.PublicPort != 0 {
			return strconv.Itoa(int(port.PublicPort))
		}
	}
	return ""
}

func getContainerIP(netSettings *types.SummaryNetworkSettings) string {
	network := netSettings.Networks[INSTANCE_NAME]
	if network != nil {
		return network.IPAddress
	}

	return netSettings.Networks["bridge"].IPAddress
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
