package handlers

import (
	"unicode"

	"github.com/dlion/spacelift-challenge/storage"
)

type HandlerManager struct {
	MinioServices []*storage.MinioService
}

const (
	MAXIMUM_ID_CHARS = 32
)

func NewHandlerManager(services []*storage.MinioService) *HandlerManager {
	return &HandlerManager{MinioServices: services}
}

func isAlphanumeric(s string) bool {
	for _, char := range s {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) {
			return false
		}
	}
	return true
}
