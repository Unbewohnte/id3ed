package util

import "testing"

func TestToString(t *testing.T) {
	someVeryNastyBytes := []byte{0, 1, 2, 3, 4, 5, 6, 50, 7, 8, 9, 10, 11, 50, 50}

	gString := ToStringLossy(someVeryNastyBytes)

	if gString != "222" {
		t.Errorf("ToString failed: expected output: %s; got %s", "222", gString)
	}
}
