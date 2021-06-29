package v1

import (
	"os"
	"path/filepath"
	"testing"
)

var TESTDATAPATH string = filepath.Join("..", "testData")

var TESTv1TAGS = &ID3v1Tags{
	SongName: "testsong",
	Artist:   "testartist",
	Album:    "testalbum",
	Year:     727,
	Comment:  "testcomment",
	Genre:    "Blues",
}

func TestGetID3v1Tags(t *testing.T) {
	testfile, err := os.OpenFile(filepath.Join(TESTDATAPATH, "testreadv1.mp3"), os.O_CREATE|os.O_RDONLY, os.ModePerm)
	if err != nil {
		t.Errorf("could not open file for testing: %s", err)
	}
	tags, err := GetID3v1Tags(testfile)
	if err != nil {
		t.Errorf("GetID3v1Tags failed: %s", err)
	}

	if tags.Comment != "Comment here " {
		t.Errorf("GetID3v1Tags failed: expected %s; got %s", "Comment here ", tags.Comment)
	}
}

func TestWriteID3v1Tags(t *testing.T) {
	f, err := os.OpenFile(filepath.Join(TESTDATAPATH, "testwritev1.mp3"), os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		t.Errorf("%s", err)
	}
	defer f.Close()

	tags := TESTv1TAGS

	// writing tags
	err = tags.Write(f)
	if err != nil {
		t.Errorf("WriteID3v1Tags failed: %s", err)
	}

	// reading tags
	readTags, err := GetID3v1Tags(f)
	if err != nil {
		t.Errorf("%s", err)
	}

	if readTags.Album != "testalbum" {
		t.Errorf("WriteID3v1Tags failed: expected %s; got %s", "testalbum", readTags.Album)
	}

	if readTags.Year != 727 {
		t.Errorf("WriteID3v1Tags failed: expected %d; got %d", 727, readTags.Year)
	}
}

func TestWriteID3v1ToFile(t *testing.T) {
	f, err := os.OpenFile(filepath.Join(TESTDATAPATH, "testwritev1.mp3"), os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		t.Errorf("%s", err)
	}

	tags := TESTv1TAGS

	err = tags.WriteToFile(f)
	if err != nil {
		t.Errorf("WriteID3v1ToFile failed: %s", err)
	}

}
