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

func TestSetBit(t *testing.T) {
	var testByte byte = 0

	if SetBit(testByte, 1) != 1 {
		t.Errorf("SetBit failed: expected output %d; got %d", 1, SetBit(testByte, 1))
	}

	if SetBit(testByte, 8) != 255 {
		t.Errorf("SetBit failed: expected output %d; got %d", 255, SetBit(testByte, 8))
	}
}
