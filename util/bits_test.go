package util

import "testing"

func TestIsSet(t *testing.T) {
	testBytes := []byte{1 << 0, 1 << 1, 1 << 2, 1 << 3}

	for index, testByte := range testBytes {
		if !IsSet(testByte, index+1) {
			t.Errorf("IsSet failed: expected %dth bit of %d to be %v; got %v",
				index+1, testByte, true, IsSet(testByte, index+1))
		}
	}
}
