package util

import (
	"bytes"
	"testing"
)

func TestWriteToExtent(t *testing.T) {
	testBuff := new(bytes.Buffer)

	testData := []byte("some data here")

	err := WriteToExtent(testBuff, testData, len(testData)+50)
	if err != nil {
		t.Errorf("WriteToExtent failed: %s", err)
	}

	if testBuff.Len() != len(testData)+50 {
		t.Errorf("WriteToExtent failed: expected length %d; got %d", len(testData)+50, testBuff.Len())
	}

	nullByteCounter := 0
	for _, b := range testBuff.Bytes() {
		if b == 0 {
			nullByteCounter++
		}
	}

	if nullByteCounter != 50 {
		t.Errorf("WriteToExtent failed: expected null bytes %d; got %d", 50, nullByteCounter)
	}
}
