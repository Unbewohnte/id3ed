package id3ed

import (
	"bytes"
	"encoding/binary"
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
	genreInt, err := binary.ReadVarint(bytes.NewBuffer(genreByte))
	if err != nil {
		return nil, err
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

// Writes given ID3v1.0 tags to dst
func WriteID3v1Tags(dst io.WriteSeeker, tags ID3v1Tags) error {
	dst.Seek(0, io.SeekEnd)

	// TAG
	_, err := dst.Write([]byte("TAG"))
	if err != nil {
		return err
	}

	// Song name
	err = writeToExtend(dst, []byte(tags.SongName), 30)
	if err != nil {
		return err
	}

	// Artist
	err = writeToExtend(dst, []byte(tags.Artist), 30)
	if err != nil {
		return err
	}

	// Album
	err = writeToExtend(dst, []byte(tags.Album), 30)
	if err != nil {
		return err
	}

	// Year
	err = writeToExtend(dst, []byte(fmt.Sprint(tags.Year)), 4)
	if err != nil {
		return err
	}

	// Comment
	err = writeToExtend(dst, []byte(tags.Comment), 30)
	if err != nil {
		return err
	}

	// Genre
	genreCode := getKey(id3v1genres, tags.Genre)
	if genreCode == -1 {
		// if no genre found - encode genre code as 255
		genreCode = 255
	}
	genrebyte := make([]byte, 1)
	binary.PutVarint(genrebyte, int64(genreCode))

	_, err = dst.Write(genrebyte)
	if err != nil {
		return err
	}

	return nil
}
