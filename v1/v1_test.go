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

func TestReadv1Tag(t *testing.T) {
	testfile, err := os.OpenFile(filepath.Join(TESTDATAPATH, "testreadv1.mp3"), os.O_CREATE|os.O_RDONLY, os.ModePerm)
	if err != nil {
		t.Errorf("could not open file for testing: %s", err)
	}
	tag, err := Readv1Tag(testfile)
	if err != nil {
		t.Errorf("GetID3v1Tag failed: %s", err)
	}

	if tag.Version != V1_1 {
		t.Errorf("GetID3v1Tag failed: expected version to be %s; got %s", V1_1, tag.Version)
	}

	if tag.Comment != "Comment here " {
		t.Errorf("GetID3v1Tag failed: expected comment to be %s; got %s", "Comment here ", tag.Comment)
	}

	if tag.Genre != "Soundtrack" {
		t.Errorf("GetID3v1Tag failed: expected genre to be %s; got %s", "Soundtrack", tag.Genre)
	}

	if tag.Track != 8 {
		t.Errorf("GetID3v1Tag failed: expected track number to be %d; got %d", 8, tag.Track)
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
	err = tag.write(f)
	if err != nil {
		t.Errorf("WriteID3v1Tag failed: %s", err)
	}

	// reading a tag
	readTag, err := Readv1Tag(f)
	if err != nil {
		t.Errorf("%s", err)
	}

	if readTag.Album != "testalbum" {
		t.Errorf("WriteID3v1Tag failed: expected %s; got %s", "testalbum", readTag.Album)
	}

	if readTag.Year != 727 {
		t.Errorf("WriteID3v1Tag failed: expected %d; got %d", 727, readTag.Year)
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
