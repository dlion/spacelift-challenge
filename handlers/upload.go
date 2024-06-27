package handlers

import (
	"io"
	"net/http"
	"strconv"
	"unicode"

	"github.com/dlion/spacelift-challenge/hash"
	"github.com/dlion/spacelift-challenge/storage"
	"github.com/docker/docker/client"
	"github.com/gorilla/mux"
)

type HandlerManager struct {
	DockerCli     *client.Client
	MinioServices []*storage.MinioService
	Instances     int
}

func (h *HandlerManager) UploadHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, exists := vars["id"]
	if !exists || len(id) > 32 || !isAlphanumeric(id) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hashManager := &hash.HashManager{}
	instanceNumber := hashManager.GetInstanceFromKey(id, h.Instances)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if ok, _ := h.MinioServices[instanceNumber].UploadFile(body, strconv.Itoa(hashManager.HashId(id))); !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func isAlphanumeric(s string) bool {
	for _, char := range s {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) {
			return false
		}
	}
	return true
}
