package v2

import (
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/Unbewohnte/id3ed/util"
)

// ID3v2.x header structure
type Header struct {
	Identifier       string
	Version          string
	Unsynchronisated bool
	Compressed       bool
	Size             int64 // size of the whole tag - 10 header bytes
}

// Reads and structuralises ID3v2 header
func GetHeader(rs io.ReadSeeker) (*Header, error) {
	var header Header

	rs.Seek(0, io.SeekStart)

	identifier, err := util.Read(rs, 3)
	if err != nil {
		return nil, err
	}
	// check if ID3v2 is used
	if !bytes.Equal([]byte(HEADERIDENTIFIER), identifier) {
		return nil, fmt.Errorf("no ID3v2 identifier found")
	}
	header.Identifier = string(identifier)

	// version
	majorVersionByte, err := util.Read(rs, 1)
	if err != nil {
		return nil, err
	}
	revisionNumberByte, err := util.Read(rs, 1)
	if err != nil {
		return nil, err
	}

	majorVersion, err := util.BytesToInt(majorVersionByte)
	if err != nil {
		return nil, err
	}
	revisionNumber, err := util.BytesToInt(revisionNumberByte)
	if err != nil {
		return nil, err
	}
	header.Version = fmt.Sprintf("%d%d", -majorVersion, revisionNumber)

	// flags
	flags, err := util.Read(rs, 1)
	if err != nil {
		return nil, err
	}
	bits := fmt.Sprintf("%08b", flags) // 1 byte is 8 bits
	if bits[0] == 1 {
		header.Unsynchronisated = true
	} else {
		header.Unsynchronisated = false
	}
	if bits[1] == 1 {
		header.Compressed = true
	} else {
		header.Compressed = false
	}

	// size
	sizeBytes, err := util.Read(rs, 4)
	if err != nil {
		return nil, err
	}

	// represent each byte in size as binary and get rid from the first useless bit,
	// then concatenate filtered parts
	var filteredSizeStr string
	for _, b := range sizeBytes {
		// the first bit is always 0, so filter it out
		filteredPart := fmt.Sprintf("%08b", b)[1:] // byte is 8 bits
		filteredSizeStr += filteredPart
	}

	// converting filtered binary size into usable int64
	size, err := strconv.ParseInt(filteredSizeStr, 2, 64)
	if err != nil {
		return nil, err
	}

	header.Size = int64(size)

	return &header, nil
}
