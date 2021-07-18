package util

// Tells if bit is set in given byte,
// if bitN <= 0 - always returns false
func GetBit(b byte, bitN int) bool {
	if bitN <= 0 {
		return false
	}
	return b&byte(1<<bitN-1) != 0
}

// Sets bit to 1 in provided byte, if bitN <= 0
// returns original b without modifications
func SetBit(b byte, bitN int) byte {
	if bitN <= 0 {
		return b
	}
	return b | byte(1<<byte(bitN)-1)
}
