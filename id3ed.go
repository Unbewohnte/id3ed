package id3ed

import (
	"fmt"
	"os"

	v1 "github.com/Unbewohnte/id3ed/v1"
	v2 "github.com/Unbewohnte/id3ed/v2"
)

type File struct {
	path          string
	ContainsID3v1 bool
	ContainsID3v2 bool
	ID3v1Tag      *v1.ID3v1Tag
	ID3v2Tag      *v2.ID3v2Tag
}

// Opens file specified by `path`, reads all ID3 tags that
// it can find and returns a *File
func Open(path string) (*File, error) {
	fhandler, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open: %s", err)
	}
	defer fhandler.Close()

	var file File

	file.path = path // saving path to the file for future writing

	v1tag, err := v1.Readv1Tag(fhandler)
	switch err {
	case nil:
		file.ContainsID3v1 = true
	case v1.ErrDoesNotUseID3v1:
		file.ContainsID3v1 = false
	default:
		return nil, fmt.Errorf("could not read ID3v1 tag from file: %s", err)
	}
	file.ID3v1Tag = v1tag

	v2tag, err := v2.ReadV2Tag(fhandler)
	switch err {
	case nil:
		file.ContainsID3v2 = true
	case v2.ErrDoesNotUseID3v2:
		file.ContainsID3v2 = false
	default:
		return nil, fmt.Errorf("could not read ID3v2 tag from file: %s", err)
	}
	file.ID3v2Tag = v2tag

	return &file, nil
}

// Writes given ID3v1 tag to file
func (f *File) WriteID3v1(tag *v1.ID3v1Tag) error {
	fhandler, err := os.OpenFile(f.path, os.O_RDWR, os.ModePerm)
	if err != nil {
		return fmt.Errorf("could not read a file: %s", err)
	}
	defer fhandler.Close()

	err = tag.WriteToFile(fhandler)
	if err != nil {
		return fmt.Errorf("could not write ID3v1 to file: %s", err)
	}

	return nil
}

// still not implemented
// func (f *File) WriteID3v2(tag *v2.ID3v2Tag) error {
// 	return nil
// }
