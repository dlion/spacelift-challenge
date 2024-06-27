package handlers

import (
	"io"
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

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Received bad file")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if ok, _ := h.MinioServices[instanceNumber].UploadFile(body, strconv.Itoa(hash.HashId(id))); !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
