package response

import (
	"fmt"
)

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	size := len(p)
	_, err := w.WriteBody([]byte(fmt.Sprintf("%x\r\n", size)))
	if err != nil {
		return 0, err
	}
	_, err = w.WriteBody(p)
	if err != nil {
		return 0, err
	}
	_, err = w.WriteBody([]byte("\r\n"))
	if err != nil {
		return 0, err
	}
	return size, nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	return w.WriteBody([]byte("0\r\n"))
}
