package util

import (
	"fmt"
	"io"
)

// Shortcut function to read n bytes from reader. The general idea peeked from here: https://github.com/dhowden/tag/blob/master/util.go
func Read(r io.Reader, n uint64) ([]byte, error) {
	read := make([]byte, n)
	_, err := r.Read(read)
	if err != nil {
		return nil, fmt.Errorf("could not read from reader: %s", err)
	}

	return read, nil
}

// Reads from rs and conversts read []byte into string, ignoring all non-printable or
// invalid characters.
func ReadToString(r io.Reader, n uint64) (string, error) {
	read := make([]byte, n)
	_, err := r.Read(read)
	if err != nil {
		return "", fmt.Errorf("could not read from reader: %s", err)
	}
	return ToStringLossy(read), nil
}
