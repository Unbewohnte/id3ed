package v1

import (
	"os"
	"path/filepath"
	"testing"
)

var TESTDATAPATH string = filepath.Join("..", "testData")

var TESTv1TAG = &ID3v1Tag{
	SongName: "testsong",
	Artist:   "testartist",
	Album:    "testalbum",
	Year:     727,
	Comment:  "testcomment",
	Genre:    "Blues",
}

func TestGetv1Tags(t *testing.T) {
	testfile, err := os.OpenFile(filepath.Join(TESTDATAPATH, "testreadv1.mp3"), os.O_CREATE|os.O_RDONLY, os.ModePerm)
	if err != nil {
		t.Errorf("could not open file for testing: %s", err)
	}
	tag, err := Getv1Tag(testfile)
	if err != nil {
		t.Errorf("GetID3v1Tag failed: %s", err)
	}

	if tag.Comment != "Comment here " {
		t.Errorf("GetID3v1Tag failed: expected %s; got %s", "Comment here ", tag.Comment)
	}
}

func TestWritev1Tags(t *testing.T) {
	f, err := os.OpenFile(filepath.Join(TESTDATAPATH, "testwritev1.mp3"), os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		t.Errorf("%s", err)
	}
	defer f.Close()

	tag := TESTv1TAG

	// writing a tag
	err = tag.Write(f)
	if err != nil {
		t.Errorf("WriteID3v1Tag failed: %s", err)
	}

	// reading a tag
	readTags, err := Getv1Tag(f)
	if err != nil {
		t.Errorf("%s", err)
	}

	if readTags.Album != "testalbum" {
		t.Errorf("WriteID3v1Tag failed: expected %s; got %s", "testalbum", readTags.Album)
	}

	if readTags.Year != 727 {
		t.Errorf("WriteID3v1Tag failed: expected %d; got %d", 727, readTags.Year)
	}
}

func TestWriteID3v1ToFile(t *testing.T) {
	f, err := os.OpenFile(filepath.Join(TESTDATAPATH, "testwritev1.mp3"), os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		t.Errorf("%s", err)
	}

	tag := TESTv1TAG

	err = tag.WriteToFile(f)
	if err != nil {
		t.Errorf("WriteID3v1ToFile failed: %s", err)
	}

}
