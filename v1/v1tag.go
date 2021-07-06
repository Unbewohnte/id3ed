package v1

// https://id3.org/ID3v1 - documentation

type ID3v1Tag struct {
	SongName string
	Artist   string
	Album    string
	Year     int
	Comment  string
	Genre    string
}

func (*ID3v1Tag) Version() int {
	return 10
}
