package v2

import (
	"fmt"
	"io"

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

// Reads ID3v2.3.0 or ID3v2.4.0 frame
func ReadFrame(rs io.Reader, version uint) (*Frame, error) {
	var frame Frame

	// ID
	identifier, err := util.ReadToString(rs, 4)
	if err != nil {
		return nil, err
	}
	frame.ID = identifier

	// Size
	framesizeBytes, err := util.Read(rs, 4)
	if err != nil {
		return nil, err
	}

	framesize, err := util.BytesToIntIgnoreFirstBit(framesizeBytes)
	if err != nil {
		return nil, err
	}

	frame.Size = framesize

	// Flags

	frameFlagsByte1, err := util.Read(rs, 1)
	if err != nil {
		return nil, err
	}

	frameFlagsByte2, err := util.Read(rs, 1)
	if err != nil {
		return nil, err
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
			return nil, err
		}
		frame.GroupByte = groupByte[0]
	}

	// Body
	frameContents, err := util.Read(rs, uint64(framesize))
	if err != nil {
		return nil, err
	}

	frame.Contents = frameContents

	return &frame, nil
}

func WriteFlag() {

}
