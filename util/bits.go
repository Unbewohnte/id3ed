package util

// Tells if bit is set in given byte
func GetBit(n byte, bitN int) bool {
	return n&byte(1<<bitN-1) != 0
}
