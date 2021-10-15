package handler

import (
	"io"
	"net/http"
)

func DefaultHandler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello world")
}
