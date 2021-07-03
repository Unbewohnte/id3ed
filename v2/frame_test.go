package v2

import (
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestReadFrame(t *testing.T) {
	f, err := os.Open(filepath.Join(TESTDATAPATH, "testreadv2.mp3"))
	if err != nil {
		t.Errorf("%s", err)
	}

	// read right after header`s bytes
	f.Seek(int64(HEADERSIZE), io.SeekStart)

	firstFrame, err := ReadFrame(f, 24)
	if err != nil {
		t.Errorf("ReadFrame failed: %s", err)
	}

	if firstFrame.Header.ID != "TRCK" {
		t.Errorf("ReadFrame failed: expected ID %s; got %s", "TRCK", firstFrame.Header.ID)
	}
}
