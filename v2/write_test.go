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

// 	// WRITING
// 	err = testTag.WriteToFile(ff)
// 	if err != nil {
// 		t.Errorf("WriteToFile failed: %s", err)
// 	}
// }
