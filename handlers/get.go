package handlers

import (
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/dlion/spacelift-challenge/hash"
	"github.com/gorilla/mux"
)

func (h *HandlerManager) GetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, exists := vars["id"]
	if !exists || len(id) > MAXIMUM_ID_CHARS || !isAlphanumeric(id) {
		log.Println("Received bad ID")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	instanceNumber := hash.GetInstanceFromKey(id, len(h.MinioServices))

	filename := strconv.Itoa(hash.HashId(id))
	fileBody, err := h.MinioServices[instanceNumber].GetFile(filename)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if _, err := io.Copy(w, fileBody); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
