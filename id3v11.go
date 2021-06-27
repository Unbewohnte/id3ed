package id3ed

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
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

// Retrieves ID3v1.1 field values of provided io.ReadSeeker
func GetID3v11Tags(rs io.ReadSeeker) (*ID3v11Tags, error) {
	// set reader to the last 128 bytes
	_, err := rs.Seek(-int64(ID3v1SIZE), io.SeekEnd)
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

	track, err := binary.ReadVarint(bytes.NewBuffer(trackByte))
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

	return &ID3v11Tags{
		SongName: songname,
		Artist:   artist,
		Album:    album,
		Year:     year,
		Comment:  comment,
		Track:    int(track),
		Genre:    genre,
	}, nil
}

// Writes given ID3v1.1 tags to dst
func (tags *ID3v11Tags) Write(dst io.WriteSeeker) error {
	dst.Seek(0, io.SeekEnd)

	// TAG
	_, err := dst.Write([]byte("TAG"))
	if err != nil {
		return fmt.Errorf("could not write to dst: %s", err)
	}

	// Song name
	err = writeToExtent(dst, []byte(tags.SongName), 30)
	if err != nil {
		return fmt.Errorf("could not write to dst: %s", err)
	}

	// Artist
	err = writeToExtent(dst, []byte(tags.Artist), 30)
	if err != nil {
		return fmt.Errorf("could not write to dst: %s", err)
	}

	// Album
	err = writeToExtent(dst, []byte(tags.Album), 30)
	if err != nil {
		return fmt.Errorf("could not write to dst: %s", err)
	}

	// Year
	err = writeToExtent(dst, []byte(fmt.Sprint(tags.Year)), 4)
	if err != nil {
		return fmt.Errorf("could not write to dst: %s", err)
	}

	// Comment
	err = writeToExtent(dst, []byte(tags.Comment), 28)
	if err != nil {
		return fmt.Errorf("could not write to dst: %s", err)
	}

	_, err = dst.Write([]byte{0})
	if err != nil {
		return fmt.Errorf("could not write to dst: %s", err)
	}

	// Track
	trackBytes := make([]byte, 1)
	binary.PutVarint(trackBytes, int64(tags.Track))
	// binary.BigEndian.PutUint16(trackBytes, uint16(tags.Track))
	_, err = dst.Write(trackBytes)
	if err != nil {
		return fmt.Errorf("could not write to dst: %s", err)
	}

	//Genre
	genreCode := getKey(id3v1genres, tags.Genre)
	if genreCode == -1 {
		// if no genre found - encode genre code as 255
		genreCode = ID3v1INVALIDGENRE
	}
	genrebyte := make([]byte, 1)
	binary.PutVarint(genrebyte, int64(genreCode))

	err = writeToExtent(dst, genrebyte, 1)
	if err != nil {
		return fmt.Errorf("could not write to dst: %s", err)
	}

	return nil
}

// Checks for existing ID3v1.1 tag in file, if present - removes it and replaces with provided tags
func (tags *ID3v11Tags) WriteToFile(f *os.File) error {
	defer f.Close()

	// check for existing ID3v1.1 tag
	f.Seek(-int64(ID3v1SIZE), io.SeekEnd)

	tag, err := read(f, 3)
	if err != nil {
		return err
	}

	if !bytes.Equal(tag, []byte("TAG")) {
		// no existing tag, just write given tags
		err = tags.Write(f)
		if err != nil {
			return err
		}
		return nil
	}

	// does contain ID3v1.1 tag. Removing it
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
