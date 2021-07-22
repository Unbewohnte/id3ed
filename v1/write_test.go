package v1

import (
	"os"
	"path/filepath"
	"testing"
)

var TESTDATAPATH string = filepath.Join("..", "testData")

func TestWriteID3v1ToFile(t *testing.T) {
	f, err := os.OpenFile(filepath.Join(TESTDATAPATH, "testwritev1.mp3"), os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		t.Errorf("%s", err)
	}

	tag := TESTv1TAG

	err = tag.WriteToFile(f)
	if err != nil {
		t.Errorf("WriteID3v1ToFile failed: %s", err)
	}

}
