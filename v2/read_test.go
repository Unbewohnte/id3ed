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

	_, err = ReadV2Tag(f)
	if err != nil {
		t.Errorf("GetV2Tag failed: %s", err)
	}
}
