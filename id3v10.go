package id3ed

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strconv"
)

// https://id3.org/ID3v1 - documentation

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
	_, err := rs.Seek(-int64(ID3v1SIZE), io.SeekEnd)
	if err != nil {
		return nil, fmt.Errorf("could not seek: %s", err)
	}

	tag, err := read(rs, 3)
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(tag, []byte(ID3v1IDENTIFIER)) {
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

// Writes given ID3v1.0 tags to given io.WriteSeeker.
func (tags *ID3v1Tags) Write(dst io.WriteSeeker) error {
	dst.Seek(0, io.SeekEnd)

	// TAG
	_, err := dst.Write([]byte(ID3v1IDENTIFIER))
	if err != nil {
		return err
	}

	// Song name
	err = writeToExtent(dst, []byte(tags.SongName), 30)
	if err != nil {
		return err
	}

	// Artist
	err = writeToExtent(dst, []byte(tags.Artist), 30)
	if err != nil {
		return err
	}

	// Album
	err = writeToExtent(dst, []byte(tags.Album), 30)
	if err != nil {
		return err
	}

	// Year
	err = writeToExtent(dst, []byte(fmt.Sprint(tags.Year)), 4)
	if err != nil {
		return err
	}

	// Comment
	err = writeToExtent(dst, []byte(tags.Comment), 30)
	if err != nil {
		return err
	}

	// Genre
	genreCode := getKey(id3v1genres, tags.Genre)
	if genreCode == -1 {
		// if no genre found - encode genre code as 255
		genreCode = ID3v1INVALIDGENRE
	}
	genrebyte := make([]byte, 1)
	binary.PutVarint(genrebyte, int64(genreCode))

	_, err = dst.Write(genrebyte)
	if err != nil {
		return err
	}

	return nil
}

// Checks for existing ID3v1 tag in file, if present - removes it and replaces with provided tags
func (tags *ID3v1Tags) WriteToFile(f *os.File) error {
	defer f.Close()

	// check for existing ID3v1 tag
	f.Seek(-int64(ID3v1SIZE), io.SeekEnd)

	tag, err := read(f, 3)
	if err != nil {
		return err
	}

	if !bytes.Equal(tag, []byte(ID3v1IDENTIFIER)) {
		// no existing tag, just write given tags
		err = tags.Write(f)
		if err != nil {
			return err
		}
		return nil
	}

	// does contain ID3v1 tag. Removing it
	fStats, err := f.Stat()
	if err != nil {
		return err
	}

	err = f.Truncate(fStats.Size() - int64(ID3v1SIZE))
	if err != nil {
		return nil
	}

	// writing new tags
	err = tags.Write(f)
	if err != nil {
		return err
	}

	return nil

}
