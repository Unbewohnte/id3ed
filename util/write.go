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
			err := buff.WriteByte(data[i])
			if err != nil {
				return err
			}
		} else {
			err := buff.WriteByte(0)
			if err != nil {
				return err
			}
		}
	}

	_, err := wr.Write(buff.Bytes())
	if err != nil {
		return err
	}

	return nil
}
