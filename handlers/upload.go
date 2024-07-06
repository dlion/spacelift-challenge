package handlers

import (
	"log"
	"net/http"

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

	consistentInstanceNumber := h.Consistent.LocateKey([]byte(id))
	if err := h.MinioServices[consistentInstanceNumber.String()].UploadFile(r.Body, id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
