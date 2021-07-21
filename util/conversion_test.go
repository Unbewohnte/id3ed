package util

import (
	"fmt"
	"strconv"
	"testing"
)

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
