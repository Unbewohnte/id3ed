package id3ed

import (
	"fmt"
	"os"
	"testing"
)

func TestGetID3v11Tags(t *testing.T) {
	testfile, err := os.Open("./testData/testmp3id3v1.mp3")
	if err != nil {
		t.Errorf("could not open file for testing: %s", err)
	}
	mp3tags, err := GetID3v11Tags(testfile)
	if err != nil {
		t.Errorf("GetID3v1Tags failed: %s", err)
	}

	if mp3tags.Artist != "Artist" {
		fmt.Printf("%v", mp3tags.Artist)
		t.Errorf("GetID3v1Tags has failed:  expected %v; got %v", "Artist", mp3tags.Artist)
	}

	if mp3tags.Track != 8 {
		t.Errorf("GetID3v1Tags has failed: expected %v; got %v", 8, mp3tags.Track)
	}
}
