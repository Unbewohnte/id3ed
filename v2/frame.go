package v2

import (
	"io"

	"github.com/Unbewohnte/id3ed/util"
)

////////////////////////////////////////////////////////////////////////////
//(ᗜˬᗜ)~⭐//Under construction//Please don`t use it in this verison//(ᗜ‸ᗜ)///
////////////////////////////////////////////////////////////////////////////

type FrameHeader struct {
	Identifier string
	FrameSize  int64
}

type Frame struct {
	Header   *FrameHeader
	Contents string
}

// NOT TESTED !
func Readv2Frame(rs io.Reader) (*Frame, error) {
	var frameHeader *FrameHeader
	var frame Frame

	identifier, err := util.ReadToString(rs, 3)
	if err != nil {
		return nil, err
	}
	frameHeader.Identifier = identifier

	framesizeBytes, err := util.Read(rs, 3)
	if err != nil {
		return nil, err
	}

	framesize, err := util.BytesToInt(framesizeBytes)
	if err != nil {
		return nil, err
	}

	frameHeader.FrameSize = framesize

	frameContents, err := util.ReadToString(rs, int(framesize))
	if err != nil {
		return nil, err
	}

	frame.Header = frameHeader
	frame.Contents = frameContents

	return &frame, nil
}

// func ReadFrame(rs io.Reader, version string) error {
// 	return nil
// }
