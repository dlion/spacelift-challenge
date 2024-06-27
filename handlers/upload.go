package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dlion/spacelift-challenge/container"
	"github.com/docker/docker/client"
	"github.com/gorilla/mux"
)

type HandlerManager struct {
	DockerCli *client.Client
}

func (h *HandlerManager) UploadHandler(w http.ResponseWriter, r *http.Request) {
	// @TODO: Use istances for the hash
	instances, err := container.GetMinioInstancesFromDocker(h.DockerCli)
	if err != nil {
		log.Fatalf("Can't get minio instances from docker")
	}

	for _, v := range instances {
		fmt.Println(v)
	}

	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Upload: %v\n", vars["id"])
}
