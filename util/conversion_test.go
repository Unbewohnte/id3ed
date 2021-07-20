package util

import (
	"fmt"
	"strconv"
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
		if len(synchSafeBytes) != 4 {
			t.Errorf("IntToBytesSynchsafe failed: expected len to be %d; got %d", 4, len(synchSafeBytes))
		}

		synchsafeInt := BytesToIntSynchsafe(synchSafeBytes)

		if synchsafeInt != testInt {
			t.Errorf("BytesToIntSynchsafe failed: expected to get %d; got %d", testInt, synchsafeInt)
		}
	}
}

func TestIntToBytes(t *testing.T) {
	var testInt uint32 = 124567
	testIntBits := fmt.Sprintf("%032b", testInt)

	gotBytes := IntToBytes(testInt)

	i := 0
	for _, gotByte := range gotBytes {
		correctByte, _ := strconv.ParseUint(testIntBits[i:i+8], 2, 8)
		if gotByte != byte(correctByte) {
			t.Errorf("IntToBytes failed: expected byte to be %d; got %d", correctByte, gotByte)
		}
		i += 8
	}
}
