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

// Reads and structuralises ID3v2 header from given bytes.
// Returns a blank header struct if encountered an error
func ReadHeader(rs io.ReadSeeker) (Header, error) {
	_, err := rs.Seek(0, io.SeekStart)
	if err != nil {
		return Header{}, fmt.Errorf("could not seek: %s", err)
	}

	hBytes, err := util.Read(rs, uint64(HEADERSIZE))
	if err != nil {
		return Header{}, fmt.Errorf("could not read from reader: %s", err)
	}

	var header Header

	identifier := hBytes[0:3]

	// check if has identifier ID3v2
	if !bytes.Equal([]byte(HEADERIDENTIFIER), identifier) {
		return Header{}, fmt.Errorf("no ID3v2 identifier found")
	}
	header.Identifier = string(identifier)

	// version
	majorVersion, err := util.ByteToInt(hBytes[3])
	if err != nil {
		return Header{}, err
	}
	revisionNumber, err := util.ByteToInt(hBytes[4])
	if err != nil {
		return Header{}, err
	}

	switch majorVersion {
	case 2:
		header.Version = V2_2
	case 3:
		header.Version = V2_3
	case 4:
		header.Version = V2_4
	default:
		return Header{}, fmt.Errorf("ID3v2.%d.%d is not supported or invalid", majorVersion, revisionNumber)
	}

	// flags
	flags := hBytes[5]

	flagBits := fmt.Sprintf("%08b", flags) // 1 byte is 8 bits

	// v3.0 and v4.0 have different amount of flags
	switch header.Version {
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
	sizeBytes := hBytes[6:]

	size, err := util.BytesToIntIgnoreFirstBit(sizeBytes)
	if err != nil {
		return Header{}, err
	}

	header.Size = size

	return header, nil
}
