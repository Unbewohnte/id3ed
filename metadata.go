package id3ed

import (
	"io"
	"os"
)

// I`m still a bit confused about interfaces,
// I`ll look into them and try to figure out
// how to use them properly
type Metadata interface {
	Read(io.ReadSeeker) error
	Write(io.WriteSeeker) error
	WriteToFile(*os.File) error
}
