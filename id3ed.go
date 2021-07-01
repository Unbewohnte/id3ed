package id3ed

import (
	"fmt"
	"io"
	"os"

	v1 "github.com/Unbewohnte/id3ed/v1"
)

type Tagger interface {
	WriteToFile(*os.File) error
	Version() int
}

// Reads certain ID3 tags from io.ReadSeeker according to given version.
// Wrapper function to v1|v2 package functions
func ReadTags(rs io.ReadSeeker, version int) (Tagger, error) {
	switch version {

	case 10:
		// ID3v1
		t, err := v1.Getv1Tags(rs)
		if err != nil {
			return nil, err
		}
		return t, nil

	case 11:
		// ID3v1.1
		t, err := v1.Getv11Tags(rs)
		if err != nil {
			return nil, err
		}
		return t, nil

	case 23:
		// ID3v2.3
		return nil, fmt.Errorf("v2.3 is not supported")

	case 24:
		// ID3v2.4
		return nil, fmt.Errorf("v2.4 is not supported")
	}

	return nil, fmt.Errorf("invalid version or not supported")
}

// DOESN`T work for some reason !? err is coming from f.Seek(), but it`s completely fine
// when run directly from v1.WriteToFile() function.
// Writes tags to file.
func WriteTags(dstf *os.File, t Tagger) error {
	err := t.WriteToFile(dstf)
	if err != nil {
		return err
	}

	return nil
}
