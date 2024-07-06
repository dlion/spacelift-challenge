package handlers

import (
	"unicode"

	"github.com/buraksezer/consistent"
	"github.com/dlion/spacelift-challenge/storage"
)

type HandlerManager struct {
	MinioServices map[string]*storage.MinioService
	Consistent    *consistent.Consistent
}

const (
	MAXIMUM_ID_CHARS = 32
)

func NewHandlerManager(services map[string]*storage.MinioService, consistent *consistent.Consistent) *HandlerManager {
	return &HandlerManager{MinioServices: services, Consistent: consistent}
}

func isAlphanumeric(s string) bool {
	for _, char := range s {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) {
			return false
		}
	}
	return true
}
