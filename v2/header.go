package v2

import (
	"bytes"
	"fmt"
	"io"

	"github.com/Unbewohnte/id3ed/util"
)

// Main header`s flags
type HeaderFlags struct {
	Unsynchronised    bool
	Compressed        bool
	HasExtendedHeader bool
	Experimental      bool
	FooterPresent     bool
}

// ID3v2.x`s main header structure
type Header struct {
	// Identifier     string
	Flags          HeaderFlags
	Version        string
	Size           uint32
	ExtendedHeader ExtendedHeader
}

// extended header`s flags
type ExtendedHeaderFlags struct {
	UpdateTag       bool
	CRCpresent      bool
	HasRestrictions bool
	Restrictions    byte // a `lazy` approach :), just for now, maybe...
}

type ExtendedHeader struct {
	Size        uint32
	Flags       ExtendedHeaderFlags
	PaddingSize uint32
	CRCdata     []byte
}

// Reads and structuralises extended header. Must
// be called AFTER the main header has beeen read (does not seek).
// ALSO ISN`T TESTED !!!
func (h *Header) readExtendedHeader(r io.Reader) error {
	h.ExtendedHeader = ExtendedHeader{}
	if !h.Flags.HasExtendedHeader {
		return nil
	}

	var extended ExtendedHeader

	// extended size
	extendedSize, err := util.Read(r, 4)
	if err != nil {
		return fmt.Errorf("could not read from reader: %s", err)
	}

	switch h.Version {
	case V2_3:
		extended.Size = util.BytesToInt(extendedSize)
	case V2_4:
		extended.Size = util.BytesToIntSynchsafe(extendedSize)
	}

	// extended flags
	switch h.Version {
	case V2_3:
		extendedFlag, err := util.Read(r, 2) // reading flag byte and a null-byte after
		if err != nil {
			return fmt.Errorf("could not read from reader: %s", err)
		}
		if util.GetBit(extendedFlag[0], 1) {
			extended.Flags.CRCpresent = true
		} else {
			extended.Flags.CRCpresent = false
		}

	case V2_4:
		// skipping `Number of flag bytes` because it`s always `1`
		_, err := util.Read(r, 1)
		if err != nil {
			return fmt.Errorf("could not read from reader: %s", err)
		}
		flagByte, err := util.Read(r, 1)
		if err != nil {
			return fmt.Errorf("could not read from reader: %s", err)
		}

		if util.GetBit(flagByte[0], 2) {
			extended.Flags.UpdateTag = true
		} else {
			extended.Flags.UpdateTag = false
		}

		if util.GetBit(flagByte[0], 3) {
			extended.Flags.CRCpresent = true
		} else {
			extended.Flags.CRCpresent = false
		}

		if util.GetBit(flagByte[0], 4) {
			extended.Flags.HasRestrictions = true
		} else {
			extended.Flags.HasRestrictions = false
		}
	}

	// extracting data given by flags
	switch h.Version {
	case V2_3:
		if extended.Flags.CRCpresent {
			crcData, err := util.Read(r, 4)
			if err != nil {
				return fmt.Errorf("could not read from reader: %s", err)
			}
			extended.CRCdata = crcData
		}

	case V2_4:
		// `Each flag that is set in the extended header has data attached`

		if extended.Flags.UpdateTag {
			// skipping null-byte length of `UpdateTag`
			_, err := util.Read(r, 1)
			if err != nil {
				return fmt.Errorf("could not read from reader: %s", err)
			}
		}

		if extended.Flags.CRCpresent {
			crclen, err := util.Read(r, 1)
			if err != nil {
				return fmt.Errorf("could not read from reader: %s", err)
			}
			crcData, err := util.Read(r, uint64(crclen[0]))
			if err != nil {
				return fmt.Errorf("could not read from reader: %s", err)
			}
			extended.CRCdata = crcData
		}

		if extended.Flags.HasRestrictions {
			// skipping one-byte length of `Restrictions`, because it`s always `1`
			_, err := util.Read(r, 1)
			if err != nil {
				return fmt.Errorf("could not read from reader: %s", err)
			}

			restrictionsByte, err := util.Read(r, 1)
			if err != nil {
				return fmt.Errorf("could not read from reader: %s", err)
			}
			// a `lazy` approach :), just for now
			extended.Flags.Restrictions = restrictionsByte[0]
		}
	}

	// extracting other version-dependent header data

	// padding if V2_3
	if h.Version == V2_3 {
		paddingSizeBytes, err := util.Read(r, 4)
		if err != nil {
			return fmt.Errorf("could not read from reader: %s", err)
		}
		paddingSize := util.BytesToInt(paddingSizeBytes)
		extended.PaddingSize = paddingSize
	}

	// finally `attaching` parsed extended header to the main *Header
	h.ExtendedHeader = extended

	return nil
}

// Reads and structuralises ID3v2 header and if present - extended header.
// Returns a blank header struct if encountered an error
func readHeader(rs io.ReadSeeker) (Header, error) {
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

	// check if has ID3v2 identifier
	if !bytes.Equal([]byte(HEADERIDENTIFIER), identifier) {
		return Header{}, ErrDoesNotUseID3v2
	}

	// version
	majorVersion := int(hBytes[3])
	revisionNumber := int(hBytes[4])

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

	// v3.0 and v4.0 have different amount of flags
	switch header.Version {
	case V2_2:
		if util.GetBit(flags, 1) {
			header.Flags.Unsynchronised = true
		} else {
			header.Flags.Unsynchronised = false
		}
		if util.GetBit(flags, 2) {
			header.Flags.Compressed = true
		} else {
			header.Flags.Compressed = false
		}
	case V2_3:
		if util.GetBit(flags, 1) {
			header.Flags.Unsynchronised = true
		} else {
			header.Flags.Unsynchronised = false
		}
		if util.GetBit(flags, 2) {
			header.Flags.HasExtendedHeader = true
		} else {
			header.Flags.HasExtendedHeader = false
		}
		if util.GetBit(flags, 3) {
			header.Flags.Experimental = true
		} else {
			header.Flags.Experimental = false
		}
		// always false, because ID3v2.3.0 does not support footers
		header.Flags.FooterPresent = false

	case V2_4:
		if util.GetBit(flags, 1) {
			header.Flags.Unsynchronised = true
		} else {
			header.Flags.Unsynchronised = false
		}
		if util.GetBit(flags, 2) {
			header.Flags.HasExtendedHeader = true
		} else {
			header.Flags.HasExtendedHeader = false
		}
		if util.GetBit(flags, 3) {
			header.Flags.Experimental = true
		} else {
			header.Flags.Experimental = false
		}
		if util.GetBit(flags, 4) {
			header.Flags.FooterPresent = true
		} else {
			header.Flags.FooterPresent = false
		}
	}

	// size
	sizeBytes := hBytes[6:]

	size := util.BytesToIntSynchsafe(sizeBytes)

	header.Size = size

	if header.Flags.HasExtendedHeader {
		err = header.readExtendedHeader(rs)
		if err != nil {
			return header, err
		}
	}

	return header, nil
}

// Converts given HeaderFlags struct into ready-to-write byte
// containing flags
func headerFlagsToByte(hf HeaderFlags, version string) byte {
	var flagsByte byte = 0
	switch version {
	case V2_2:
		if hf.Unsynchronised {
			flagsByte = util.SetBit(flagsByte, 8)
		}
		if hf.Compressed {
			flagsByte = util.SetBit(flagsByte, 7)
		}

	case V2_3:
		if hf.Unsynchronised {
			flagsByte = util.SetBit(flagsByte, 8)
		}
		if hf.HasExtendedHeader {
			flagsByte = util.SetBit(flagsByte, 7)
		}
		if hf.Experimental {
			flagsByte = util.SetBit(flagsByte, 6)
		}

	case V2_4:
		if hf.Unsynchronised {
			flagsByte = util.SetBit(flagsByte, 8)
		}
		if hf.HasExtendedHeader {
			flagsByte = util.SetBit(flagsByte, 7)
		}
		if hf.Experimental {
			flagsByte = util.SetBit(flagsByte, 6)
		}
		if hf.FooterPresent {
			flagsByte = util.SetBit(flagsByte, 5)
		}

	}
	return flagsByte
}

// Converts given header into ready-to-write bytes
func (h *Header) toBytes() []byte {
	buff := new(bytes.Buffer)

	// id
	buff.Write([]byte(HEADERIDENTIFIER))

	// version
	version := []byte{0, 0}
	switch h.Version {
	case V2_2:
		version = []byte{2, 0}
	case V2_3:
		version = []byte{3, 0}
	case V2_4:
		version = []byte{4, 0}
	}
	buff.Write(version)

	// flags
	flagByte := headerFlagsToByte(h.Flags, h.Version)
	buff.WriteByte(flagByte)

	// size
	tagSize := util.IntToBytesSynchsafe(h.Size)
	buff.Write(tagSize)

	// extended header
	if !h.Flags.HasExtendedHeader {
		return buff.Bytes()
	}

	// double check for possible errors
	if h.Version == V2_2 {
		return buff.Bytes()
	}

	// size
	extSize := util.IntToBytes(h.ExtendedHeader.Size)
	buff.Write(extSize)

	// flags and other version specific fields
	switch h.Version {
	case V2_3:
		// flags
		flagBytes := []byte{0, 0}
		if h.ExtendedHeader.Flags.CRCpresent {
			flagBytes[0] = util.SetBit(flagBytes[0], 8)
		}
		buff.Write(flagBytes)

		// crc data
		if h.ExtendedHeader.Flags.CRCpresent {
			buff.Write(h.ExtendedHeader.CRCdata)
		}

		// padding size
		paddingSize := util.IntToBytes(h.ExtendedHeader.PaddingSize)
		buff.Write(paddingSize)

	case V2_4:
		numberOfFlagBytes := byte(1)
		buff.WriteByte(numberOfFlagBytes)

		extFlags := byte(0)
		if h.ExtendedHeader.Flags.UpdateTag {
			extFlags = util.SetBit(extFlags, 7)
		}
		if h.ExtendedHeader.Flags.CRCpresent {
			extFlags = util.SetBit(extFlags, 6)
		}
		if h.ExtendedHeader.Flags.HasRestrictions {
			extFlags = util.SetBit(extFlags, 5)
		}

		buff.WriteByte(extFlags)

		// writing data, provided by flags
		if h.ExtendedHeader.Flags.UpdateTag {
			// data len
			buff.WriteByte(0)
		}

		if h.ExtendedHeader.Flags.CRCpresent {
			// data len
			buff.WriteByte(5)
			// data
			buff.Write(h.ExtendedHeader.CRCdata)
		}

		if h.ExtendedHeader.Flags.HasRestrictions {
			// data len
			buff.WriteByte(1)
			// data
			buff.WriteByte(h.ExtendedHeader.Flags.Restrictions)
		}
	}

	return buff.Bytes()
}
