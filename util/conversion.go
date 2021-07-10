package util

import (
	"fmt"
	"strconv"
	"strings"

	euni "golang.org/x/text/encoding/unicode"
)

// Decodes given byte into integer
func ByteToInt(gByte byte) (int, error) {
	integer, err := strconv.Atoi(fmt.Sprintf("%d", gByte))
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
func ToStringLossy(gBytes []byte) string {
	var runes []rune
	for _, b := range gBytes {
		if b <= 31 {
			continue
		}
		runes = append(runes, rune(b))
	}

	return strings.ToValidUTF8(string(runes), "")
}

// Decodes the given frame`s contents
func DecodeText(fContents []byte) string {
	textEncoding := fContents[0] // the first byte is the encoding

	switch textEncoding {
	case 0:
		// ISO-8859-1
		return ToStringLossy(fContents[1:])
	case 1:
		// UTF-16 with BOM
		encoding := euni.UTF16(euni.BigEndian, euni.ExpectBOM)
		decoder := encoding.NewDecoder()

		decodedBytes := make([]byte, len(fContents)*2)
		_, _, err := decoder.Transform(decodedBytes, fContents[1:], true)
		if err != nil {
			return ""
		}

		return string(decodedBytes)

	case 2:
		// UTF-16
		encoding := euni.UTF16(euni.BigEndian, euni.IgnoreBOM)
		decoder := encoding.NewDecoder()

		decodedBytes := make([]byte, len(fContents)*2)
		_, _, err := decoder.Transform(decodedBytes, fContents[1:], true)
		if err != nil {
			return ""
		}

		return string(decodedBytes)

	case 3:
		// UTF-8
		return ToStringLossy(fContents[1:])
	}

	return ""
}
