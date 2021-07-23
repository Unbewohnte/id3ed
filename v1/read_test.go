package v1

import (
	"os"
	"path/filepath"
	"testing"
)

var TESTv1TAG = &ID3v1Tag{
	SongName:       "testsong",
	Artist:         "testartist",
	Album:          "testalbum",
	Year:           727,
	Comment:        "testcomment",
	Genre:          "Blues",
	HasEnhancedTag: true,
	EnhancedTag: EnhancedID3v1Tag{
		Artist:   "ARRRTIST",
		Album:    "ALLLLBUUUM",
		SongName: "NAME",
		Speed:    EnhancedSpeed[2],
	},
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

	if tag.version != V1_1 {
		t.Errorf("GetID3v1Tag failed: expected version to be %s; got %s", V1_1, tag.version)
	}

	if tag.Comment != "testcomment" {
		t.Errorf("GetID3v1Tag failed: expected comment to be %s; got %s", "testcomment", tag.Comment)
	}

	if tag.Genre != id3v1genres[0] {
		t.Errorf("GetID3v1Tag failed: expected genre to be %s; got %s", id3v1genres[0], tag.Genre)
	}

	if tag.Track != 8 {
		t.Errorf("GetID3v1Tag failed: expected track number to be %d; got %d", 8, tag.Track)
	}
}
