package util

import (
	"fmt"
	"strconv"
	"strings"
)

// Decodes given byte into integer
func ByteToInt(gByte byte) (int64, error) {
	integer, err := strconv.ParseInt(fmt.Sprintf("%d", gByte), 10, 64)
	if err != nil {
		return 0, err
	}
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

// Converts given bytes into string, ignoring the first 31 non-printable ASCII characters.
// (LOSSY, if given bytes contain some nasty ones)
func ToString(gBytes []byte) string {
	var filteredBytes []byte
	for _, b := range gBytes {
		if b <= 31 {
			continue
		}
		filteredBytes = append(filteredBytes, b)
	}

	return strings.ToValidUTF8(string(filteredBytes), "")
}
