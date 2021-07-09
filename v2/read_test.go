package v2

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Unbewohnte/id3ed/util"
)

func TestReadV2Tag(t *testing.T) {
	f, err := os.Open(filepath.Join(TESTDATAPATH, "testreadv2.mp3"))
	if err != nil {
		t.Errorf("%s", err)
	}

	tag, err := ReadV2Tag(f)
	if err != nil {
		t.Errorf("GetV2Tag failed: %s", err)
	}

	titleFrame := tag.Frames["TIT2"]

	if util.ToStringLossy(titleFrame.Contents) != "title" {
		t.Errorf("ReadV2Tag failed: expected contents to be %s; got %s",
			"title", util.ToStringLossy(titleFrame.Contents))
	}
}
