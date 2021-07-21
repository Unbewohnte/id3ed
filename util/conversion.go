package util

import (
	"encoding/binary"
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

// Simply converts given uint32 into synch unsafe bytes
func IntToBytes(gInt uint32) []byte {
	buff := make([]byte, 4)
	binary.BigEndian.PutUint32(buff, gInt)
	return buff
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
	var synchsafeIBytes []byte

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
