package util

import (
	"strings"

	euni "golang.org/x/text/encoding/unicode"
)

// got the logic from: https://github.com/bogem/id3v2 , thank you very much.
const first7BitsMask = uint32(254 << 24) // shifting 11111110 to the end of uint32

// Converts given bytes into integer
func BytesToInt(gBytes []byte) uint32 {
	var integer uint32 = 0
	for _, b := range gBytes {
		integer = integer << 8
		integer = integer | uint32(b)
	}
	return integer
}

// Decodes given integer bytes into integer, ignores the first bit
// of every given byte in binary form
func BytesToIntSynchsafe(gBytes []byte) uint32 {
	var integer uint32 = 0
	for _, b := range gBytes {
		integer = integer << 7
		integer = integer | uint32(b)
	}

	return integer
}

// The exact opposite of what `BytesToIntSynchsafe` does
// Finally understood with the help of: https://github.com/bogem/id3v2/blob/master/size.go ,
// thank you very much !
func IntToBytesSynchsafe(gInt uint32) []byte {
	synchsafeIBytes := make([]byte, 4)

	// skip 4 0`ed bits
	gInt = gInt << 4

	// int32 == 4 bytes
	for i := 0; i < 32/8; i++ {
		// get first 7 bits
		first7Bits := gInt & first7BitsMask

		// shift captured bits to the beginning
		first7Bits = first7Bits >> (3*8 + 1)

		b := byte(first7Bits)
		synchsafeIBytes = append(synchsafeIBytes, b)

		// prepare next 7 bits for the next iteration
		gInt = gInt << 7
	}
	return synchsafeIBytes
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
