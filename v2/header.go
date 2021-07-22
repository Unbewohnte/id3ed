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
	flags          HeaderFlags
	version        string
	size           uint32
	extendedHeader ExtendedHeader
}

// Extended header`s flags
type ExtendedHeaderFlags struct {
	UpdateTag       bool
	CRCpresent      bool
	HasRestrictions bool
	Restrictions    byte // a `lazy` approach :), just for now, maybe...
}

type ExtendedHeader struct {
	size        uint32
	flags       ExtendedHeaderFlags
	paddingSize uint32
	crcData     []byte
}

// using ONLY getters on header, because
// header MUST NOT be changed manually
// from the outside of the package

func (h *Header) Version() string {
	return h.version
}

func (h *Header) Flags() HeaderFlags {
	return h.flags
}

func (h *Header) Size() uint32 {
	return h.size
}

func (h *Header) ExtendedHeader() *ExtendedHeader {
	if h.flags.HasExtendedHeader {
		return &h.extendedHeader
	}
	return nil
}

// Reads and structuralises extended header. Must
// be called AFTER the main header has beeen read (does not seek).
// ALSO ISN`T TESTED !!!
func (h *Header) readExtendedHeader(r io.Reader) error {
	// h.ExtendedHeader = ExtendedHeader{}
	if !h.Flags().HasExtendedHeader {
		return nil
	}

	var extended ExtendedHeader

	// extended size
	extendedSize, err := util.Read(r, 4)
	if err != nil {
		return fmt.Errorf("could not read from reader: %s", err)
	}

	switch h.Version() {
	case V2_3:
		extended.size = util.BytesToInt(extendedSize)
	case V2_4:
		extended.size = util.BytesToIntSynchsafe(extendedSize)
	}

	// extended flags
	switch h.Version() {
	case V2_3:
		extendedFlag, err := util.Read(r, 2) // reading flag byte and a null-byte after
		if err != nil {
			return fmt.Errorf("could not read from reader: %s", err)
		}
		if util.GetBit(extendedFlag[0], 1) {
			extended.flags.CRCpresent = true
		} else {
			extended.flags.CRCpresent = false
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
			extended.flags.UpdateTag = true
		} else {
			extended.flags.UpdateTag = false
		}

		if util.GetBit(flagByte[0], 3) {
			extended.flags.CRCpresent = true
		} else {
			extended.flags.CRCpresent = false
		}

		if util.GetBit(flagByte[0], 4) {
			extended.flags.HasRestrictions = true
		} else {
			extended.flags.HasRestrictions = false
		}
	}

	// extracting data given by flags
	switch h.Version() {
	case V2_3:
		if extended.flags.CRCpresent {
			crcData, err := util.Read(r, 4)
			if err != nil {
				return fmt.Errorf("could not read from reader: %s", err)
			}
			extended.crcData = crcData
		}

	case V2_4:
		// `Each flag that is set in the extended header has data attached`

		if extended.flags.UpdateTag {
			// skipping null-byte length of `UpdateTag`
			_, err := util.Read(r, 1)
			if err != nil {
				return fmt.Errorf("could not read from reader: %s", err)
			}
		}

		if extended.flags.CRCpresent {
			crclen, err := util.Read(r, 1)
			if err != nil {
				return fmt.Errorf("could not read from reader: %s", err)
			}
			crcData, err := util.Read(r, uint64(crclen[0]))
			if err != nil {
				return fmt.Errorf("could not read from reader: %s", err)
			}
			extended.crcData = crcData
		}

		if extended.flags.HasRestrictions {
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
			extended.flags.Restrictions = restrictionsByte[0]
		}
	}

	// extracting other version-dependent header data

	// padding if V2_3
	if h.Version() == V2_3 {
		paddingSizeBytes, err := util.Read(r, 4)
		if err != nil {
			return fmt.Errorf("could not read from reader: %s", err)
		}
		paddingSize := util.BytesToInt(paddingSizeBytes)
		extended.paddingSize = paddingSize
	}

	// finally `attaching` parsed extended header to the main *Header
	h.extendedHeader = extended

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
		header.version = V2_2
	case 3:
		header.version = V2_3
	case 4:
		header.version = V2_4
	default:
		return Header{}, fmt.Errorf("ID3v2.%d.%d is not supported or invalid", majorVersion, revisionNumber)
	}

	// flags
	flags := hBytes[5]

	// v3.0 and v4.0 have different amount of flags
	switch header.Version() {
	case V2_2:
		if util.GetBit(flags, 1) {
			header.flags.Unsynchronised = true
		} else {
			header.flags.Unsynchronised = false
		}
		if util.GetBit(flags, 2) {
			header.flags.Compressed = true
		} else {
			header.flags.Compressed = false
		}
	case V2_3:
		if util.GetBit(flags, 1) {
			header.flags.Unsynchronised = true
		} else {
			header.flags.Unsynchronised = false
		}
		if util.GetBit(flags, 2) {
			header.flags.HasExtendedHeader = true
		} else {
			header.flags.HasExtendedHeader = false
		}
		if util.GetBit(flags, 3) {
			header.flags.Experimental = true
		} else {
			header.flags.Experimental = false
		}
		// always false, because ID3v2.3.0 does not support footers
		header.flags.FooterPresent = false

	case V2_4:
		if util.GetBit(flags, 1) {
			header.flags.Unsynchronised = true
		} else {
			header.flags.Unsynchronised = false
		}
		if util.GetBit(flags, 2) {
			header.flags.HasExtendedHeader = true
		} else {
			header.flags.HasExtendedHeader = false
		}
		if util.GetBit(flags, 3) {
			header.flags.Experimental = true
		} else {
			header.flags.Experimental = false
		}
		if util.GetBit(flags, 4) {
			header.flags.FooterPresent = true
		} else {
			header.flags.FooterPresent = false
		}
	}

	// size
	sizeBytes := hBytes[6:]

	size := util.BytesToIntSynchsafe(sizeBytes)

	header.size = size

	if header.flags.HasExtendedHeader {
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
	switch h.Version() {
	case V2_2:
		version = []byte{2, 0}
	case V2_3:
		version = []byte{3, 0}
	case V2_4:
		version = []byte{4, 0}
	}
	buff.Write(version)

	// flags
	flagByte := headerFlagsToByte(h.flags, h.version)
	buff.WriteByte(flagByte)

	// size
	tagSize := util.IntToBytesSynchsafe(h.size)
	buff.Write(tagSize)

	// extended header
	if !h.flags.HasExtendedHeader {
		return buff.Bytes()
	}

	// double check for possible errors
	if h.Version() == V2_2 {
		return buff.Bytes()
	}

	// size
	extSize := util.IntToBytes(h.extendedHeader.size)
	buff.Write(extSize)

	// flags and other version specific fields
	switch h.Version() {
	case V2_3:
		// flags
		flagBytes := []byte{0, 0}
		if h.extendedHeader.flags.CRCpresent {
			flagBytes[0] = util.SetBit(flagBytes[0], 8)
		}
		buff.Write(flagBytes)

		// crc data
		if h.extendedHeader.flags.CRCpresent {
			buff.Write(h.extendedHeader.crcData)
		}

		// padding size
		paddingSize := util.IntToBytes(h.extendedHeader.paddingSize)
		buff.Write(paddingSize)

	case V2_4:
		numberOfFlagBytes := byte(1)
		buff.WriteByte(numberOfFlagBytes)

		extFlags := byte(0)
		if h.extendedHeader.flags.UpdateTag {
			extFlags = util.SetBit(extFlags, 7)
		}
		if h.extendedHeader.flags.CRCpresent {
			extFlags = util.SetBit(extFlags, 6)
		}
		if h.extendedHeader.flags.HasRestrictions {
			extFlags = util.SetBit(extFlags, 5)
		}

		buff.WriteByte(extFlags)

		// writing data, provided by flags
		if h.extendedHeader.flags.UpdateTag {
			// data len
			buff.WriteByte(0)
		}

		if h.extendedHeader.flags.CRCpresent {
			// data len
			buff.WriteByte(5)
			// data
			buff.Write(h.extendedHeader.crcData)
		}

		if h.extendedHeader.flags.HasRestrictions {
			// data len
			buff.WriteByte(1)
			// data
			buff.WriteByte(h.extendedHeader.flags.Restrictions)
		}
	}

	return buff.Bytes()
}
