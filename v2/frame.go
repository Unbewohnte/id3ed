package v2

import (
	"fmt"
	"io"

	"github.com/Unbewohnte/id3ed/util"
)

var ErrGotPadding error = fmt.Errorf("got padding")
var ErrBiggerThanSize error = fmt.Errorf("frame size is bigger than size of the whole tag")
var ErrInvalidFHeaderSize error = fmt.Errorf("frame header must be 6 or 10 bytes long")

type FrameFlags struct {
	TagAlterPreservation  bool
	FileAlterPreservation bool
	ReadOnly              bool
	Compressed            bool
	Encrypted             bool
	InGroup               bool
}

type FrameHeader struct {
	ID    string
	Size  int64
	Flags FrameFlags
}

type Frame struct {
	Header   FrameHeader
	Contents []byte
}

// Structuralises frame header from given bytes. For versions see: constants.
func getFrameHeader(fHeaderbytes []byte, version string) (FrameHeader, error) {
	// validation check
	if int(len(fHeaderbytes)) != int(10) && int(len(fHeaderbytes)) != int(6) {
		return FrameHeader{}, ErrInvalidFHeaderSize
	}

	var header FrameHeader

	switch version {
	case V2_2:
		header.ID = string(fHeaderbytes[0:3])

		framesizeBytes, err := util.BytesToIntIgnoreFirstBit(fHeaderbytes[3:6])
		if err != nil {
			return FrameHeader{}, err
		}
		header.Size = framesizeBytes

	case V2_3:
		fallthrough

	case V2_4:
		fallthrough

	default:
		// ID
		header.ID = string(fHeaderbytes[0:4])

		// Size
		framesizeBytes := fHeaderbytes[4:8]

		framesize, err := util.BytesToIntIgnoreFirstBit(framesizeBytes)
		if err != nil {
			return FrameHeader{}, err
		}

		header.Size = framesize

		// Flags
		frameFlagsByte1 := fHeaderbytes[8]
		frameFlagsByte2 := fHeaderbytes[9]

		// I don`t have enough knowledge to handle this more elegantly

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

		header.Flags = flags
	}

	return header, nil
}

// Reads ID3v2.3.0 or ID3v2.4.0 frame from given frame bytes.
// Returns a blank Frame struct if encountered an error, amount of
// bytes read from io.Reader.
func ReadNextFrame(r io.Reader, h Header) (Frame, uint64, error) {
	var frame Frame
	var read uint64 = 0

	// Frame header
	headerBytes, err := util.Read(r, uint64(HEADERSIZE))
	if err != nil {
		return Frame{}, 0, err
	}

	read += uint64(HEADERSIZE)

	frameHeader, err := getFrameHeader(headerBytes, h.Version)
	if err == ErrGotPadding {
		return Frame{}, read, err
	} else if err != nil {
		return Frame{}, read, fmt.Errorf("could not get header of a frame: %s", err)
	}

	frame.Header = frameHeader

	// Contents
	contents, err := util.Read(r, uint64(frameHeader.Size))
	if err != nil {
		return Frame{}, read, err
	}

	frame.Contents = contents

	read += uint64(frameHeader.Size)

	return frame, read, err
}

// Returns decoded string from f.Contents.
// Note that it can and probably will return
// corrupted data if you use it on non-text frames such as APIC
// for such cases please deal with raw []byte
func (f *Frame) Text() string {
	return util.DecodeText(f.Contents)
}
