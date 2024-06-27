package main

import (
	"log"
	"net/http"
	"os"

	"github.com/dlion/spacelift-challenge/docker"
	"github.com/dlion/spacelift-challenge/handlers"
	"github.com/dlion/spacelift-challenge/storage"
	"github.com/docker/docker/client"
	"github.com/gorilla/mux"
)

func main() {

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatalf("Failed to create Docker client: %v", err)
	}

	//@TODO: Use istances for the hash
	instances, err := storage.GetMinioInstancesFromDocker(cli)
	if err != nil {
		log.Fatalf("Can't get minio instances from docker")
	}

	for _, v := range instances {
		log.Println(v)
	}

	r := mux.NewRouter()
	r.HandleFunc("/object/{id}", handlers.UploadHandler).Methods("PUT")
	r.HandleFunc("/object/{id}", handlers.GetHandler).Methods("GET")

	serverAddress := getServerAddress(cli)
	log.Printf("Starting server on %s\n", serverAddress)
	log.Fatal(http.ListenAndServe(serverAddress, r))
}

func getServerAddress(cli *client.Client) string {
	containerInsp, err := docker.InspectContainerByID(cli, os.Getenv("HOSTNAME"))
	if err != nil {
		log.Fatalf("Can't inspect the gateway container from docker")
	}
	serverAddress := docker.GetIPAddressFromTheContainer(containerInsp) + ":3000"
	return serverAddress
}
