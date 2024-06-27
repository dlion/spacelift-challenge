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
	r := mux.NewRouter()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatalf("Failed to create Docker client: %v", err)
	}

	minioInstances, err := container.GetMinioInstancesFromDocker(cli)
	if err != nil {
		log.Fatalf("Failed to get Minio Instances: %v", err)
	}

	services := make([]*storage.MinioService, len(minioInstances))
	for i, v := range minioInstances {
		client, err := storage.NewMinioClient(v.URL, v.Access, v.Secret)
		if err != nil {
			log.Fatalf("Failed to initialize Minio Clients: %v", err)
		}
		services[i] = storage.NewMinioService(client)
	}

	handlerManager := handlers.HandlerManager{
		DockerCli:     cli,
		Instances:     len(minioInstances),
		MinioServices: services,
	}
	r.HandleFunc("/object/{id}", handlerManager.UploadHandler).Methods("PUT")
	r.HandleFunc("/object/{id}", handlers.GetHandler).Methods("GET")

	serverAddress := getServerAddress(cli)
	log.Printf("Starting server on %s\n", serverAddress)
	log.Fatal(http.ListenAndServe(serverAddress, r))
}

func getServerAddress(cli *client.Client) string {
	containerInsp, err := container.InspectContainerByID(cli, os.Getenv("HOSTNAME"))
	if err != nil {
		log.Fatalf("Can't inspect the gateway container from docker")
	}
	serverAddress := container.GetIPAddressFromTheContainer(containerInsp) + ":3000"
	return serverAddress
}
