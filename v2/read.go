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

	var read uint64 = 0
	var frames []Frame
	for {
		if read == uint64(header.Size()) {
			break
		} else if read > uint64(header.Size()) {
			// read more than required, but did not
			// encouter padding, something is wrong here
			return nil, ErrReadMoreThanSize
		}

		frame, err := readNextFrame(rs, header.Version())
		switch err {
		case nil:
		case ErrGotPadding:
			// expected error, just return what we`ve collected
			return &ID3v2Tag{
				Header: header,
				Frames: frames,
			}, nil
		case ErrInvalidID:
			// expected error, just return what we`ve collected
			return &ID3v2Tag{
				Header: header,
				Frames: frames,
			}, nil

		default:
			return nil, err
		}

		frames = append(frames, frame)

		// counting how many bytes read
		if header.Version() == V2_2 {
			read += uint64(V2_2FrameHeaderSize) + uint64(frame.Header.Size())
		} else {
			read += uint64(V2_3FrameHeaderSize) + uint64(frame.Header.Size())
		}
	}

	return &ID3v2Tag{
		Header: header,
		Frames: frames,
	}, nil
}
