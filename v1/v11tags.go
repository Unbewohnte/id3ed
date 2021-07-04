package v1

type ID3v11Tags struct {
	SongName string
	Artist   string
	Album    string
	Year     int
	Comment  string
	Track    int
	Genre    string
}

func (tags *ID3v11Tags) Version() int {
	return 11
}
