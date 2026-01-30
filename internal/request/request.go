package request

import (
	"errors"
	"fmt"
	"httpserver/internal/headers"
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"
)

type Request struct {
	RequestLine RequestLine
	Body        []byte
	Headers     headers.Headers
	state       requestState
}

type requestState int

const (
	requestStateInitialized requestState = iota
	requestStateParsingHeaders
	requestStateParsingBody
	requestStateDone
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	bufferSize := 8
	request := &Request{
		state:   requestStateInitialized,
		Headers: headers.NewHeaders(),
	}
	buf := make([]byte, bufferSize)
	readToIndex := 0
	for {
		if request.state == requestStateDone {
			break
		}
		if readToIndex >= len(buf) {
			bufferSize *= 2
			newBuf := make([]byte, bufferSize)
			copy(newBuf, buf[:readToIndex])
			buf = newBuf
		}
		n, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if err == io.EOF {
				if request.state == requestStateDone {
					break
				}
				return nil, errors.New("unexpected EOF")
			}
			return nil, err
		}
		readToIndex += n
		bConsumed, err := request.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}
		if bConsumed == 0 {
			continue
		}
		copy(buf, buf[bConsumed:readToIndex])
		readToIndex -= bConsumed
	}

	return request, nil
}

func parseRequestLine(line []byte) (RequestLine, int, error) {
	if !strings.Contains(string(line), "\r\n") {
		return RequestLine{}, 0, nil
	}

	request := strings.SplitN(string(line), "\r\n", 2)
	requestLine := strings.Split(string(request[0]), " ")

	if len(requestLine) != 3 {
		return RequestLine{}, 0, errors.New("invalid request line")
	}

	bytesConsumed := len(request[0]) + 2

	hasAlpha, err := regexp.MatchString(`^[a-zA-Z]+$`, requestLine[0])
	if err != nil {
		return RequestLine{}, 0, err
	}

	if !hasAlpha {
		return RequestLine{}, 0, errors.New("invalid method")
	}

	HTTPVersion := strings.Split(requestLine[2], "/")

	if HTTPVersion[0] != "HTTP" || HTTPVersion[1] != "1.1" {
		return RequestLine{}, 0, errors.New("invalid HTTP version")
	}

	return RequestLine{
		HttpVersion:   HTTPVersion[1],
		RequestTarget: requestLine[1],
		Method:        requestLine[0],
	}, bytesConsumed, nil
}

func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0
	for r.state != requestStateDone {
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return 0, err
		}
		if n == 0 {
			break
		}
		totalBytesParsed += n
	}
	return totalBytesParsed, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.state {
	case requestStateInitialized:
		requestLine, bytesConsumed, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if bytesConsumed == 0 {
			return 0, nil
		}
		r.RequestLine = requestLine
		r.state = requestStateParsingHeaders
		return bytesConsumed, nil

	case requestStateParsingHeaders:
		bytesConsumed, complete, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if bytesConsumed == 0 {
			return 0, nil
		}
		if complete {
			r.state = requestStateParsingBody
		}
		return bytesConsumed, nil

	case requestStateParsingBody:
		dataLen, ok := r.Headers.Get("Content-Length")
		if !ok {
			r.state = requestStateDone
			return 0, nil
		}
		dataLenNum, err := strconv.Atoi(dataLen)
		if err != nil {
			return 0, err
		}
		if len(r.Body) == dataLenNum {
			r.state = requestStateDone
			return 0, nil
		}
		bytesConsumed := len(data)
		if bytesConsumed+len(r.Body) > dataLenNum {
			log.Printf("Data Overflow; Data Expected: %v Received: %v: %s", dataLenNum, bytesConsumed+len(r.Body), data)
			return 0, errors.New("invalid content length")
		}
		r.Body = append(r.Body, data...)
		return bytesConsumed, nil

	case requestStateDone:
		return 0, fmt.Errorf("error: trying to read data in a done state")

	default:
		return 0, errors.New("invalid request state")
	}
}
