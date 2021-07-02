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

// NOT TESTED ! Reads v2.3 | v2.4 frame
func ReadFrame(rs io.Reader) (*Frame, error) {
	var frameHeader FrameHeader
	var frame Frame

	identifier, err := util.ReadToString(rs, 4)
	if err != nil {
		return nil, err
	}
	frameHeader.ID = identifier

	framesizeBytes, err := util.Read(rs, 4)
	if err != nil {
		return nil, err
	}

	framesize, err := util.BytesToIntIgnoreFirstBit(framesizeBytes)
	if err != nil {
		return nil, err
	}

	frameHeader.FrameSize = framesize

	frameFlagsBytes, err := util.Read(rs, 2)
	if err != nil {
		return nil, err
	}

	// STILL NOT IMPLEMENTED FLAG HANDLING  !
	frameFlags, err := util.BytesToInt(frameFlagsBytes)
	if err != nil {
		return nil, err
	}
	frameHeader.Flags = int(frameFlags)

	frameContents, err := util.Read(rs, uint64(framesize))
	if err != nil {
		return nil, err
	}

	frame.Header = frameHeader
	frame.Contents = string(frameContents)

	return &frame, nil
}
