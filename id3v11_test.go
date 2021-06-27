package id3ed

import (
	"fmt"
	"os"
	"testing"
)

var TESTV11TAGS = &ID3v11Tags{
	SongName: "testsong",
	Artist:   "testartist",
	Album:    "testalbum",
	Year:     727,
	Comment:  "testcomment",
	Track:    5,
	Genre:    "Blues",
}

func TestGetID3v11Tags(t *testing.T) {
	testfile, err := os.Open("./testData/testreadv1.mp3")
	if err != nil {
		t.Errorf("could not open file for testing: %s", err)
	}
	mp3tags, err := GetID3v11Tags(testfile)
	if err != nil {
		t.Errorf("GetID3v11Tags failed: %s", err)
	}

	if mp3tags.Artist != "Artist" {
		fmt.Printf("%v", mp3tags.Artist)
		t.Errorf("GetID3v11Tags failed:  expected %s; got %s", "Artist", mp3tags.Artist)
	}

	if mp3tags.Track != 4 {
		t.Errorf("GetID3v11Tags failed: expected %d; got %d", 4, mp3tags.Track)
	}
}

func TestWriteID3v11Tags(t *testing.T) {
	os.Remove("./testData/testwritev1.mp3")

	f, err := os.Create("./testData/testwritev1.mp3")
	if err != nil {
		t.Errorf("%s", err)
	}
	defer f.Close()

	tags := TESTV11TAGS

	err = tags.Write(f)
	if err != nil {
		t.Errorf("WriteID3v1Tags failed: %s", err)
	}

	readTags, err := GetID3v11Tags(f)
	if err != nil {
		t.Errorf("%s", err)
	}

	if readTags.Album != "testalbum" {
		t.Errorf("WriteID3v11Tags failed: expected %s; got %s", "testalbum", readTags.Album)
	}

	if readTags.Year != 727 {
		t.Errorf("WriteID3v11Tags failed: expected %d; got %d", 727, readTags.Year)
	}

	if readTags.Track != 5 {
		t.Errorf("WriteID3v11Tags failed: expected %d; got %d", 5, readTags.Track)
	}
}

func TestWriteID3v11ToFile(t *testing.T) {
	f, err := os.Open("./testData/testwritev1.mp3")
	if err != nil {
		t.Errorf("%s", err)
	}

	tags := TESTV11TAGS

	err = tags.WriteToFile(f)
	if err != nil {
		t.Errorf("WriteID3v1ToFile failed: %s", err)
	}
}
