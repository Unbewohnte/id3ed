package id3ed

//////////////////////////////////////
//(ᗜˬᗜ)~⭐//Under construction//(ᗜ‸ᗜ)//
//////////////////////////////////////

import (
	"bytes"
	"fmt"
	"io"
)

type Header struct {
	Identifier string
	Version    int
	Flags      int
	Size       int64
}

func GetHeader(rs io.ReadSeeker) (*Header, error) {
	var header Header

	rs.Seek(0, io.SeekStart)

	identifier, err := read(rs, 3)
	if err != nil {
		return nil, err
	}
	// check if ID3v2 is used
	if !bytes.Equal([]byte(ID3v2IDENTIFIER), identifier) {
		return nil, fmt.Errorf("does not use ID3v2")
	}
	////

	header.Identifier = string(identifier)

	return &header, nil
}
