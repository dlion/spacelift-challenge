package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/dlion/spacelift-challenge/hash"
	"github.com/gorilla/mux"
)

func (h *HandlerManager) UploadHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, exists := vars["id"]
	if !exists || len(id) > MAXIMUM_ID_CHARS || !isAlphanumeric(id) {
		log.Println("Received bad ID")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	instanceNumber := hash.GetInstanceFromKey(id, len(h.MinioServices))

	if err := h.MinioServices[instanceNumber].UploadFile(r.Body, strconv.Itoa(hash.HashId(id))); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
