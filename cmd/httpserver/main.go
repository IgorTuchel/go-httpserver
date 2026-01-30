package main

import (
	"httpserver/internal/request"
	"httpserver/internal/response"
	"httpserver/internal/server"
	"log"
	"os"
	"os/signal"
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
	headers := response.GetDefaultHeaders(len(message), "text/html")
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
	headers := response.GetDefaultHeaders(len(message), "text/html")
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
	headers := response.GetDefaultHeaders(len(message), "text/html")
	w.WriteHeaders(headers)
	w.WriteBody(message)
}
