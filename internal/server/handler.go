package server

import (
	"httpserver/internal/request"
	"httpserver/internal/response"
)

type HandlerError struct {
	Message string
	Code    response.StatusCode
}
type Handler func(w *response.Writer, req *request.Request)
