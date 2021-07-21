package util

import (
	"strings"
	"unicode"

	euni "golang.org/x/text/encoding/unicode"
)

// Checks if given characters are in ASCII range
func InASCII(chars string) bool {
	for i := 0; i < len(chars); i++ {
		if chars[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
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

const (
	EncodingISO8859 byte = iota
	EncodingUTF16BOM
	EncodingUTF16
	EncodingUTF8
)

// Decodes the given frame`s contents
func DecodeText(fContents []byte) string {
	textEncoding := fContents[0] // the first byte is the encoding

	switch textEncoding {
	case EncodingISO8859:
		// ISO-8859-1
		return ToStringLossy(fContents[1:])
	case EncodingUTF16BOM:
		// UTF-16 with BOM
		encoding := euni.UTF16(euni.BigEndian, euni.ExpectBOM)
		decoder := encoding.NewDecoder()

		decodedBytes := make([]byte, len(fContents)*2)
		_, _, err := decoder.Transform(decodedBytes, fContents[1:], true)
		if err != nil {
			return ""
		}

		return string(decodedBytes)

	case EncodingUTF16:
		// UTF-16
		encoding := euni.UTF16(euni.BigEndian, euni.IgnoreBOM)
		decoder := encoding.NewDecoder()

		decodedBytes := make([]byte, len(fContents)*2)
		_, _, err := decoder.Transform(decodedBytes, fContents[1:], true)
		if err != nil {
			return ""
		}

		return string(decodedBytes)

	case EncodingUTF8:
		// UTF-8
		return ToStringLossy(fContents[1:])
	}

	return ""
}
