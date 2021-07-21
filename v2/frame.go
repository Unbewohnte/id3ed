package v2

import (
	"bytes"
	"io"
	"unicode"

	"github.com/Unbewohnte/id3ed/util"
)

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

// Checks if given identifier is valid by specification
func isValidID(frameID string) bool {
	// check if id is in ASCII table
	if !util.InASCII(frameID) {
		return false
	}

	// check if id is in upper case
	for _, char := range frameID {
		if !unicode.IsUpper(char) {
			return false
		}
	}

	return true
}

func getV22FrameHeader(fHeaderbytes []byte) FrameHeader {
	var header FrameHeader

	header.ID = string(fHeaderbytes[0:3])

	framesizeBytes := util.BytesToIntSynchsafe(fHeaderbytes[3:6])
	header.Size = framesizeBytes

	return header
}

func getV23FrameHeader(fHeaderbytes []byte) FrameHeader {
	var header FrameHeader

	// ID
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

	return header
}

func getV24FrameHeader(fHeaderbytes []byte) FrameHeader {
	var header FrameHeader

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

	return header
}

// Structuralises frame header from given bytes. For versions see: constants.
func getFrameHeader(fHeaderbytes []byte, version string) FrameHeader {
	var header FrameHeader

	switch version {
	case V2_2:
		header = getV22FrameHeader(fHeaderbytes)

	case V2_3:
		header = getV23FrameHeader(fHeaderbytes)

	case V2_4:
		header = getV24FrameHeader(fHeaderbytes)
	}

	return header
}

// Reads a frame from r.
// Returns a blank Frame struct if encountered an error.
func readNextFrame(r io.Reader, version string) (Frame, error) {
	var frame Frame

	// Frame header
	var headerBytes []byte
	var err error
	switch version {
	case V2_2:
		headerBytes, err = util.Read(r, uint64(V2_2FrameHeaderSize))
		if err != nil {
			return Frame{}, err
		}
	default:
		headerBytes, err = util.Read(r, uint64(V2_3FrameHeaderSize))
		if err != nil {
			return Frame{}, err
		}
	}

	// check for padding and validate ID characters
	if bytes.Contains(headerBytes[0:3], []byte{0}) {
		return Frame{}, ErrGotPadding
	}

	if !isValidID(string(headerBytes[0:3])) {
		return Frame{}, ErrInvalidID
	}

	frameHeader := getFrameHeader(headerBytes, version)
	frame.Header = frameHeader

	// Contents
	contents, err := util.Read(r, uint64(frameHeader.Size))
	if err != nil {
		return Frame{}, err
	}
	frame.Contents = contents

	return frame, nil
}

// Returns decoded string from f.Contents.
// Note that it can and probably will return
// corrupted data if you use it on non-text frames such as APIC
// for such cases please deal with raw []byte
func (f *Frame) Text() string {
	return util.DecodeText(f.Contents)
}

func frameFlagsToBytes(ff FrameFlags, version string) []byte {
	var flagBytes = []byte{0, 0}

	switch version {
	case V2_2:
		return nil

	case V2_3:
		if ff.TagAlterPreservation {
			flagBytes[0] = util.SetBit(flagBytes[0], 8)
		}
		if ff.FileAlterPreservation {
			flagBytes[0] = util.SetBit(flagBytes[0], 7)
		}
		if ff.ReadOnly {
			flagBytes[0] = util.SetBit(flagBytes[0], 6)
		}

		if ff.Compressed {
			flagBytes[1] = util.SetBit(flagBytes[1], 8)
		}
		if ff.Encrypted {
			flagBytes[1] = util.SetBit(flagBytes[1], 7)
		}
		if ff.InGroup {
			flagBytes[1] = util.SetBit(flagBytes[1], 6)
		}
		return flagBytes

	case V2_4:
		if ff.TagAlterPreservation {
			flagBytes[0] = util.SetBit(flagBytes[0], 7)
		}
		if ff.FileAlterPreservation {
			flagBytes[0] = util.SetBit(flagBytes[0], 6)
		}
		if ff.ReadOnly {
			flagBytes[0] = util.SetBit(flagBytes[0], 5)
		}

		if ff.InGroup {
			flagBytes[1] = util.SetBit(flagBytes[1], 7)
		}
		if ff.Compressed {
			flagBytes[1] = util.SetBit(flagBytes[1], 4)
		}
		if ff.Encrypted {
			flagBytes[1] = util.SetBit(flagBytes[1], 3)
		}
		if ff.Unsyrchronised {
			flagBytes[1] = util.SetBit(flagBytes[1], 2)
		}
		if ff.HasDataLengthIndicator {
			flagBytes[1] = util.SetBit(flagBytes[1], 1)
		}
		return flagBytes

	default:
		return nil
	}
}

// Converts frame to ready-to-write bytes
func (f *Frame) toBytes(version string) []byte {
	buff := new(bytes.Buffer)

	// identifier
	buff.Write([]byte(f.Header.ID))

	// size
	buff.Write(util.IntToBytesSynchsafe(f.Header.Size))

	// flags
	flagBytes := frameFlagsToBytes(f.Header.Flags, version)
	if flagBytes != nil {
		buff.Write(flagBytes)
	}

	// contents
	buff.Write(f.Contents)

	return buff.Bytes()
}
