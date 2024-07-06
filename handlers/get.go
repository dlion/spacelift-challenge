package handlers

import (
	"io"
	"log"
	"net/http"

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

	consistentInstanceNumber := h.Consistent.LocateKey([]byte(id))
	fileBody, err := h.MinioServices[consistentInstanceNumber.String()].GetFile(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if _, err := io.Copy(w, fileBody); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
