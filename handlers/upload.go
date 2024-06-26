package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Upload: %v\n", vars["id"])
}
