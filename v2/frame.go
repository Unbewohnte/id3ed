package v2

import (
	"fmt"
	"io"
	"strings"

	"github.com/Unbewohnte/id3ed/util"
)

var ErrGotPadding error = fmt.Errorf("got padding")
var ErrBiggerThanSize error = fmt.Errorf("frame size is bigger than size of the whole tag")
var ErrInvalidFHeaderSize error = fmt.Errorf("frame header must be 6 or 10 bytes long")
var ErrInvalidID error = fmt.Errorf("invalid identifier")

type FrameFlags struct {
	TagAlterPreservation   bool
	FileAlterPreservation  bool
	ReadOnly               bool
	Compressed             bool
	Encrypted              bool
	InGroup                bool
	Unsyrchronised         bool
	HasDataLengthIndicator bool
}

type FrameHeader struct {
	ID    string
	Size  uint32
	Flags FrameFlags
}

type Frame struct {
	Header   FrameHeader
	Contents []byte
}

// Checks if provided frame identifier is valid by
// its length and presence of invalid characters.
func isValidFrameID(frameID []byte) bool {
	if len(frameID) != 3 && len(frameID) != 4 {
		return false
	}
	str := strings.ToValidUTF8(string(frameID), "invalidChar")
	if len(str) != 3 && len(str) != 4 {
		return false
	}

	return true
}

func getV22FrameHeader(fHeaderbytes []byte) (FrameHeader, error) {
	var header FrameHeader

	if !isValidFrameID(fHeaderbytes[0:3]) {
		return FrameHeader{}, ErrInvalidID
	}
	header.ID = string(fHeaderbytes[0:3])

	framesizeBytes := util.BytesToIntSynchsafe(fHeaderbytes[3:6])
	header.Size = framesizeBytes

	return header, nil
}

func getV23FrameHeader(fHeaderbytes []byte) (FrameHeader, error) {
	var header FrameHeader

	// ID
	if !isValidFrameID(fHeaderbytes[0:4]) {
		return FrameHeader{}, ErrInvalidID
	}
	header.ID = string(fHeaderbytes[0:4])

	// Size
	framesizeBytes := fHeaderbytes[4:8]

	framesize := util.BytesToIntSynchsafe(framesizeBytes)

	header.Size = framesize

	// Flags
	frameFlags1 := fHeaderbytes[8]
	frameFlags2 := fHeaderbytes[9]

	var flags FrameFlags

	if util.GetBit(frameFlags1, 1) {
		flags.TagAlterPreservation = true
	} else {
		flags.TagAlterPreservation = false
	}
	if util.GetBit(frameFlags1, 2) {
		flags.FileAlterPreservation = true
	} else {
		flags.FileAlterPreservation = false
	}
	if util.GetBit(frameFlags1, 3) {
		flags.ReadOnly = true
	} else {
		flags.ReadOnly = false
	}
	if util.GetBit(frameFlags2, 1) {
		flags.Compressed = true
	} else {
		flags.Compressed = false
	}
	if util.GetBit(frameFlags2, 1) {
		flags.Encrypted = true
	} else {
		flags.Encrypted = false
	}
	if util.GetBit(frameFlags2, 1) {
		flags.InGroup = true
	} else {
		flags.InGroup = false
	}

	header.Flags = flags

	return header, nil
}

func getV24FrameHeader(fHeaderbytes []byte) (FrameHeader, error) {
	var header FrameHeader

	// ID
	if !isValidFrameID(fHeaderbytes[0:4]) {
		return FrameHeader{}, ErrInvalidID
	}
	header.ID = string(fHeaderbytes[0:4])

	// Size
	framesizeBytes := fHeaderbytes[4:8]

	framesize := util.BytesToIntSynchsafe(framesizeBytes)

	header.Size = framesize

	// Flags
	frameFlags1 := fHeaderbytes[8]
	frameFlags2 := fHeaderbytes[9]

	var flags FrameFlags

	if util.GetBit(frameFlags1, 2) {
		flags.TagAlterPreservation = true
	} else {
		flags.TagAlterPreservation = false
	}
	if util.GetBit(frameFlags1, 3) {
		flags.FileAlterPreservation = true
	} else {
		flags.FileAlterPreservation = false
	}
	if util.GetBit(frameFlags1, 4) {
		flags.ReadOnly = true
	} else {
		flags.ReadOnly = false
	}
	if util.GetBit(frameFlags2, 2) {
		flags.InGroup = true
	} else {
		flags.InGroup = false
	}
	if util.GetBit(frameFlags2, 5) {
		flags.Compressed = true
	} else {
		flags.Compressed = false
	}
	if util.GetBit(frameFlags2, 6) {
		flags.Encrypted = true
	} else {
		flags.Encrypted = false
	}
	if util.GetBit(frameFlags2, 7) {
		flags.Unsyrchronised = true
	} else {
		flags.Unsyrchronised = false
	}
	if util.GetBit(frameFlags2, 8) {
		flags.HasDataLengthIndicator = true
	} else {
		flags.HasDataLengthIndicator = false
	}

	header.Flags = flags

	return header, nil
}

// Structuralises frame header from given bytes. For versions see: constants.
func getFrameHeader(fHeaderbytes []byte, version string) (FrameHeader, error) {
	// validation check
	if int(len(fHeaderbytes)) != int(10) && int(len(fHeaderbytes)) != int(6) {
		return FrameHeader{}, ErrInvalidFHeaderSize
	}

	var header FrameHeader
	var err error

	switch version {
	case V2_2:
		header, err = getV22FrameHeader(fHeaderbytes)
		if err != nil {
			return FrameHeader{}, err
		}

	case V2_3:
		header, err = getV23FrameHeader(fHeaderbytes)
		if err != nil {
			return FrameHeader{}, err
		}

	case V2_4:
		header, err = getV24FrameHeader(fHeaderbytes)
		if err != nil {
			return FrameHeader{}, err
		}
	}

	return header, nil
}

// Reads ID3v2.3.0 or ID3v2.4.0 frame from given frame bytes.
// Returns a blank Frame struct if encountered an error, amount of
// bytes read from io.Reader.
func readNextFrame(r io.Reader, h Header) (Frame, uint64, error) {
	var frame Frame
	var read uint64 = 0

	// Frame header
	headerBytes, err := util.Read(r, uint64(HEADERSIZE))
	if err != nil {
		return Frame{}, 0, err
	}

	read += uint64(HEADERSIZE)

	frameHeader, err := getFrameHeader(headerBytes, h.Version)
	if err == ErrGotPadding || err == ErrInvalidID {
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

// Converts frame to ready-to-write bytes
// func (f *Frame) ToBytes() []byte {
// 	buff := new(bytes.Buffer)

// 	// identifier
// 	buff.Write([]byte(f.Header.ID))
// 	// size
// 	buff.Write(util.IntToBytesSynchsafe(f.Header.Size))
// 	// flags

// 	return buff.Bytes()
// }

// Returns bytes of the frame that can be
// written into a file.
// func (f *Frame) Bytes() ([]byte, error) {
// 	header := f.Header
// 	contents := f.Contents

// 	var headerbytes []byte

// 	identifierBytes := []byte(header.ID)
// 	// sizeBytes

// 	return nil, nil
// }
