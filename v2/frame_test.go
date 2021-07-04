package v2

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/Unbewohnte/id3ed/util"
)

func TestReadFrame(t *testing.T) {
	f, err := os.Open(filepath.Join(TESTDATAPATH, "testreadv2.mp3"))
	if err != nil {
		t.Errorf("%s", err)
	}

	// read right after header`s bytes
	f.Seek(int64(HEADERSIZE), io.SeekStart)

	firstFrame, err := ReadFrame(f)
	if err != nil {
		t.Errorf("ReadFrame failed: %s", err)
	}

	if firstFrame.ID != "TRCK" {
		t.Errorf("ReadFrame failed: expected ID %s; got %s", "TRCK", firstFrame.ID)
	}

	if firstFrame.Flags.Encrypted != false {
		t.Errorf("ReadFrame failed: expected compressed flag to be %v; got %v", false, firstFrame.Flags.Encrypted)
	}

	secondFrame, err := ReadFrame(f)
	if err != nil {
		t.Errorf("ReadFrame failed: %s", err)
	}

	if secondFrame.ID != "TDRC" {
		t.Errorf("ReadFrame failed: expected ID %s; got %s", "TDRC", secondFrame.ID)
	}

	if util.ToString(secondFrame.Contents) != "2006" {
		t.Errorf("ReadFrame failed: expected contents to be %s; got %s", "2006", util.ToString(secondFrame.Contents))
	}
}
