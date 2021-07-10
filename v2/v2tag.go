package v2

import "github.com/Unbewohnte/id3ed/util"

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
		return util.ToStringLossy(tag.GetFrame("TT2").Contents)
	default:
		return util.ToStringLossy(tag.GetFrame("TIT2").Contents)
	}
}

// Returns the contents for the album frame
func (tag *ID3v2Tag) Album() string {
	switch tag.Header.Version {
	case V2_2:
		return util.ToStringLossy(tag.GetFrame("TAL").Contents)
	default:
		return util.ToStringLossy(tag.GetFrame("TALB").Contents)
	}
}

// Returns the contents for the artist frame
func (tag *ID3v2Tag) Artist() string {
	switch tag.Header.Version {
	case V2_2:
		return util.ToStringLossy(tag.GetFrame("TP1").Contents)
	default:
		return util.ToStringLossy(tag.GetFrame("TPE1").Contents)
	}
}

// Returns the contents for the year frame
func (tag *ID3v2Tag) Year() string {
	switch tag.Header.Version {
	case V2_2:
		return util.ToStringLossy(tag.GetFrame("TYE").Contents)
	default:
		return util.ToStringLossy(tag.GetFrame("TYER").Contents)
	}
}

// Returns the contents for the comment frame
func (tag *ID3v2Tag) Comment() string {
	switch tag.Header.Version {
	case V2_2:
		return util.ToStringLossy(tag.GetFrame("COM").Contents)
	default:
		return util.ToStringLossy(tag.GetFrame("COMM").Contents)
	}
}
