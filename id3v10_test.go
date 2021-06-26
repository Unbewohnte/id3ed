package id3ed

import (
	"os"
	"testing"
)

var TESTv1TAGS = &ID3v1Tags{
	SongName: "testsong",
	Artist:   "testartist",
	Album:    "testalbum",
	Year:     727,
	Comment:  "testcomment",
	Genre:    "Blues",
}

func TestGetID3v1Tags(t *testing.T) {
	testfile, err := os.Open("./testData/testreadv1.mp3")
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
	os.Remove("./testData/testwritev1.mp3")

	f, err := os.Create("./testData/testwritev1.mp3")
	if err != nil {
		t.Errorf("%s", err)
	}
	defer f.Close()

	tags := TESTv1TAGS

	err = WriteID3v1Tags(f, tags)
	if err != nil {
		t.Errorf("WriteID3v1Tags failed: %s", err)
	}

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
	f, err := os.Open("./testData/testwritev1.mp3")
	if err != nil {
		t.Errorf("%s", err)
	}

	tags := TESTv1TAGS

	err = WriteID3v1ToFile(f, tags)
	if err != nil {
		t.Errorf("WriteID3v1ToFile failed: %s", err)
	}

}
