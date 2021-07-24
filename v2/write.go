package v2

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Unbewohnte/id3ed/util"
)

// Writes ID3v2Tag to ws
func (tag *ID3v2Tag) write(ws io.WriteSeeker) error {
	_, err := ws.Seek(0, io.SeekStart)
	if err != nil {
		return fmt.Errorf("could not seek: %s", err)
	}

	// write header
	_, err = ws.Write(tag.Header.toBytes())
	if err != nil {
		return fmt.Errorf("could not write to writer: %s", err)
	}

	// write frames
	for _, frame := range tag.Frames {
		_, err = ws.Write(frame.toBytes())
		if err != nil {
			return fmt.Errorf("could not write to writer: %s", err)
		}
	}

	// write padding if has any
	if tag.Padding != 0 {
		util.WriteToExtent(ws, []byte{0}, int(tag.Padding))
	}

	return nil
}

// Writes ID3v2Tag to file, removing already existing tag if found
func (tag *ID3v2Tag) WriteToFile(f *os.File) error {
	defer f.Close()

	_, err := f.Seek(0, io.SeekStart)
	if err != nil {
		return fmt.Errorf("could not seek: %s", err)
	}

	// check if there`s content at all
	fStats, err := f.Stat()
	if err != nil {
		return err
	}

	if fStats.Size() < 3 {
		// there`s no way that the file can contain TAG,
		// just write and exit

		// `write` for some reason removes all contents if there`s no tag, so
		// we need forcefully store already existing data and
		// write it again afterwards

		_, err := f.Seek(0, io.SeekStart)
		if err != nil {
			return fmt.Errorf("could not seek: %s", err)
		}

		contents, err := util.Read(f, uint64(fStats.Size()))
		if err != nil {
			return err
		}

		err = tag.write(f)
		if err != nil {
			return err
		}

		_, err = f.Write(contents)
		if err != nil {
			return err
		}

		// apparently, there are 3 zerobytes
		// that appear after writing the contents for some
		// alien-like reason so we need to remove them.
		fStats, err = f.Stat()
		if err != nil {
			return err
		}

		err = f.Truncate(fStats.Size() - 3)
		if err != nil {
			return err
		}

		return nil
	}

	// check for an existing tag
	possibleHeaderID, err := util.ReadToString(f, 3)
	if err != nil {
		return err
	}

	if possibleHeaderID != HEADERIDENTIFIER {
		// No existing tag, just write what we have
		// and exit

		// `write` for some reason removes all contents if there`s no tag, so
		// we need forcefully store already existing data and
		// write it again afterwards

		_, err := f.Seek(0, io.SeekStart)
		if err != nil {
			return fmt.Errorf("could not seek: %s", err)
		}

		contents, err := util.Read(f, uint64(fStats.Size()))
		if err != nil {
			return err
		}

		err = tag.write(f)
		if err != nil {
			return err
		}

		_, err = f.Write(contents)
		if err != nil {
			return err
		}

		// apparently, there are 3 zerobytes
		// that appear after writing the contents for some
		// alien-like reason so we need to remove them.
		fStats, err = f.Stat()
		if err != nil {
			return err
		}

		err = f.Truncate(fStats.Size() - 3)
		if err != nil {
			return err
		}

		return nil
	}
	// there is an existing tag, remove it
	// and write a new one

	// get size of the existing tag
	existingHeader, err := readHeader(f)
	if err != nil {
		return err
	}
	existingHeaderSize := existingHeader.Size()

	// cannot truncate just the existing tag with f.Truncate(),
	// so we need to improvise and have a temporary copy of the mp3,
	// wipe the original file, write our tag and place the actual
	// music without the old tag from the temporary copy.

	// create a temporary file
	temporaryDir := os.TempDir()
	tmpF, err := os.CreateTemp(temporaryDir, fmt.Sprintf("%s_TEMP", filepath.Base(f.Name())))
	if err != nil {
		return err
	}

	defer tmpF.Close()
	// remove it afterwards
	defer os.Remove(filepath.Join(temporaryDir, tmpF.Name()))

	// copy contents of the original mp3 to a temporary one
	_, err = io.Copy(tmpF, f)
	if err != nil {
		return err
	}

	// fully remove contents of the original file
	err = f.Truncate(0)
	if err != nil {
		return err
	}

	// write our tag to the original file, which is at that moment is
	// empty
	err = tag.write(f)
	if err != nil {
		return err
	}

	tmpFStats, err := tmpF.Stat()
	if err != nil {
		return err
	}

	// read all contents of the temporary file, except the existing tag

	musicDataSize := int64(tmpFStats.Size() - int64(existingHeaderSize))

	_, err = tmpF.Seek(int64(existingHeaderSize), io.SeekStart)
	if err != nil {
		return fmt.Errorf("could not seek: %s", err)
	}

	musicData, err := util.Read(tmpF, uint64(musicDataSize))
	if err != nil {
		return err
	}

	// and write them into the original file, which
	// contains only the new tag
	_, err = f.Write(musicData)
	if err != nil {
		return err
	}

	return nil
}
