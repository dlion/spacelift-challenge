package main

import (
	"log"
	"net/http"
	"os"

	"github.com/dlion/spacelift-challenge/docker"
	"github.com/dlion/spacelift-challenge/handlers"
	"github.com/docker/docker/client"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatalf("Failed to create Docker client: %v", err)
	}
	prova := handlers.HandlerManager{DockerCli: cli}

	r.HandleFunc("/object/{id}", prova.UploadHandler).Methods("PUT")
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
