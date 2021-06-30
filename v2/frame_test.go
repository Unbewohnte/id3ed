package v2

// func TestReadFrame(t *testing.T) {
// 	f, err := os.Open(filepath.Join(TESTDATAPATH, "testreadv2.mp3"))
// 	if err != nil {
// 		t.Errorf("%s", err)
// 	}

// 	// read right after header`s bytes
// 	f.Seek(int64(HEADERSIZE), io.SeekStart)

// 	_, err = Readv2Frame(f)
// 	if err != nil {
// 		t.Errorf("ReadFrame failed: %s", err)
// 	}
// }
