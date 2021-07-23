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

	firstFrame, err := readNextFrame(f, header.Version())
	if err != nil {
		t.Errorf("ReadFrame failed: %s", err)
	}

	if firstFrame.Header.ID() != "TRCK" {
		t.Errorf("GetFrame failed: expected ID %s; got %s",
			"TRCK", firstFrame.Header.ID())
	}

	if firstFrame.Header.Flags().Encrypted != false {
		t.Errorf("ReadFrame failed: expected compressed flag to be %v; got %v",
			false, firstFrame.Header.Flags().Encrypted)
	}

	secondFrame, err := readNextFrame(f, header.Version())
	if err != nil {
		t.Errorf("ReadFrame failed: %s", err)
	}

	if secondFrame.Header.ID() != "TDRC" {
		t.Errorf("ReadFrame failed: expected ID %s; got %s",
			"TDRC", secondFrame.Header.ID())
	}

	if util.ToStringLossy(secondFrame.Contents) != "2006" {
		t.Errorf("ReadFrame failed: expected contents to be %s; got %s",
			"2006", secondFrame.Contents)
	}
}

func TestFrameFlagsToBytes(t *testing.T) {
	testFlags := FrameFlags{
		TagAlterPreservation: true,
		ReadOnly:             true,
	}

	versions := []string{V2_2, V2_3, V2_4}

	for _, version := range versions {
		flagBytes := frameFlagsToBytes(testFlags, version)
		if version == V2_2 && flagBytes != nil {
			t.Errorf("frameFlagsToBytes failed: V2_2, expected flagbytes to be nil; got %v", flagBytes)
		}
		if version != V2_2 && len(flagBytes) != 2 {
			t.Errorf("frameFlagsToBytes failed: expected flagbytes to be len of 2; got %v", len(flagBytes))
		}
	}
}

func TestFrameToBytes(t *testing.T) {
	testframe := Frame{
		Header: FrameHeader{
			id:    "TEST",
			flags: FrameFlags{}, // all false
			size:  4,
		},
		Contents: []byte{util.EncodingUTF8, 60, 60, 60}, // 60 == <
	}

	frameBytes := testframe.toBytes()

	// t.Errorf("%+v", frameBytes)
	// 84 69 83 84 0 0 0 4 0 0 3 60 60 60

	// 84 69 83 84 - id (4)
	// 0 0 - flags (2)
	// 0 4 0 0 - size (4)
	// header - 4 + 4 + 2 = 10 bytes (success)
	// 3 60 60 60 - contents

	if len(frameBytes)-int(testframe.Header.Size()) != HEADERSIZE {
		t.Errorf("FrameToBytes failed: expected header size to be %d; got %d",
			HEADERSIZE, len(frameBytes)-int(testframe.Header.Size()))
	}

	if util.DecodeText(frameBytes[10:]) != "<<<" {
		t.Errorf("FrameToBytes failed: expected contents to be %v; got %v",
			testframe.Contents, frameBytes[10:])
	}
}

func TestNewFrame(t *testing.T) {
	gotFrame, err := NewFrame("TMRK", []byte("Very cool"), true)
	if err != nil {
		t.Errorf("CreateFrame failed: %s", err)
	}

	// check for encoding byte
	if gotFrame.Contents[0] != util.EncodingUTF8 {
		t.Errorf("CreateFrame failed: contents are expected to have an encoding byte %v; got %v",
			util.EncodingUTF8, gotFrame.Contents[0])
	}

	if gotFrame.Header.Size() != uint32(len(gotFrame.Contents)) {
		t.Errorf("CreateFrame failed: expected size to be %d; got %d",
			len(gotFrame.Contents), gotFrame.Header.Size())
	}

}
