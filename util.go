package id3ed

import (
	"fmt"
	"io"
)

// Shortcut function to read n bytes from reader. Peeked from here: https://github.com/dhowden/tag/blob/master/util.go
func read(rs io.Reader, n int) ([]byte, error) {
	read := make([]byte, n)
	_, err := rs.Read(read)
	if err != nil {
		return nil, fmt.Errorf("could not read from reader: %s", err)
	}

	return read, nil
}

// Shortcut function to read n bytes and convert them into string
func readToString(rs io.Reader, n int) (string, error) {
	read := make([]byte, n)
	_, err := rs.Read(read)
	if err != nil {
		return "", fmt.Errorf("could not read from reader: %s", err)
	}
	return string(read), nil
}
