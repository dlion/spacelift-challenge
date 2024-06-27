package container

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func GetContainerList(client *client.Client, ctx context.Context) ([]types.Container, error) {
	containers, err := client.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		return nil, err
	}

	return containers, nil
}

func InspectContainerByID(client *client.Client, id string) (types.ContainerJSON, error) {
	containerInspect, err := client.ContainerInspect(context.Background(), id)
	if err != nil {
		return types.ContainerJSON{}, fmt.Errorf("failed to inspect container %s: %v", id, err)
	}

	return containerInspect, nil
}

func GetIPAddressFromTheContainer(ci types.ContainerJSON) string {
	if ci.NetworkSettings != nil {
		for _, network := range ci.NetworkSettings.Networks {
			return network.IPAddress
		}
	}
	return "0.0.0.0"
}

func GetPortFromTheContainer(ci types.ContainerJSON, evaluateBindings bool) string {
	for portMapping, bindings := range ci.NetworkSettings.Ports {
		if evaluateBindings && len(bindings) > 0 {
			return bindings[0].HostPort
		}

		if !evaluateBindings {
			return portMapping.Port()
		}
	}

	return ""
}

func GetAccessSecretKeyFromTheContainer(ci types.ContainerJSON) map[string]string {
	envMap := make(map[string]string)

	for _, env := range ci.Config.Env {
		pair := strings.SplitN(env, "=", 2)
		if len(pair) == 2 {
			envMap[pair[0]] = pair[1]
		}
	}

	return envMap
}
