package utils

import (
	"log"
	"net/http"
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Yellow = "\033[33m"
	Green  = "\033[32m"
)

type Error struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

type Success struct {
	Success bool        `json:"success"`
	Message interface{} `json:"message"`
}

func SuccessResponse(statusCode int, resonseJSON, method, path string, w http.ResponseWriter) {
	log.Printf("%s[%s] %s StatusCode:%d%s", Green, method, path, statusCode, Reset)
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(resonseJSON))
}

func ClientErrorResponse(statusCode int, resonseJSON, method, path string, w http.ResponseWriter) {
	log.Printf("%s[%s] %s StatusCode:%d%s", Yellow, method, path, statusCode, Reset)
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(resonseJSON))
}

func ServerErrorResponse(statusCode int, resonse, method, path string, w http.ResponseWriter) {
	log.Printf("%s[%s] %s StatusCode:%d%s", Red, method, path, statusCode, Reset)
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "plain/text")
	w.Write([]byte(resonse))
}
