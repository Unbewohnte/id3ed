package v2

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Unbewohnte/id3ed/util"
)

func TestReadNextFrame(t *testing.T) {
	f, err := os.Open(filepath.Join(TESTDATAPATH, "testreadv2.mp3"))
	if err != nil {
		t.Errorf("%s", err)
	}

	header, err := readHeader(f)
	if err != nil {
		t.Errorf("%s", err)
	}

	firstFrame, _, err := readNextFrame(f, header)
	if err != nil {
		t.Errorf("ReadFrame failed: %s", err)
	}

	if firstFrame.Header.ID != "TRCK" {
		t.Errorf("GetFrame failed: expected ID %s; got %s",
			"TRCK", firstFrame.Header.ID)
	}

	if firstFrame.Header.Flags.Encrypted != false {
		t.Errorf("ReadFrame failed: expected compressed flag to be %v; got %v",
			false, firstFrame.Header.Flags.Encrypted)
	}

	secondFrame, _, err := readNextFrame(f, header)
	if err != nil {
		t.Errorf("ReadFrame failed: %s", err)
	}

	if secondFrame.Header.ID != "TDRC" {
		t.Errorf("ReadFrame failed: expected ID %s; got %s",
			"TDRC", secondFrame.Header.ID)
	}

	if util.ToStringLossy(secondFrame.Contents) != "2006" {
		t.Errorf("ReadFrame failed: expected contents to be %s; got %s",
			"2006", secondFrame.Contents)
	}
}
