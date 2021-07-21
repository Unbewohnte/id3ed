package v2

// func TestWrite(t *testing.T) {
// 	f, err := os.Open(filepath.Join(TESTDATAPATH, "testreadv2.mp3"))
// 	if err != nil {
// 		t.Errorf("%s", err)
// 	}
// 	defer f.Close()

// 	testTag, err := ReadV2Tag(f)
// 	if err != nil {
// 		t.Errorf("%s", err)
// 	}

// 	ff, err := os.OpenFile(filepath.Join(TESTDATAPATH, "testwritev2.mp3"),
// 		os.O_CREATE|os.O_RDWR, os.ModePerm)
// 	if err != nil {
// 		t.Errorf("%s", err)
// 	}
// 	defer ff.Close()

// 	err = testTag.write(ff)
// 	if err != nil {
// 		t.Errorf("%s", err)
// 	}

// 	wroteTag, err := ReadV2Tag(ff)
// 	if err != nil {
// 		t.Errorf("%s", err)
// 	}

// 	// t.Errorf("ORIGINAL: %+v", testTag)
// 	// t.Errorf("WRITTEN: %+v", wroteTag)
// 	for _, origfr := range testTag.Frames {
// 		t.Errorf("ORIG Fr: %+v\n", origfr)
// 	}

// 	for _, wrtfr := range wroteTag.Frames {
// 		t.Errorf("WRITTEN Fr: %+v\n", wrtfr)
// 	}
// }
