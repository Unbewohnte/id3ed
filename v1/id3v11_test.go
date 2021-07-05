package v1

import (
	"fmt"
	"os"
	"path/filepath"
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
	testfile, err := os.Open(filepath.Join(TESTDATAPATH, "testreadv1.mp3"))
	if err != nil {
		t.Errorf("could not open file for testing: %s", err)
	}
	mp3tags, err := Getv11Tags(testfile)
	if err != nil {
		t.Errorf("GetID3v11Tags failed: %s", err)
	}

	if mp3tags.Artist != "Artist" {
		fmt.Printf("%v", mp3tags.Artist)
		t.Errorf("GetID3v11Tags failed:  expected artist %s; got %s", "Artist", mp3tags.Artist)
	}

	if mp3tags.Track != 8 {
		t.Errorf("GetID3v11Tags failed: expected track %d; got %d", 8, mp3tags.Track)
	}
}

// WILL ADD NEW "TAG" WITHOUT REMOVING THE OLD ONE !!!
func TestWriteID3v11Tags(t *testing.T) {
	f, err := os.OpenFile(filepath.Join(TESTDATAPATH, "testwritev1.mp3"), os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		t.Errorf("%s", err)
	}
	defer f.Close()

	tags := TESTV11TAGS

	err = tags.Write(f)
	if err != nil {
		t.Errorf("WriteID3v1Tags failed: %s", err)
	}

	readTags, err := Getv11Tags(f)
	if err != nil {
		t.Errorf("%s", err)
	}

	if readTags.Album != "testalbum" {
		t.Errorf("WriteID3v11Tags failed: expected album %s; got %s", "testalbum", readTags.Album)
	}

	if readTags.Year != 727 {
		t.Errorf("WriteID3v11Tags failed: expected year %d; got %d", 727, readTags.Year)
	}

	if readTags.Track != 10 {
		t.Errorf("WriteID3v11Tags failed: expected track %d; got %d", 10, readTags.Track)
	}
}

func TestWriteID3v11ToFile(t *testing.T) {
	f, err := os.OpenFile(filepath.Join(TESTDATAPATH, "testwritev1.mp3"), os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		t.Errorf("%s", err)
	}

	tags := TESTV11TAGS

	err = tags.WriteToFile(f)
	if err != nil {
		t.Errorf("WriteID3v1ToFile failed: %s", err)
	}
}
