package v2

import (
	"fmt"
	"io"
)

// Writes ID3v2Tag to ws
func (tag *ID3v2Tag) write(ws io.WriteSeeker) error {
	_, err := ws.Seek(0, io.SeekStart)
	if err != nil {
		return fmt.Errorf("could not seek: %s", err)
	}

	// write header
	ws.Write(tag.Header.toBytes())

	// write frames
	for _, frame := range tag.Frames {
		ws.Write(frame.toBytes(tag.Header.Version))
	}

	return nil
}

// func (tag *ID3v2Tag) WriteToFile(f *os.File) error {
// 	return nil
// }
