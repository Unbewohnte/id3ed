package v2

import (
	"os"
	"path/filepath"
	"testing"
)

var TESTDATAPATH string = filepath.Join("..", "testData")

func TestReadHeader(t *testing.T) {
	f, err := os.Open(filepath.Join(TESTDATAPATH, "testreadv2.mp3"))
	if err != nil {
		t.Errorf("%s", err)
	}

	header, err := readHeader(f)
	if err != nil {
		t.Errorf("GetHeader failed: %s", err)
	}

	if header.Identifier != "ID3" {
		t.Errorf("GetHeader failed: expected identifier %s; got %s", "ID3", header.Identifier)
	}

	if header.Flags.HasExtendedHeader != false {
		t.Errorf("GetHeader failed: expected flag %v; got %v", false, header.Flags.HasExtendedHeader)
	}

	if header.Flags.Unsynchronisated != false {
		t.Errorf("GetHeader failed: expected flag %v; got %v", false, header.Flags.Unsynchronisated)
	}

	if header.Size != 1138 {
		t.Errorf("GetHeader failed: expected size %v; got %v", 1138, header.Size)
	}
}
