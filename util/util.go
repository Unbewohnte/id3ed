package util

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

// Shortcut function to read n bytes from reader. Peeked from here: https://github.com/dhowden/tag/blob/master/util.go
func Read(rs io.Reader, n int) ([]byte, error) {
	read := make([]byte, n)
	_, err := rs.Read(read)
	if err != nil {
		return nil, fmt.Errorf("could not read from reader: %s", err)
	}

	return read, nil
}

// Shortcut function to read n bytes and convert them into string.
// If encountered zero-byte - converts to string only previously read bytes
func ReadToString(rs io.Reader, n int) (string, error) {
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

// Writes data to wr, if len(data) is less than lenNeeded - adds null bytes until written lenNeeded bytes
func WriteToExtent(wr io.Writer, data []byte, lenNeeded int) error {
	if len(data) > lenNeeded {
		return fmt.Errorf("length of given data bytes is bigger than length needed")
	}

	buff := new(bytes.Buffer)
	for i := 0; i < lenNeeded; i++ {
		if i < len(data) {
			err := buff.WriteByte(data[i])
			if err != nil {
				return err
			}
		} else {
			err := buff.WriteByte(0)
			if err != nil {
				return err
			}
		}
	}

	_, err := wr.Write(buff.Bytes())
	if err != nil {
		return err
	}

	return nil
}

// Returns found key (int) in provided map by value (string);
// If key does not exist in map - returns -1
func GetKey(mp map[int]string, givenValue string) int {
	for key, value := range mp {
		if value == givenValue {
			return key
		}
	}
	return -1
}

// Decodes given integer bytes into integer
func BytesToInt(gBytes []byte) (int64, error) {
	buff := bytes.NewBuffer(gBytes)
	integer, err := binary.ReadVarint(buff)
	if err != nil {
		return 0, err
	}
	buff = nil
	return integer, nil
}