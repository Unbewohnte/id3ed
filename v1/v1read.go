package v1

import (
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/Unbewohnte/id3ed/util"
)

// Retrieves ID3v1 field values of provided io.ReadSeeker (usually a file)
func Getv1Tags(rs io.ReadSeeker) (*ID3v1Tags, error) {
	// set reader to the last 128 bytes
	_, err := rs.Seek(-int64(ID3v1SIZE), io.SeekEnd)
	if err != nil {
		return nil, fmt.Errorf("could not seek: %s", err)
	}

	tag, err := util.Read(rs, 3)
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(tag, []byte(ID3v1IDENTIFIER)) {
		// no TAG, given file does not use ID3v1
		return nil, fmt.Errorf("does not use ID3v1: expected %s; got %s", ID3v1IDENTIFIER, tag)
	}

	songname, err := util.ReadToString(rs, 30)
	if err != nil {
		return nil, err
	}

	artist, err := util.ReadToString(rs, 30)
	if err != nil {
		return nil, err
	}

	album, err := util.ReadToString(rs, 30)
	if err != nil {
		return nil, err
	}

	yearStr, err := util.ReadToString(rs, 4)
	if err != nil {
		return nil, err
	}
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return nil, fmt.Errorf("could not convert yearbytes into int: %s", err)
	}

	comment, err := util.ReadToString(rs, 30)
	if err != nil {
		return nil, err
	}

	genreByte, err := util.Read(rs, 1)
	if err != nil {
		return nil, err
	}
	genreInt, err := util.BytesToInt(genreByte)
	if err != nil {
		return nil, fmt.Errorf("cannot convert bytes to int: %s", err)
	}
	genre, exists := id3v1genres[int(genreInt)]
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
