package id3ed

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
)

// https://id3.org/ID3v1 - documentation

const ID3V1SIZE int = 128

type ID3v1Tags struct {
	SongName string
	Artist   string
	Album    string
	Year     int
	Comment  string
	Genre    string
}

// Retrieves ID3v1 field values of provided io.ReadSeeker (usually a file)
func GetID3v1Tags(rs io.ReadSeeker) (*ID3v1Tags, error) {
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

	comment, err := readToString(rs, 30)
	if err != nil {
		return nil, err
	}

	genreByte, err := read(rs, 1)
	if err != nil {
		return nil, err
	}
	// genre is one byte by specification
	genre, exists := id3v1genres[int(genreByte[0])]
	if !exists {
		genre = ""
	}

	return &ID3v1Tags{
		SongName: songname,
		Artist:   artist,
		Album:    album,
		Year:     year,
		Comment:  comment,
		Genre:    genre,
	}, nil
}

// Writes given ID3v1.0 tags to dst
func SetID3v1Tags(dst io.WriteSeeker, tags ID3v11Tags) error {
	dst.Seek(0, io.SeekEnd)

	return nil
}
