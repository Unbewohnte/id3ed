package v2

import (
	"bytes"
	"fmt"
	"io"

	"github.com/Unbewohnte/id3ed/util"
)

type HeaderFlags struct {
	Unsynchronisated  bool
	HasExtendedHeader bool
	Experimental      bool
	FooterPresent     bool
}

// ID3v2.x header structure
type Header struct {
	Identifier string
	Flags      HeaderFlags
	Version    string
	Size       int64 // size of the whole tag - 10 header bytes
}

// Reads and structuralises ID3v2.3.0 or ID3v2.4.0 header.
// Returns a blank header struct if encountered an error
func GetHeader(rs io.ReadSeeker) (Header, error) {
	var header Header

	rs.Seek(0, io.SeekStart)

	identifier, err := util.Read(rs, 3)
	if err != nil {
		return Header{}, err
	}
	// check if ID3v2 is used
	if !bytes.Equal([]byte(HEADERIDENTIFIER), identifier) {
		return Header{}, fmt.Errorf("no ID3v2 identifier found")
	}
	header.Identifier = string(identifier)

	// version
	VersionBytes, err := util.Read(rs, 2)
	if err != nil {
		return Header{}, err
	}

	majorVersion, err := util.ByteToInt(VersionBytes[0])
	if err != nil {
		return Header{}, err
	}
	revisionNumber, err := util.ByteToInt(VersionBytes[1])
	if err != nil {
		return Header{}, err
	}

	var version string
	switch majorVersion {
	case 2:
		version = V2_2
	case 3:
		version = V2_3
	case 4:
		version = V2_4
	default:
		return Header{}, fmt.Errorf("ID3v2.%d.%d is not supported", majorVersion, revisionNumber)
	}

	header.Version = version

	// flags
	flags, err := util.Read(rs, 1)
	if err != nil {
		return Header{}, err
	}

	flagBits := fmt.Sprintf("%08b", flags) // 1 byte is 8 bits

	// v3.0 and v4.0 have different amount of flags
	switch version {
	case V2_3:
		if flagBits[0] == 1 {
			header.Flags.Unsynchronisated = true
		} else {
			header.Flags.Unsynchronisated = false
		}
		if flagBits[1] == 1 {
			header.Flags.HasExtendedHeader = true
		} else {
			header.Flags.HasExtendedHeader = false
		}
		if flagBits[2] == 1 {
			header.Flags.Experimental = true
		} else {
			header.Flags.Experimental = false
		}
		// always false, because ID3v2.3.0 does not support footers
		header.Flags.FooterPresent = false

	case V2_4:
		if flagBits[0] == 1 {
			header.Flags.Unsynchronisated = true
		} else {
			header.Flags.Unsynchronisated = false
		}
		if flagBits[1] == 1 {
			header.Flags.HasExtendedHeader = true
		} else {
			header.Flags.HasExtendedHeader = false
		}
		if flagBits[2] == 1 {
			header.Flags.Experimental = true
		} else {
			header.Flags.Experimental = false
		}
		if flagBits[3] == 1 {
			header.Flags.FooterPresent = true
		} else {
			header.Flags.FooterPresent = false
		}
	}

	// size
	sizeBytes, err := util.Read(rs, 4)
	if err != nil {
		return Header{}, err
	}

	size, err := util.BytesToIntIgnoreFirstBit(sizeBytes)
	if err != nil {
		return Header{}, err
	}

	header.Size = size

	return header, nil
}
