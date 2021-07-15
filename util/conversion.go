package util

import (
	"strings"

	euni "golang.org/x/text/encoding/unicode"
)

// got the logic from: https://github.com/bogem/id3v2 , thank you very much.

const first7BitsMask = uint32(254) << 24 // shifting 11111110 to the end of uint32

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

// The exact opposite of what `BytesToIntIgnoreFirstBit` does
func SynchsafeIntToBytes(gInt uint32) []byte {
	bytes := make([]byte, 32)

	// looping 4 times (32 bits / 8 bits (4 bytes in int32))
	for i := 0; i < 32; i += 8 {
		gIntCopy := gInt                    //11010101 11001011 00100000 10111111
		first7 := gIntCopy & first7BitsMask //11010100 00000000 00000000 00000000
		shifted := first7 >> 25             //00000000 00000000 00000000 01101010
		bytes = append(bytes, byte(shifted))
	}

	return bytes
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
