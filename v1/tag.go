package v1

// https://id3.org/ID3v1 - documentation

type ID3v1Tag struct {
	version        string
	SongName       string
	Artist         string
	Album          string
	Year           int
	Comment        string
	Track          uint8 // basically a byte, but converted to int for convenience
	Genre          string
	HasEnhancedTag bool
	EnhancedTag    EnhancedID3v1Tag
}

// from https://en.wikipedia.org/wiki/ID3

type EnhancedID3v1Tag struct {
	SongName  string
	Artist    string
	Album     string
	Speed     string
	Genre     string
	StartTime string
	EndTime   string
}

var EnhancedSpeed = map[int]string{
	0: "Unset",
	1: "Slow",
	2: "Medium",
	3: "Fast",
	4: "Hardcore",
}
