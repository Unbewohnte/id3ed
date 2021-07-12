package id3ed

import (
	"path/filepath"
	"testing"

	v1 "github.com/Unbewohnte/id3ed/v1"
)

var TESTDATAPATH string = "testData"

func TestOpen(t *testing.T) {
	file, err := Open(filepath.Join(TESTDATAPATH, "testreadv2.mp3"))
	if err != nil {
		t.Errorf("Open failed: %s", err)
	}

	if file.ContainsID3v1 {
		t.Error("Open failed: expected testing file to not contain ID3v1")
	}

	if !file.ContainsID3v2 {
		t.Error("Open failed: expected testing file to contain ID3v2")
	}
}

func TestWriteID3v1(t *testing.T) {
	file, err := Open(filepath.Join(TESTDATAPATH, "testwritev1.mp3"))
	if err != nil {
		t.Errorf("Open failed: %s", err)
	}
	v1tag := &v1.ID3v1Tag{
		SongName: "testsong",
		Artist:   "testartist",
		Album:    "testalbum",
		Year:     727,
		Comment:  "testcomment",
		Genre:    "Blues",
	}

	err = file.WriteID3v1(v1tag)
	if err != nil {
		t.Errorf("WriteID3v1 failed: %s", err)
	}
}
