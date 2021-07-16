package util

import "testing"

func TestGetBit(t *testing.T) {
	testBytes := []byte{1 << 0, 1 << 1, 1 << 2, 1 << 3}

	for index, testByte := range testBytes {
		if !GetBit(testByte, index+1) {
			t.Errorf("IsSet failed: expected %dth bit of %d to be %v; got %v",
				index+1, testByte, true, GetBit(testByte, index+1))
		}
	}
}
