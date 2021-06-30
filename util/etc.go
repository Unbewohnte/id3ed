package util

import (
	"bytes"
	"encoding/binary"
)

// Returns found key (int) in provided map by value (string);
// If key does not exist in map - returns -1
func GetKey(mp map[int]string, givenValue string) int {
	for key, value := range mp {
		if value == givenValue {
			return key
		}
	}
	return -1
}

// Decodes given integer bytes into integer
func BytesToInt(gBytes []byte) (int64, error) {
	buff := bytes.NewBuffer(gBytes)
	integer, err := binary.ReadVarint(buff)
	if err != nil {
		return 0, err
	}
	buff = nil
	return integer, nil
}
