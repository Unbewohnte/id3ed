package v1

import (
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/Unbewohnte/id3ed/util"
)

// Retrieves ID3v1 field values of provided io.ReadSeeker (usually a file)
func Getv1Tag(rs io.ReadSeeker) (*ID3v1Tag, error) {
	var tag ID3v1Tag

	// set reader to the last 128 bytes
	_, err := rs.Seek(-int64(ID3v1SIZE), io.SeekEnd)
	if err != nil {
		return nil, fmt.Errorf("could not seek: %s", err)
	}

	// ID
	identifier, err := util.Read(rs, 3)
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(identifier, []byte(ID3v1IDENTIFIER)) {
		// no identifier, given file does not use ID3v1
		return nil, fmt.Errorf("does not use ID3v1: expected %s; got %s", ID3v1IDENTIFIER, identifier)
	}

	// Songname
	songname, err := util.ReadToString(rs, 30)
	if err != nil {
		return nil, err
	}
	tag.SongName = songname

	// Artist
	artist, err := util.ReadToString(rs, 30)
	if err != nil {
		return nil, err
	}
	tag.Artist = artist

	// Album name
	album, err := util.ReadToString(rs, 30)
	if err != nil {
		return nil, err
	}
	tag.Album = album

	// Year
	yearStr, err := util.ReadToString(rs, 4)
	if err != nil {
		return nil, err
	}
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return nil, fmt.Errorf("could not convert yearbytes into int: %s", err)
	}
	tag.Year = year

	// Comment and Track
	comment, err := util.Read(rs, 30)
	if err != nil {
		return nil, err
	}
	tag.Comment = util.ToString(comment)
	tag.Track = 0

	var track int = 0
	// check if 29th byte is null byte (v1.0 or v1.1)
	if comment[28] == 0 {
		// it is v1.1, track number exists
		track, err = util.ByteToInt(comment[29])
		if err != nil {
			return nil, fmt.Errorf("could not get int from byte: %s", err)
		}
		tag.Track = uint8(track)

		comment = comment[0:28]
		tag.Comment = util.ToString(comment)
	}

	// Genre
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
	tag.Genre = genre

	if track == 0 {
		tag.Version = V1_0
	} else {
		tag.Version = V1_1
	}

	return &tag, nil
}
