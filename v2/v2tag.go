package v2

type ID3v2Tag struct {
	Header Header
	Frames []Frame
}

// Searches for frame with the same identifier as id in tag,
// returns &it if found
func (tag *ID3v2Tag) GetFrame(id string) *Frame {
	for _, frame := range tag.Frames {
		if frame.Header.ID == id {
			return &frame
		}
	}
	return nil
}

// Returns the contents for the title frame
func (tag *ID3v2Tag) Title() string {
	switch tag.Header.Version {
	case V2_2:
		return tag.GetFrame("TT2").Text()
	default:
		return tag.GetFrame("TIT2").Text()
	}
}

// Returns the contents for the album frame
func (tag *ID3v2Tag) Album() string {
	switch tag.Header.Version {
	case V2_2:
		return tag.GetFrame("TAL").Text()
	default:
		return tag.GetFrame("TALB").Text()
	}
}

// Returns the contents for the artist frame
func (tag *ID3v2Tag) Artist() string {
	switch tag.Header.Version {
	case V2_2:
		return tag.GetFrame("TP1").Text()
	default:
		return tag.GetFrame("TPE1").Text()
	}
}

// Returns the contents for the year frame
func (tag *ID3v2Tag) Year() string {
	switch tag.Header.Version {
	case V2_2:
		return tag.GetFrame("TYE").Text()
	default:
		return tag.GetFrame("TYER").Text()
	}
}

// Returns the contents for the comment frame
func (tag *ID3v2Tag) Comment() string {
	switch tag.Header.Version {
	case V2_2:
		return tag.GetFrame("COM").Text()
	default:
		return tag.GetFrame("COMM").Text()
	}
}
