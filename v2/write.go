package v2

// // Writes ID3v2Tag to ws
// func (tag *ID3v2Tag) write(ws io.WriteSeeker) error {
// 	_, err := ws.Seek(0, io.SeekStart)
// 	if err != nil {
// 		return fmt.Errorf("could not seek: %s", err)
// 	}

// 	// write header
// 	_, err = ws.Write(tag.Header.toBytes())
// 	if err != nil {
// 		return fmt.Errorf("could not write to writer: %s", err)
// 	}

// 	// write frames
// 	for _, frame := range tag.Frames {
// 		_, err = ws.Write(frame.toBytes())
// 		if err != nil {
// 			return fmt.Errorf("could not write to writer: %s", err)
// 		}
// 	}

// 	return nil
// }

// // Writes ID3v2Tag to file, removing already existing tag if found
// func (tag *ID3v2Tag) WriteToFile(f *os.File) error {
// 	defer f.Close()

// 	_, err := f.Seek(0, io.SeekStart)
// 	if err != nil {
// 		return fmt.Errorf("could not seek: %s", err)
// 	}

// 	// check for existing tag
// 	possibleHeaderID, err := util.ReadToString(f, 3)
// 	if err != nil {
// 		return err
// 	}

// 	if possibleHeaderID != HEADERIDENTIFIER {
// 		// No existing tag, just write what we have
// 		// and exit
// 		tag.write(f)

// 		return nil
// 	}
// 	// there is an existing tag, remove it
// 	// and write a new one

// 	// get size of the existing tag
// 	existingHeader, err := readHeader(f)
// 	if err != nil {
// 		return err
// 	}
// 	existingHSize := existingHeader.Size()

// 	// cannot truncate just the existing tag with f.Truncate(),
// 	// so we need to improvise and have a temporary copy of the mp3,
// 	// wipe the original file, write our tag and place the actual
// 	// music without the old tag from the temporary copy.

// 	// create a temporary file
// 	temporaryDir := os.TempDir()
// 	tmpF, err := os.CreateTemp(temporaryDir, fmt.Sprintf("%s_TEMP", filepath.Base(f.Name())))
// 	if err != nil {
// 		return err
// 	}

// 	defer tmpF.Close()
// 	// remove it afterwards
// 	defer os.Remove(filepath.Join(temporaryDir, tmpF.Name()))

// 	tmpFStats, err := tmpF.Stat()
// 	if err != nil {
// 		return err
// 	}

// 	// copy contents from the original mp3 to a temporary one
// 	_, err = io.Copy(tmpF, f)
// 	if err != nil {
// 		return err
// 	}

// 	// fully remove contents from the original file
// 	err = f.Truncate(0)
// 	if err != nil {
// 		return err
// 	}

// 	// write our tag
// 	tag.write(f)

// 	// read all contents of the temporary file, except the existing tag
// 	tmpF.Seek(int64(existingHSize), io.SeekStart)

// 	musicDataSize := uint64(tmpFStats.Size() - int64(existingHSize))

// 	musicData, err := util.Read(tmpF, musicDataSize)
// 	if err != nil {
// 		return err
// 	}

// 	// and write them into the original file, which
// 	// contains only the new tag
// 	_, err = f.Write(musicData)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
