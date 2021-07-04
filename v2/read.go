package v2

// Reads ID3v2 frames from rs. NOT TESTED !!!!
// func GetFrames(rs io.ReadSeeker) ([]*Frame, error) {
// 	header, err := GetHeader(rs)
// 	if err != nil {
// 		return nil, fmt.Errorf("could not get header: %s", err)
// 	}
// 	tagsize := header.Size

// 	var frames []*Frame
// 	var read uint64 = 0
// 	for {
// 		if read == uint64(tagsize) {
// 			break
// 		}

// 		frame, err := ReadFrame(rs)
// 		if err != nil {
// 			return frames, fmt.Errorf("could not read frame: %s", err)
// 		}
// 		frames = append(frames, frame)

// 		// counting how many bytes has been read
// 		read += 10 // frame header
// 		if frame.Flags.InGroup {
// 			// header has 1 additional byte
// 			read += 1
// 		}
// 		read += uint64(frame.Size) // and the contents itself
// 	}

// 	return frames, nil
// }
