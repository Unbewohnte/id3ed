package v1

import (
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
	_, err := dst.Seek(0, io.SeekEnd)
	if err != nil {
		return fmt.Errorf("could not seek: %s", err)
	}

	// write enhanced, if uses one
	if tag.HasEnhancedTag {
		// IDentifier
		err = util.WriteToExtent(dst, []byte(ENHANCEDIDENTIFIER), 4)
		if err != nil {
			return err
		}

		// Songname
		err = util.WriteToExtent(dst, []byte(tag.EnhancedTag.SongName), 60)
		if err != nil {
			return err
		}

		// Artist
		err = util.WriteToExtent(dst, []byte(tag.EnhancedTag.Artist), 60)
		if err != nil {
			return err
		}

		// Album
		err = util.WriteToExtent(dst, []byte(tag.EnhancedTag.Album), 60)
		if err != nil {
			return err
		}

		// Speed
		_, err = dst.Write([]byte(tag.EnhancedTag.Speed))
		if err != nil {
			return err
		}

		// Genre
		err = util.WriteToExtent(dst, []byte(tag.EnhancedTag.Genre), 30)
		if err != nil {
			return err
		}

		// Time
		err = util.WriteToExtent(dst, []byte(tag.EnhancedTag.StartTime), 6)
		if err != nil {
			return err
		}

		err = util.WriteToExtent(dst, []byte(tag.EnhancedTag.EndTime), 6)
		if err != nil {
			return err
		}
	}

	// write a regular ID3v1

	// ID
	_, err = dst.Write([]byte(IDENTIFIER))
	if err != nil {
		return err
	}

	// Song name
	err = util.WriteToExtent(dst, []byte(tag.SongName), 30)
	if err != nil {
		return err
	}

	// Artist
	err = util.WriteToExtent(dst, []byte(tag.Artist), 30)
	if err != nil {
		return err
	}

	// Album
	err = util.WriteToExtent(dst, []byte(tag.Album), 30)
	if err != nil {
		return err
	}

	// Year
	err = util.WriteToExtent(dst, []byte(fmt.Sprint(tag.Year)), 4)
	if err != nil {
		return err
	}

	// Comment and Track

	// check for track number, if specified and valid - comment must be shrinked to 28 bytes and 29th
	// byte must be 0 byte (use ID3v1.1 instead of v1.0)
	if tag.Track == 0 {
		// write only 30 bytes long comment without track
		err = util.WriteToExtent(dst, []byte(tag.Comment), 30)
		if err != nil {
			return err
		}
	} else {
		// write 28 bytes long shrinked comment
		err = util.WriteToExtent(dst, []byte(tag.Comment), 28)
		if err != nil {
			return err
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
		// if no genre found - set genre code to 255
		genreCode = INVALIDGENRE
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

	fStats, err := f.Stat()
	if err != nil {
		return fmt.Errorf("cannot get file stats: %s", err)
	}

	filesize := fStats.Size()

	// process all possible scenarios
	switch {

	case containsEnhancedTAG(f) && containsTAG(f):
		// remove both
		err = f.Truncate(filesize - int64(TAGSIZE) - int64(ENHANCEDSIZE))
		if err != nil {
			return fmt.Errorf("could not truncate file %s", err)
		}
		// write the new one
		err = tag.write(f)
		if err != nil {
			return err
		}

	case containsEnhancedTAG(f) && !containsTAG(f):
		// remove enhanced tag, replace with new
		err = f.Truncate(filesize - int64(ENHANCEDSIZE))
		if err != nil {
			return fmt.Errorf("could not truncate file %s", err)
		}

		err = tag.write(f)
		if err != nil {
			return err
		}

	case !containsEnhancedTAG(f) && containsTAG(f):
		// remove regular one, replace with new
		err = f.Truncate(filesize - int64(TAGSIZE))
		if err != nil {
			return fmt.Errorf("could not truncate file %s", err)
		}

		err = tag.write(f)
		if err != nil {
			return err
		}

	case !containsEnhancedTAG(f) && !containsTAG(f):
		// no existing TAGs, simply write what we have
		err := tag.write(f)
		if err != nil {
			return err
		}
	}

	return nil
}
