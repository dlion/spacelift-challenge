package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func (h *HandlerManager) UploadHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, exists := vars["id"]
	if checkID(id, exists) {
		log.Println("Received bad ID")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	instanceID := h.Consistent.LocateKey([]byte(id))
	if err := h.MinioServices[instanceID.String()].UploadFile(r.Body, id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
