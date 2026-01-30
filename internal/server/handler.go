package server

import (
	"httpserver/internal/request"
	"httpserver/internal/response"
	"io"
)

type HandlerError struct {
	Message string
	Code    response.StatusCode
}

func NewHandlerError(w io.Writer, err HandlerError) {
	response.WriteStatusLine(w, err.Code)
	headers := response.GetDefaultHeaders(len(err.Message), "text/plain")
	response.WriteHeaders(w, headers)
	w.Write([]byte(err.Message))

}

type Handler func(w *response.Writer, req *request.Request)
