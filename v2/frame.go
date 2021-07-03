package v2

import (
	"io"

	"github.com/Unbewohnte/id3ed/util"
)

////////////////////////////////////////////////////////////////////////////
//(ᗜˬᗜ)~⭐//Under construction//Please don`t use it in this verison//(ᗜ‸ᗜ)///
////////////////////////////////////////////////////////////////////////////

type FrameHeader struct {
	ID        string
	FrameSize int64
	Flags     int
}

type Frame struct {
	Header   FrameHeader
	Contents string
}

// Reads ID3v2.3.0 or ID3v2.4.0 frame
func ReadFrame(rs io.Reader, version uint) (*Frame, error) {
	var frameHeader FrameHeader
	var frame Frame

	// ID
	identifier, err := util.ReadToString(rs, 4)
	if err != nil {
		return nil, err
	}
	frameHeader.ID = identifier

	// Size
	framesizeBytes, err := util.Read(rs, 4)
	if err != nil {
		return nil, err
	}

	framesize, err := util.BytesToIntIgnoreFirstBit(framesizeBytes)
	if err != nil {
		return nil, err
	}

	frameHeader.FrameSize = framesize

	// Flags
	frameFlagsBytes, err := util.Read(rs, 2)
	if err != nil {
		return nil, err
	}

	frameFlags, err := util.BytesToInt(frameFlagsBytes)
	if err != nil {
		return nil, err
	}

	frameHeader.Flags = int(frameFlags)

	// Body
	frameContents, err := util.Read(rs, uint64(framesize))
	if err != nil {
		return nil, err
	}

	frame.Header = frameHeader
	frame.Contents = string(frameContents)

	return &frame, nil
}
