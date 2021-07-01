package util

import (
	"bytes"
	"testing"
)

func TestRead(t *testing.T) {
	testBuffer := new(bytes.Buffer)

	_, err := testBuffer.WriteString("Writing some data here. ᗜˬᗜ")
	if err != nil {
		t.Errorf("%s", err)
	}

	buffLen := testBuffer.Len()
	buffBytes := testBuffer.Bytes()

	read, err := Read(testBuffer, uint64(buffLen))
	if err != nil {
		t.Errorf("Read failed: %s", err)
	}

	if int(len(read)) != int(buffLen) {
		t.Errorf("Read failed: expected length %d; got %d", buffLen, len(read))
	}

	for index, b := range read {
		if buffBytes[index] != b {
			t.Errorf("Read failed: expected byte %v; got %v at index %d", buffBytes[index], b, index)
		}
	}
}

func TestReadToString(t *testing.T) {
	testBuffer := new(bytes.Buffer)

	_, err := testBuffer.WriteString("This is literally the same as before. ᗜˬᗜ")
	if err != nil {
		t.Errorf("%s", err)
	}

	buffLen := testBuffer.Len()
	buffBytes := testBuffer.Bytes()

	read, err := Read(testBuffer, uint64(buffLen))
	if err != nil {
		t.Errorf("Read failed: %s", err)
	}

	if int(len(read)) != int(buffLen) {
		t.Errorf("Read failed: expected length %d; got %d", buffLen, len(read))
	}

	for index, b := range read {
		if buffBytes[index] != b {
			t.Errorf("Read failed: expected byte %v; got %v at index %d", buffBytes[index], b, index)
		}
	}
}
