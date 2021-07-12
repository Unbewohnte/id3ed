package v2

import (
	"fmt"
	"io"
)

// Reads the whole ID3v2 tag from rs
func ReadV2Tag(rs io.ReadSeeker) (*ID3v2Tag, error) {
	header, err := readHeader(rs)
	if err == ErrDoesNotUseID3v2 {
		return nil, err
	} else if err != nil {
		return nil, fmt.Errorf("could not get header: %s", err)

	}

	// collect frames
	var read uint64 = 10 // because already read header
	var frames []Frame
	for {
		if read > uint64(header.Size) {
			break
		}

		frame, r, err := readNextFrame(rs, header)
		if err == ErrGotPadding || err == ErrBiggerThanSize || err == ErrInvalidID {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("could not read frame: %s", err)
		}

		read += r

		frames = append(frames, frame)
	}

	return &ID3v2Tag{
		Header: header,
		Frames: frames,
	}, nil
}
