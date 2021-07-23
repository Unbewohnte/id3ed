package v2

import "testing"

func TestNewTAG(t *testing.T) {
	frame1, err := NewFrame("TTST", []byte("TEST text FRAME (ᗜˬᗜ)"), true)
	if err != nil {
		t.Errorf("%s", err)
	}
	frame2, err := NewFrame("BNRY", []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 255}, false)
	if err != nil {
		t.Errorf("%s", err)
	}

	newtag := NewTAG([]Frame{*frame1, *frame2})

	if newtag.Header.Version() != V2_4 {
		t.Errorf("NewTAG failed: expected version to be %s; got %s",
			V2_4, newtag.Header.Version())
	}

	var size uint32 = 0
	for _, frame := range newtag.Frames {
		size += uint32(len(frame.toBytes()))
	}

	if newtag.Header.Size() != size {
		t.Errorf("NewTAG failed: expected size to be %d; got %d",
			size, newtag.Header.Size())
	}

	// t.Errorf("%+v", newtag)

}
