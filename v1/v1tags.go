package v1

// https://id3.org/ID3v1 - documentation

type ID3v1Tags struct {
	SongName string
	Artist   string
	Album    string
	Year     int
	Comment  string
	Genre    string
}

func (tags *ID3v1Tags) Version() int {
	return 10
}
