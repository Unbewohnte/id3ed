package id3ed

import (
	"os"
	"path/filepath"
	"testing"
)

var TESTDATAPATH string = filepath.Join("testData")

func TestReadTags(t *testing.T) {
	f, err := os.OpenFile(filepath.Join(TESTDATAPATH, "testreadv1.mp3"), os.O_RDONLY, os.ModePerm)
	if err != nil {
		t.Errorf("%s", err)
	}

	_, err = ReadTags(f, 11)
	if err != nil {
		t.Errorf("ReadTags failed: %s", err)
	}

	_, err = ReadTags(f, 10)
	if err != nil {
		t.Errorf("ReadTags failed: %s", err)
	}
}

// func TestWriteTags(t *testing.T) {
// 	f, err := os.OpenFile(filepath.Join(TESTDATAPATH, "testreadv1.mp3"), os.O_RDONLY, os.ModePerm)
// 	if err != nil {
// 		t.Errorf("%s", err)
// 	}

// 	tagger, err := ReadTags(f, 11)
// 	if err != nil {
// 		t.Errorf("%s", err)
// 	}

// 	f2, err := os.OpenFile(filepath.Join(TESTDATAPATH, "testWriteTags.mp3"), os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
// 	if err != nil {
// 		t.Errorf("%s", err)
// 	}

// 	err = WriteTags(f2, tagger)
// 	if err != nil {
// 		t.Errorf("WriteTags failed: %s", err)
// 	}
// }
