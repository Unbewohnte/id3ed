package v2

import (
	"fmt"
	"io"
	"strings"

	"github.com/Unbewohnte/id3ed/util"
)

type FrameFlags struct {
	TagAlterPreservation  bool
	FileAlterPreservation bool
	ReadOnly              bool
	Compressed            bool
	Encrypted             bool
	InGroup               bool
}

type Frame struct {
	ID        string
	Size      int64
	Flags     FrameFlags
	GroupByte byte
	Contents  []byte
}

var ErrGotPadding error = fmt.Errorf("got padding")

// Reads next ID3v2.3.0 or ID3v2.4.0 frame.
// Returns a blank Frame struct if encountered an error
func ReadFrame(rs io.Reader) (Frame, error) {
	var frame Frame

	// ID
	identifier, err := util.ReadToString(rs, 4)
	if err != nil {
		return Frame{}, err
	}
	if len(identifier) < 1 {
		// probably read all frames and got padding as identifier

		// I know that it`s a terrible desicion, but with my current
		// implementation it`s the only way I can see that will somewhat work

		return Frame{}, ErrGotPadding
	}
	frame.ID = identifier

	// Size
	framesizeBytes, err := util.Read(rs, 4)
	if err != nil {
		return Frame{}, err
	}

	framesize, err := util.BytesToIntIgnoreFirstBit(framesizeBytes)
	if err != nil {
		return Frame{}, err
	}

	frame.Size = framesize

	// Flags

	frameFlagsByte1, err := util.Read(rs, 1)
	if err != nil {
		return Frame{}, err
	}

	frameFlagsByte2, err := util.Read(rs, 1)
	if err != nil {
		return Frame{}, err
	}

	// I don`t have enough knowledge to handle this more elegantly
	// Any pointers ?

	flagsByte1Bits := fmt.Sprintf("%08b", frameFlagsByte1)
	flagsByte2Bits := fmt.Sprintf("%08b", frameFlagsByte2)
	var flags FrameFlags

	if flagsByte1Bits[0] == 1 {
		flags.TagAlterPreservation = true
	} else {
		flags.TagAlterPreservation = false
	}
	if flagsByte1Bits[1] == 1 {
		flags.FileAlterPreservation = true
	} else {
		flags.FileAlterPreservation = false
	}
	if flagsByte1Bits[2] == 1 {
		flags.ReadOnly = true
	} else {
		flags.ReadOnly = false
	}
	if flagsByte2Bits[0] == 1 {
		flags.Compressed = true
	} else {
		flags.Compressed = false
	}
	if flagsByte2Bits[1] == 1 {
		flags.Encrypted = true
	} else {
		flags.Encrypted = false
	}
	if flagsByte2Bits[2] == 1 {
		flags.InGroup = true
	} else {
		flags.InGroup = false
	}

	frame.Flags = flags

	if flags.InGroup {
		groupByte, err := util.Read(rs, 1)
		if err != nil {
			return Frame{}, err
		}
		frame.GroupByte = groupByte[0]
	}

	// Body
	frameContents, err := util.Read(rs, uint64(framesize))
	if err != nil {
		return Frame{}, err
	}

	frame.Contents = frameContents

	return frame, nil
}

// Reads all ID3v2 frames from rs.
// Returns a nil as []Frame if encountered an error
func GetFrames(rs io.ReadSeeker) ([]Frame, error) {
	// skip header
	_, err := rs.Seek(10, io.SeekStart)
	if err != nil {
		return nil, fmt.Errorf("could not skip header: %s", err)
	}

	var frames []Frame
	for {
		frame, err := ReadFrame(rs)
		if err == ErrGotPadding {
			return frames, nil
		}

		if err != nil {
			return nil, fmt.Errorf("could not read frame: %s", err)
		}

		frames = append(frames, frame)
	}
}

// Looks for a certain identificator in given frames and returns frame if found
func GetFrame(id string, frames []Frame) Frame {
	for _, frame := range frames {
		if strings.Contains(frame.ID, id) {
			return frame
		}
	}
	return Frame{}
}
