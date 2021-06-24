package id3ed

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
)

type ID3v11Tags struct {
	SongName string
	Artist   string
	Album    string
	Year     int
	Comment  string
	Track    int
	Genre    string
}

// Retrieves ID3v1.1 field values of provided io.ReadSeeker (usually a file)
func GetID3v11Tags(rs io.ReadSeeker) (*ID3v11Tags, error) {
	// set reader to the last 128 bytes
	_, err := rs.Seek(-int64(ID3V1SIZE), io.SeekEnd)
	if err != nil {
		return nil, fmt.Errorf("could not seek: %s", err)
	}

	tag, err := read(rs, 3)
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(tag, []byte("TAG")) {
		// no TAG, given file does not use ID3v1
		return nil, fmt.Errorf("does not use ID3v1: expected %s; got %s", "TAG", tag)
	}

	songname, err := readToString(rs, 30)
	if err != nil {
		return nil, err
	}

	artist, err := readToString(rs, 30)
	if err != nil {
		return nil, err
	}

	album, err := readToString(rs, 30)
	if err != nil {
		return nil, err
	}

	yearStr, err := readToString(rs, 4)
	if err != nil {
		return nil, err
	}
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return nil, fmt.Errorf("could not convert yearbytes into int: %s", err)
	}

	comment, err := readToString(rs, 28)
	if err != nil {
		return nil, err
	}

	// skip 1 null byte
	_, err = read(rs, 1)
	if err != nil {
		return nil, err
	}

	trackByte, err := read(rs, 1)
	if err != nil {
		return nil, err
	}

	// track is one byte by specification
	track := int(trackByte[0])

	genreByte, err := read(rs, 1)
	if err != nil {
		return nil, err
	}
	// genre is one byte by specification
	genre, exists := id3v1genres[int(genreByte[0])]
	if !exists {
		genre = ""
	}

	return &ID3v11Tags{
		SongName: songname,
		Artist:   artist,
		Album:    album,
		Year:     year,
		Comment:  comment,
		Track:    track,
		Genre:    genre,
	}, nil
}

// Writes given ID3v1.1 tags to dst
func SetID3v11Tags(dst io.WriteSeeker, tags ID3v11Tags) error {
	dst.Seek(0, io.SeekEnd)

	return nil
}
