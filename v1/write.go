package v1

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/Unbewohnte/id3ed/util"
)

// Writes given ID3v1.0 or ID3v1.1 tag to given io.WriteSeeker.
// NOTE: will not remove already existing ID3v1 tag if it`s present,
// use ⁕WriteToFile⁕ method if you`re working with REAL mp3 files !!!
func (tag *ID3v1Tag) write(dst io.WriteSeeker) error {
	dst.Seek(0, io.SeekEnd)

	// ID
	_, err := dst.Write([]byte(ID3v1IDENTIFIER))
	if err != nil {
		return fmt.Errorf("could not write to writer: %s", err)
	}

	// Song name
	err = util.WriteToExtent(dst, []byte(tag.SongName), 30)
	if err != nil {
		return fmt.Errorf("could not write to writer: %s", err)
	}

	// Artist
	err = util.WriteToExtent(dst, []byte(tag.Artist), 30)
	if err != nil {
		return fmt.Errorf("could not write to writer: %s", err)
	}

	// Album
	err = util.WriteToExtent(dst, []byte(tag.Album), 30)
	if err != nil {
		return fmt.Errorf("could not write to writer: %s", err)
	}

	// Year
	err = util.WriteToExtent(dst, []byte(fmt.Sprint(tag.Year)), 4)
	if err != nil {
		return fmt.Errorf("could not write to writer: %s", err)
	}

	// Comment and Track

	// check for track number, if specified and valid - comment must be shrinked to 28 bytes and 29th
	// byte must be 0 byte (use ID3v1.1 instead of v1.0)
	if tag.Track == 0 {
		// write only 30 bytes long comment without track
		err = util.WriteToExtent(dst, []byte(tag.Comment), 30)
		if err != nil {
			return fmt.Errorf("could not write to writer: %s", err)
		}
	} else {
		// write 28 bytes long shrinked comment
		err = util.WriteToExtent(dst, []byte(tag.Comment), 28)
		if err != nil {
			return fmt.Errorf("could not write to writer: %s", err)
		}

		// write 0 byte as padding
		_, err = dst.Write([]byte{0})
		if err != nil {
			return fmt.Errorf("could not write to writer: %s", err)
		}

		// write track byte
		_, err = dst.Write([]byte{byte(tag.Track)})
		if err != nil {
			return fmt.Errorf("could not write to writer: %s", err)
		}
	}

	// Genre
	genreCode := util.GetKey(id3v1genres, tag.Genre)
	if genreCode == -1 {
		// if no genre found - encode genre code as 255
		genreCode = ID3v1INVALIDGENRE
	}
	genrebyte := make([]byte, 1)
	binary.PutVarint(genrebyte, int64(genreCode))

	_, err = dst.Write(genrebyte)
	if err != nil {
		return fmt.Errorf("could not write to writer: %s", err)
	}

	return nil
}

// Checks for existing ID3v1 or ID3v1.1 tag in file, if present - removes it and replaces with provided tag
func (tag *ID3v1Tag) WriteToFile(f *os.File) error {
	defer f.Close()

	// check for existing ID3v1 tag
	f.Seek(-int64(ID3v1SIZE), io.SeekEnd)

	identifier, err := util.Read(f, 3)
	if err != nil {
		return err
	}

	if !bytes.Equal(identifier, []byte(ID3v1IDENTIFIER)) {
		// no existing identifier, just write given tag
		err = tag.write(f)
		if err != nil {
			return err
		}
		return nil
	}

	// does contain ID3v1 tag. Removing it
	fStats, err := f.Stat()
	if err != nil {
		return fmt.Errorf("cannot get file stats: %s", err)
	}

	err = f.Truncate(fStats.Size() - int64(ID3v1SIZE))
	if err != nil {
		return fmt.Errorf("could not truncate file %s", err)
	}

	// writing a new tag
	err = tag.write(f)
	if err != nil {
		return fmt.Errorf("could not write to writer: %s", err)
	}

	return nil

}
