package v2

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Unbewohnte/id3ed/util"
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

	if header.Flags.HasExtendedHeader != false {
		t.Errorf("GetHeader failed: expected flag %v; got %v", false, header.Flags.HasExtendedHeader)
	}

	if header.Flags.Unsynchronised != false {
		t.Errorf("GetHeader failed: expected flag %v; got %v", false, header.Flags.Unsynchronised)
	}

	if header.Size != 1138 {
		t.Errorf("GetHeader failed: expected size %v; got %v", 1138, header.Size)
	}
}

func TestHeaderFlagsToByte(t *testing.T) {
	hf := HeaderFlags{
		Experimental:  true,
		FooterPresent: true,
	}
	var correctFlagsByte byte = 0
	correctFlagsByte = util.SetBit(correctFlagsByte, 5)
	correctFlagsByte = util.SetBit(correctFlagsByte, 6)

	gotByte := headerFlagsToByte(hf, V2_4)

	if gotByte != correctFlagsByte {
		t.Errorf("headerFlagsToByte failed: expected to get %d; got %d", correctFlagsByte, gotByte)
	}
}

func TestHeaderToBytes(t *testing.T) {
	testHeader := Header{
		Version:        V2_4,
		Flags:          HeaderFlags{}, // all false
		Size:           12345,
		ExtendedHeader: ExtendedHeader{},
	}

	hBytes := testHeader.toBytes()

	// t.Errorf("%v", hBytes)
	// 73 68 51 4 0 0 0 0 96 57
	// 73 68 51 - identifier
	// 4 0 - version
	// 0 - flags
	// 0 0 96 57 - size

	if string(hBytes[0:3]) != HEADERIDENTIFIER {
		t.Errorf("expected to get %s, got %s", HEADERIDENTIFIER, string(hBytes[0:3]))
	}

	if util.BytesToIntSynchsafe(hBytes[6:10]) != testHeader.Size {
		t.Errorf("toBytes failed: expected size to be %d; got %d",
			testHeader.Size, util.BytesToIntSynchsafe(hBytes[7:10]))
	}
}
