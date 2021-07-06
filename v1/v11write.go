package v1

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/Unbewohnte/id3ed/util"
)

// Writes given ID3v1.1 tag to dst
// NOTE: will not remove already existing ID3v1.1 tag if it`s present,
// use ⁕WriteToFile⁕ method if you`re working with REAL mp3 files !!!
func (tag *ID3v11Tag) Write(dst io.WriteSeeker) error {
	dst.Seek(0, io.SeekEnd)

	// TAG
	_, err := dst.Write([]byte(ID3v1IDENTIFIER))
	if err != nil {
		return fmt.Errorf("could not write to dst: %s", err)
	}

	// Song name
	err = util.WriteToExtent(dst, []byte(tag.SongName), 30)
	if err != nil {
		return fmt.Errorf("could not write to dst: %s", err)
	}

	// Artist
	err = util.WriteToExtent(dst, []byte(tag.Artist), 30)
	if err != nil {
		return fmt.Errorf("could not write to dst: %s", err)
	}

	// Album
	err = util.WriteToExtent(dst, []byte(tag.Album), 30)
	if err != nil {
		return fmt.Errorf("could not write to dst: %s", err)
	}

	// Year
	err = util.WriteToExtent(dst, []byte(fmt.Sprint(tag.Year)), 4)
	if err != nil {
		return fmt.Errorf("could not write to dst: %s", err)
	}

	// Comment
	err = util.WriteToExtent(dst, []byte(tag.Comment), 28)
	if err != nil {
		return fmt.Errorf("could not write to dst: %s", err)
	}

	_, err = dst.Write([]byte{0})
	if err != nil {
		return fmt.Errorf("could not write to dst: %s", err)
	}

	// Track
	trackBytes := make([]byte, 1)
	binary.PutVarint(trackBytes, int64(tag.Track))
	_, err = dst.Write(trackBytes)
	if err != nil {
		return fmt.Errorf("could not write to dst: %s", err)
	}

	//Genre
	genreCode := util.GetKey(id3v1genres, tag.Genre)
	if genreCode == -1 {
		// if no genre found - encode genre code as 255
		genreCode = ID3v1INVALIDGENRE
	}
	genrebyte := make([]byte, 1)
	binary.PutVarint(genrebyte, int64(genreCode))

	err = util.WriteToExtent(dst, genrebyte, 1)
	if err != nil {
		return fmt.Errorf("could not write to dst: %s", err)
	}

	return nil
}

// Checks for existing ID3v1.1 tag in file, if present - removes it and replaces with provided tag
func (tag *ID3v11Tag) WriteToFile(f *os.File) error {
	defer f.Close()

	// check for existing ID3v1.1 tag
	_, err := f.Seek(-int64(ID3v1SIZE), io.SeekEnd)
	if err != nil {
		return fmt.Errorf("could not seek: %s", err)
	}

	identifier, err := util.Read(f, 3)
	if err != nil {
		// return err
		return err

	}

	if !bytes.Equal(identifier, []byte(ID3v1IDENTIFIER)) {
		// no existing identifier, just write given tag
		err = tag.Write(f)
		if err != nil {
			return fmt.Errorf("could not write to writer: %s", err)
		}
		return nil
	}

	// does contain ID3v1.1 tag. Removing it
	fStats, err := f.Stat()
	if err != nil {
		return fmt.Errorf("cannot get file stats: %s", err)
	}

	err = f.Truncate(fStats.Size() - int64(ID3v1SIZE))
	if err != nil {
		return fmt.Errorf("could not truncate file %s", err)
	}

	// writing new tag
	err = tag.Write(f)
	if err != nil {
		return fmt.Errorf("could not write to writer: %s", err)
	}

	return nil
}
