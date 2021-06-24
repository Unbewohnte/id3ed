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

// Shortcut function to read n bytes and convert them into string.
// If encountered zero-byte - converts to string only previously read bytes
func readToString(rs io.Reader, n int) (string, error) {
	read := make([]byte, n)
	_, err := rs.Read(read)
	if err != nil {
		return "", fmt.Errorf("could not read from reader: %s", err)
	}

	var readString string
	for _, b := range read {
		if b == 0 {
			break
		}
		readString += string(b)
	}

	return readString, nil
}
