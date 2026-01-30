package main

import (
	"httpserver/internal/request"
	"httpserver/internal/response"
	"httpserver/internal/server"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handlerYProblem)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handlerYProblem(w io.Writer, req *request.Request) *server.HandlerError {
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		return &server.HandlerError{Code: response.BAD_REQUEST, Message: "Your problem is not my problem\n"}
	case "/myproblem":
		return &server.HandlerError{Code: response.INTERNAL_SERVER_ERROR, Message: "Woopsie, my bad\n"}
	default:
		w.Write([]byte("All good, frfr\n"))
		return nil
	}
}
