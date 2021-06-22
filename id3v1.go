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
	Genre    int
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
		return nil, fmt.Errorf("does not use ID3v1")
	}

	songname, err := read(rs, 30)
	if err != nil {
		return nil, err
	}

	artist, err := read(rs, 30)
	if err != nil {
		return nil, err
	}

	album, err := read(rs, 30)
	if err != nil {
		return nil, err
	}

	yearBytes, err := read(rs, 4)
	if err != nil {
		return nil, err
	}
	year, err := strconv.Atoi(string(yearBytes))
	if err != nil {
		return nil, fmt.Errorf("could not convert yearbytes into int: %s", err)
	}

	comment, err := read(rs, 30)
	if err != nil {
		return nil, err
	}

	genreByte, err := read(rs, 1)
	if err != nil {
		return nil, err
	}
	// genre is one byte by specification
	genre := int(genreByte[0])

	return &ID3v1Tags{
		SongName: string(songname),
		Artist:   string(artist),
		Album:    string(album),
		Year:     year,
		Comment:  string(comment),
		Genre:    genre,
	}, nil
}

func (t *ID3v1Tags) GetSongName() string {
	return t.SongName
}

func (t *ID3v1Tags) GetArtist() string {
	return t.Artist
}

func (t *ID3v1Tags) GetAlbum() string {
	return t.Album
}

func (t *ID3v1Tags) GetYear() int {
	return t.Year
}

func (t *ID3v1Tags) GetComment() string {
	return t.Comment
}

func (t *ID3v1Tags) GetGenre() int {
	return t.Genre
}
