package server

import (
	"bytes"
	"httpserver/internal/request"
	"httpserver/internal/response"
	"log"
	"net"
	"strconv"
	"sync/atomic"
)

type Server struct {
	listener net.Listener
	closed   atomic.Bool
	handler  Handler
}

func Serve(port int, handler Handler) (*Server, error) {
	s := &Server{handler: handler}
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}
	s.listener = listener
	go s.listen()
	return s, nil
}

func (s *Server) Close() error {
	if s.closed.Load() {
		return nil
	}
	s.closed.Store(true)
	s.listener.Close()
	return nil
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	req, err := request.RequestFromReader(conn)
	if err != nil {
		log.Printf("Error reading request: %v", err)
		return
	}

	res := bytes.Buffer{}

	handlerErr := s.handler(&res, req)
	if handlerErr != nil {
		NewHandlerError(conn, *handlerErr)
		return
	}
	err = response.WriteStatusLine(conn, response.OK)
	if err != nil {
		log.Printf("Error writing status line: %v", err)
		return
	}
	headers := response.GetDefaultHeaders(len(res.Bytes()))
	err = response.WriteHeaders(conn, headers)
	if err != nil {
		log.Printf("Error writing headers: %v", err)
		return
	}
	conn.Write(res.Bytes())
}
