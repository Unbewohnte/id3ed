package v1

type ID3v11Tag struct {
	SongName string
	Artist   string
	Album    string
	Year     int
	Comment  string
	Track    int
	Genre    string
}

func (*ID3v11Tag) Version() int {
	return 11
}
