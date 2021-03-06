package v2

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadV2Tag(t *testing.T) {
	f, err := os.Open(filepath.Join(TESTDATAPATH, "testreadv2.mp3"))
	if err != nil {
		t.Errorf("%s", err)
	}

	tag, err := ReadV2Tag(f)
	if err != nil {
		t.Errorf("GetV2Tag failed: %s", err)
	}

	if tag.Padding != 1024 {
		t.Errorf("GetV2Tag failed: expected to have %d padding bytes: got %d", 1024, tag.Padding)
	}

	titleFrame := tag.GetFrame("TIT2")

	if titleFrame.Text() != "title" {
		t.Errorf("ReadV2Tag failed: expected contents of the title frame to be %s; got %s",
			"title", titleFrame.Text())
	}

	album := tag.Album()
	if album != "album" {
		t.Errorf("ReadV2Tag failed: expected contents of the album frame to be %s; got %s",
			"album", album)
	}

	picture := tag.Picture()
	if picture != nil {
		t.Errorf("ReadV2Tag failed: expected file not to have a picture")
	}

	genre := tag.GetFrame("TCON")
	if genre == nil {
		t.Errorf("ReadV2Tag failed: expected genre to be %s; got %v", "anime", genre)
	}
}
