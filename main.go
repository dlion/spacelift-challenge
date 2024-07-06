package main

import (
	"hash/fnv"
	"log"
	"net/http"
	"os"

	"github.com/buraksezer/consistent"
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
	consistent := addMinioInstancesToTheRing(minioInstances)
	router := defineHandlers(services, consistent)
	serverAddress := getServerAddress(dockerClient)
	log.Printf("Starting server on %s\n", serverAddress)
	log.Fatal(http.ListenAndServe(serverAddress, router))
}

type hasher struct{}

func (h hasher) Sum64(data []byte) uint64 {
	fh := fnv.New64()
	fh.Write(data)
	return fh.Sum64()
}

func addMinioInstancesToTheRing(minioInstances []container.MinioInstance) *consistent.Consistent {
	cfg := consistent.Config{
		PartitionCount:    len(minioInstances),
		ReplicationFactor: 0,
		Load:              1.25,
		Hasher:            hasher{},
	}

	c := consistent.New(nil, cfg)
	for _, node := range minioInstances {
		c.Add(node)
	}

	return c
}

func defineHandlers(services map[string]*storage.MinioService, consistent *consistent.Consistent) *mux.Router {
	handlerManager := handlers.NewHandlerManager(services, consistent)
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

func getMinioServicesFromMinioInstances(minioInstances []container.MinioInstance) map[string]*storage.MinioService {
	services := make(map[string]*storage.MinioService, len(minioInstances))
	for _, v := range minioInstances {
		client, err := storage.NewMinioClient(v.URL, v.Access, v.Secret)
		if err != nil {
			log.Fatalf("Failed to initialize Minio Clients: %v", err)
		}
		services[v.String()] = storage.NewMinioService(client)
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
