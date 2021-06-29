package v2

import (
	"os"
	"path/filepath"
	"testing"
)

var TESTDATAPATH string = filepath.Join("..", "testData")

func TestGetHeader(t *testing.T) {
	f, err := os.Open(filepath.Join(TESTDATAPATH, "testreadv2.mp3"))
	if err != nil {
		t.Errorf("%s", err)
	}

	header, err := GetHeader(f)
	if err != nil {
		t.Errorf("GetHeader failed: %s", err)
	}

	if header.Identifier != "ID3" {
		t.Errorf("GetHeader failed: expected identifier %s; got %s", "ID3", header.Identifier)
	}

	if header.Compressed != false {
		t.Errorf("GetHeader failed: expected flag %v; got %v", false, header.Compressed)
	}

	if header.Unsynchronisated != false {
		t.Errorf("GetHeader failed: expected flag %v; got %v", false, header.Unsynchronisated)
	}

	if header.Size != 1138 {
		t.Errorf("GetHeader failed: expected size %v; got %v", 1138, header.Size)
	}
}
