package main

import (
	"crypto/sha256"
	"fmt"
	"httpserver/internal/headers"
	"httpserver/internal/request"
	"httpserver/internal/response"
	"httpserver/internal/server"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handler)
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

func handler(w *response.Writer, req *request.Request) {
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin") {
		handlerProxyHTTPBin(w, req)
		return
	}
	switch req.RequestLine.RequestTarget {
	case "/":
		defaultHandler(w, req)
		return
	case "/yourproblem":
		handlerYourProblem(w, req)
		return
	case "/myproblem":
		handlerMyProblem(w, req)
		return
	case "/video":
		handlerVideo(w, req)
		return
	default:
		w.Write([]byte("All good, frfr\n"))
		return
	}
}

func handlerMyProblem(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.INTERNAL_SERVER_ERROR)
	message := []byte(`
	<html>
		<head>
			<title>500 Internal Server Error</title>
		</head>
		<body>
			<h1>Internal Server Error</h1>
			<p>Okay, you know what? This one is on me.</p>
		</body>
	</html>`)
	headers := response.GetDefaultHeaders(len(message))
	headers.Overwrite("Content-Type", "text/html")
	w.WriteHeaders(headers)
	w.WriteBody(message)
}

func handlerYourProblem(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.BAD_REQUEST)
	message := []byte(`
	<html>
		<head>
			<title>400 Bad Request</title>
		</head>
		<body>
			<h1>Bad Request</h1>
			<p>Your request honestly kinda sucked.</p>
		</body>
	</html>`)
	headers := response.GetDefaultHeaders(len(message))
	headers.Overwrite("Content-Type", "text/html")
	w.WriteHeaders(headers)
	w.WriteBody(message)
}

func defaultHandler(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.NOT_FOUND)
	message := []byte(`
	<html>
		<head>
			<title>404 Not Found</title>
		</head>
		<body>
			<h1>Not Found</h1>
			<p>The requested resource could not be found.</p>
		</body>
	</html>`)
	headers := response.GetDefaultHeaders(len(message))
	headers.Overwrite("Content-Type", "text/html")
	w.WriteHeaders(headers)
	w.WriteBody(message)
}

func handlerProxyHTTPBin(w *response.Writer, req *request.Request) {
	buf := make([]byte, 1024)
	stream := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin")
	resp, err := http.Get("https://httpbin.org" + stream)
	if err != nil {
		w.WriteStatusLine(response.INTERNAL_SERVER_ERROR)
		message := []byte(`
	<html>
		<head>
			<title>500 Internal Server Error</title>
		</head>
		<body>
			<h1>Internal Server Error</h1>
			<p>Okay, you know what? This one is on me.</p>
		</body>
	</html>`)
		headers := response.GetDefaultHeaders(len(message))
		headers.Overwrite("Content-Type", "text/html")
		w.WriteHeaders(headers)
		w.WriteBody(message)
		return
	}
	defer resp.Body.Close()
	w.WriteStatusLine(response.StatusCode(resp.StatusCode))
	header := headers.Headers{
		"Transfer-Encoding": "chunked",
		"Content-Type":      resp.Header.Get("Content-Type"),
	}
	header.Set("Trailer", "X-Content-SHA256")
	header.Set("Trailer", "X-Content-Length")
	w.WriteHeaders(header)
	acc := []byte{}
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			w.WriteChunkedBody(buf[:n])
			acc = append(acc, buf[:n]...)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}
	}
	w.WriteChunkedBodyDone()
	trailerHeaders := headers.NewHeaders()
	trailerHeaders.Set("X-Content-SHA256", fmt.Sprintf("%x", sha256.Sum256(acc)))
	trailerHeaders.Set("X-Content-Length", strconv.Itoa(len(acc)))
	w.WriteTrailers(trailerHeaders)
	log.Println(trailerHeaders)
}

func handlerVideo(w *response.Writer, req *request.Request) {
	w.WriteStatusLine(response.OK)
	header := headers.Headers{
		"Content-Type": "video/mp4",
	}
	vidBuff, err := os.ReadFile("./assets/vim.mp4")
	if err != nil {
		w.WriteStatusLine(response.INTERNAL_SERVER_ERROR)
		message := []byte(`
	<html>
		<head>
			<title>500 Internal Server Error</title>
		</head>
		<body>
			<h1>Internal Server Error</h1>
			<p>Okay, you know what? This one is on me.</p>
		</body>
	</html>`)
		headers := response.GetDefaultHeaders(len(message))
		headers.Overwrite("Content-Type", "text/html")
		w.WriteHeaders(headers)
		w.WriteBody(message)
		return
	}
	header.Set("Content-Length", strconv.Itoa(len(vidBuff)))
	w.WriteHeaders(header)
	w.WriteBody(vidBuff)
}
