# Hello, World!

This is a tiny HTTP server written in Go.

It listens on `127.0.0.1:42069`, manually parses incoming HTTP/1.1 requests from a raw TCP connection, and writes responses (status line, headers, and body) back to the client without using Go's `net/http` package.
