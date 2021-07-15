package v1

import (
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/Unbewohnte/id3ed/util"
)

var errDoesNotUseEnhancedID3v1 error = fmt.Errorf("does not use enhanced ID3v1 tag")

// Checks if enhanced tag is used
func usesEnhancedTag(rs io.ReadSeeker) bool {
	_, err := rs.Seek(-int64(ID3v1SIZE+ENHANCEDSIZE), io.SeekEnd)
	if err != nil {
		return false
	}
	identifier, err := util.Read(rs, 4)
	if err != nil {
		return false
	}
	if !bytes.Equal(identifier, []byte(ID3v1ENHANCEDIDENTIFIER)) {
		return false
	}

	return true
}

// Tries to read enhanced ID3V1 tag from rs
func readEnhancedTag(rs io.ReadSeeker) (EnhancedID3v1Tag, error) {

	if !usesEnhancedTag(rs) {
		return EnhancedID3v1Tag{}, errDoesNotUseEnhancedID3v1
	}

	var enhanced EnhancedID3v1Tag

	// set reader into the position
	_, err := rs.Seek(-int64(ID3v1SIZE+ENHANCEDSIZE), io.SeekEnd)
	if err != nil {
		return enhanced, fmt.Errorf("could not seek: %s", err)
	}

	// songname
	songName, err := util.ReadToString(rs, 60)
	if err != nil {
		return EnhancedID3v1Tag{}, err
	}
	enhanced.SongName = songName

	artist, err := util.ReadToString(rs, 60)
	if err != nil {
		return enhanced, err
	}
	enhanced.Artist = artist

	// album
	album, err := util.ReadToString(rs, 60)
	if err != nil {
		return enhanced, err
	}
	enhanced.Album = album

	// speed
	speedByte, err := util.Read(rs, 1)
	if err != nil {
		return enhanced, err
	}

	var speed string
	switch speedByte[0] {
	case 0:
		speed = "Unset"
	case 1:
		speed = "Slow"
	case 2:
		speed = "Medium"
	case 3:
		speed = "Fast"
	case 4:
		speed = "Hardcore"
	}
	enhanced.Speed = speed

	// genre
	genre, err := util.ReadToString(rs, 30)
	if err != nil {
		return enhanced, err
	}
	enhanced.Genre = genre

	// time
	startTime, err := util.ReadToString(rs, 6)
	if err != nil {
		return enhanced, err
	}
	enhanced.StartTime = startTime

	endtime, err := util.ReadToString(rs, 6)
	if err != nil {
		return enhanced, err
	}
	enhanced.EndTime = endtime

	return enhanced, nil
}

// Retrieves ID3v1 field values from rs.
func Readv1Tag(rs io.ReadSeeker) (*ID3v1Tag, error) {
	var tag ID3v1Tag

	// check if uses enhanced tag
	if usesEnhancedTag(rs) {
		enhanced, _ := readEnhancedTag(rs)
		tag.EnhancedTag = enhanced
	}

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
		return nil, ErrDoesNotUseID3v1
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
	tag.Comment = util.ToStringLossy(comment)
	tag.Track = 0

	var track int = 0
	// check if 29th byte is null byte (v1.0 or v1.1)
	if comment[28] == 0 {
		// it is v1.1, track number exists
		track = int(comment[29])

		tag.Track = uint8(track)

		comment = comment[0:28]
		tag.Comment = util.ToStringLossy(comment)
	}

	// Genre
	genreByte, err := util.Read(rs, 1)
	if err != nil {
		return nil, err
	}
	genreInt := int(genreByte[0])

	genre, exists := id3v1genres[int(genreInt)]
	if !exists {
		genre = ""
	}
	tag.Genre = genre

	if track == 0 {
		tag.version = V1_0
	} else {
		tag.version = V1_1
	}

	return &tag, nil
}
