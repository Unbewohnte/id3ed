package id3ed

import (
	"os"
	"testing"
)

func TestGetID3v1Tags(t *testing.T) {
	testfile, err := os.Open("./testData/testmp3id3v1.mp3")
	if err != nil {
		t.Errorf("could not open file for testing: %s", err)
	}
	tags, err := GetID3v1Tags(testfile)
	if err != nil {
		t.Errorf("GetID3v1Tags failed: %s", err)
	}

	comment := tags.GetComment()

	if comment != "Comment here " {
		t.Errorf("GetID3v1Tags failed: expected %s; got %s", "Comment here ", comment)
	}
}
