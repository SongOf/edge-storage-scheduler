package handler

import (
	"io"
	"net/http"
)

func HealthHandler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "edge storage agent is alive")
}
