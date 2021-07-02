package util

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
)

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
		return 0, fmt.Errorf("could not decode integer: %s", err)
	}
	buff = nil
	return integer, nil
}

// Decodes given integer bytes into integer, ignores the first bit
// of every given byte in binary form
func BytesToIntIgnoreFirstBit(gBytes []byte) (int64, error) {
	// represent each byte in size as binary and get rid from the first bit,
	// then concatenate filtered parts
	var filteredBits string
	for _, b := range gBytes {
		// ignore the first bit
		filteredPart := fmt.Sprintf("%08b", b)[1:] // byte is 8 bits
		filteredBits += filteredPart
	}

	// convert filtered binary into usable int64
	integer, err := strconv.ParseInt(filteredBits, 2, 64)
	if err != nil {
		return -1, err
	}

	return integer, nil
}
