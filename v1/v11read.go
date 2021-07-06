package v1

import (
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/Unbewohnte/id3ed/util"
)

// Retrieves ID3v1.1 field values of provided io.ReadSeeker
func Getv11Tag(rs io.ReadSeeker) (*ID3v11Tag, error) {
	// set reader to the last 128 bytes
	_, err := rs.Seek(-int64(ID3v1SIZE), io.SeekEnd)
	if err != nil {
		return nil, fmt.Errorf("could not seek: %s", err)
	}

	identifier, err := util.Read(rs, 3)
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(identifier, []byte(ID3v1IDENTIFIER)) {
		// no identifier, given file does not use ID3v1
		return nil, fmt.Errorf("does not use ID3v1: expected %s; got %s", ID3v1IDENTIFIER, identifier)
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

	comment, err := util.ReadToString(rs, 28)
	if err != nil {
		return nil, err
	}

	// skip 1 null byte
	_, err = util.Read(rs, 1)
	if err != nil {
		return nil, err
	}

	trackByte, err := util.Read(rs, 1)
	if err != nil {
		return nil, err
	}

	track, err := util.ByteToInt(trackByte[0])
	if err != nil {
		return nil, fmt.Errorf("cannot convert bytes to int: %s", err)
	}

	genreByte, err := util.Read(rs, 1)
	if err != nil {
		return nil, err
	}
	genreInt, err := util.ByteToInt(genreByte[0])
	if err != nil {
		return nil, fmt.Errorf("cannot convert bytes to int: %s", err)
	}
	genre, exists := id3v1genres[int(genreInt)]
	if !exists {
		genre = ""
	}

	return &ID3v11Tag{
		SongName: songname,
		Artist:   artist,
		Album:    album,
		Year:     year,
		Comment:  comment,
		Track:    int(track),
		Genre:    genre,
	}, nil
}
