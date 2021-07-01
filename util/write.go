package util

import (
	"bytes"
	"fmt"
	"io"
)

// Writes data to wr, if len(data) is less than lenNeeded - adds null bytes until written lenNeeded bytes
func WriteToExtent(wr io.Writer, data []byte, lenNeeded int) error {
	if len(data) > lenNeeded {
		return fmt.Errorf("length of given data bytes is bigger than length needed")
	}

	buff := new(bytes.Buffer)
	for i := 0; i < lenNeeded; i++ {
		if i < len(data) {
			// write given data
			err := buff.WriteByte(data[i])
			if err != nil {
				return fmt.Errorf("could not write byte: %s", err)
			}
		} else {
			// write null bytes
			err := buff.WriteByte(0)
			if err != nil {
				return fmt.Errorf("could not write byte: %s", err)
			}
		}
	}

	// write constructed buffer`s bytes
	_, err := wr.Write(buff.Bytes())
	if err != nil {
		return fmt.Errorf("could not write to writer: %s", err)
	}

	return nil
}
