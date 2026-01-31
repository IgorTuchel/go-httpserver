package response

import (
	"httpserver/internal/headers"
	"io"
	"strconv"
)

type Writer struct {
	io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{w}
}

type StatusCode int

const (
	OK                    StatusCode = 200
	BAD_REQUEST           StatusCode = 400
	NOT_FOUND             StatusCode = 404
	INTERNAL_SERVER_ERROR StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	switch statusCode {
	case OK:
		_, err := w.Write([]byte("HTTP/1.1 200 OK\r\n"))
		if err != nil {
			return err
		}
		return nil
	case BAD_REQUEST:
		_, err := w.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
		if err != nil {
			return err
		}
		return nil
	case NOT_FOUND:
		_, err := w.Write([]byte("HTTP/1.1 404 Not Found\r\n"))
		if err != nil {
			return err
		}
		return nil
	case INTERNAL_SERVER_ERROR:
		_, err := w.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
		if err != nil {
			return err
		}
		return nil
	default:
		_, err := w.Write([]byte("HTTP/1.1 " + strconv.Itoa(int(statusCode)) + " Unknown Status Code\r\n"))
		if err != nil {
			return err
		}
		return nil
	}
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	return headers.Headers{
		"Content-Length": strconv.Itoa(contentLen),
		"Content-Type":   "text/plain",
		"Connection":     "close",
	}
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for key, value := range headers {
		_, err := w.Write([]byte(key + ": " + value + "\r\n"))
		if err != nil {
			return err
		}
	}
	_, err := w.Write([]byte("\r\n"))
	if err != nil {
		return err
	}
	return nil
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	return WriteStatusLine(w, statusCode)
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	return WriteHeaders(w, headers)
}

func (w *Writer) WriteBody(body []byte) (int, error) {
	n, err := w.Write(body)
	if err != nil {
		return n, err
	}
	return n, nil
}

func (w *Writer) WriteTrailers(h headers.Headers) error {
	for key, value := range h {
		_, err := w.Write([]byte(key + ": " + value + "\r\n"))
		if err != nil {
			return err
		}
	}
	_, err := w.Write([]byte("\r\n"))
	if err != nil {
		return err
	}
	return nil
}
