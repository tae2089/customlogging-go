package main

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
)

type CustomWriter struct{}

func (e CustomWriter) Write(p []byte) (int, error) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, p, "", "    "); err != nil {
		return 0, err
	}
	n, err := os.Stdout.Write(prettyJSON.Bytes())
	if err != nil {
		return n, err
	}
	if n != len(prettyJSON.Bytes()) {
		return n, io.ErrShortWrite
	}
	return len(p), nil
}
