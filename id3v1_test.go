package id3ed

import (
	"os"
	"testing"
)

func TestGetID3v1Tags(t *testing.T) {
	testfile, err := os.Open("./testData/id3v1.txt")
	if err != nil {
		t.Errorf("could not open file for testing: %s", err)
	}
	_, err = GetID3v1Tags(testfile)
	if err != nil {
		t.Errorf("GetID3v1Tags failed: %s", err)
	}
}
