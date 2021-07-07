package v1

// https://id3.org/ID3v1 - documentation

type ID3v1Tag struct {
	Version  string
	SongName string
	Artist   string
	Album    string
	Year     int
	Comment  string
	Track    uint8 // basically a byte, but converted to int for convenience
	Genre    string
}
