package util

import (
	"testing"
)

func TestToStringLossy(t *testing.T) {
	someVeryNastyBytes := []byte{0, 1, 2, 3, 4, 5, 6, 50, 7, 8, 9, 10, 11, 50, 50}

	gString := ToStringLossy(someVeryNastyBytes)

	if gString != "222" {
		t.Errorf("ToString failed: expected output: %s; got %s", "222", gString)
	}
}

func TestDecodeText(t *testing.T) {
	// 3 - UTF-8 encoding
	someFrameContents := []byte{3, 50, 50, 50, 50, 0, 0, 0, 0, 50}

	decodedUtf8text := DecodeText(someFrameContents)

	if decodedUtf8text != "22222" {
		t.Errorf("DecodeText failed: expected text %s, got %s", "22222", decodedUtf8text)
	}
}

// func TestIntToBytesFirstBitZeroed(t *testing.T) {
// 	var testint uint32 = 123456

// 	intbytes := IntToBytesFirstBitZeroed(testint)

// 	if BytesToIntIgnoreFirstBit(intbytes) != testint {
// 		t.Errorf("IntToBytesFirstBitZeroed failed: expected to get %v; got %v",
// 			testint, BytesToIntIgnoreFirstBit(intbytes))
// 	}
// }
