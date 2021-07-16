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

func TestIntToBytesSynchsafe(t *testing.T) {
	testInts := []uint32{
		1234,
		12,
		1,
		0,
		99999,
		87654321,
	}

	for _, testInt := range testInts {
		synchSafeBytes := IntToBytesSynchsafe(testInt)

		synchsafeInt := BytesToIntSynchsafe(synchSafeBytes)

		if synchsafeInt != testInt {
			t.Errorf("BytesToIntSynchsafe failed: expected to get %d; got %d", testInt, synchsafeInt)
		}
	}
}
