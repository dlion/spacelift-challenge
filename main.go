package main

import (
	"log"
	"net/http"
	"os"

	"github.com/dlion/spacelift-challenge/container"
	"github.com/dlion/spacelift-challenge/handlers"
	"github.com/dlion/spacelift-challenge/storage"
	"github.com/docker/docker/client"
	"github.com/gorilla/mux"
)

func main() {
	dockerClient := getDockerClientFromEnv()
	minioInstances := getMinioInstancesFromDocker(dockerClient)
	services := getMinioServicesFromMinioInstances(minioInstances)
	router := defineHandlers(services)
	serverAddress := getServerAddress(dockerClient)
	log.Printf("Starting server on %s\n", serverAddress)
	log.Fatal(http.ListenAndServe(serverAddress, router))
}

func defineHandlers(services []*storage.MinioService) *mux.Router {
	handlerManager := handlers.NewHandlerManager(services)
	router := mux.NewRouter()
	router.HandleFunc("/object/{id}", handlerManager.UploadHandler).Methods("PUT")
	router.HandleFunc("/object/{id}", handlerManager.GetHandler).Methods("GET")
	return router
}

func getDockerClientFromEnv() *client.Client {
	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatalf("Failed to create Docker client: %v", err)
	}

	return dockerClient
}

func getMinioServicesFromMinioInstances(minioInstances []container.MinioInstance) []*storage.MinioService {
	services := make([]*storage.MinioService, len(minioInstances))
	for i, v := range minioInstances {
		client, err := storage.NewMinioClient(v.URL, v.Access, v.Secret)
		if err != nil {
			log.Fatalf("Failed to initialize Minio Clients: %v", err)
		}
		services[i] = storage.NewMinioService(client)
	}
	return services
}

func getMinioInstancesFromDocker(dockerClient *client.Client) []container.MinioInstance {
	minioInstances, err := container.GetMinioInstancesFromDocker(dockerClient)
	if err != nil {
		log.Fatalf("Failed to get Minio Instances: %v", err)
	}
	return minioInstances
}

func getServerAddress(cli *client.Client) string {
	containerInsp, err := container.InspectContainerByID(cli, os.Getenv("HOSTNAME"))
	if err != nil {
		log.Fatalf("Can't inspect the gateway container from docker")
	}
	return container.GetIPAddressFromTheContainer(containerInsp) + ":" + container.GetPortFromTheContainer(containerInsp, true)
}
