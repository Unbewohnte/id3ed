package v2

import "strings"

type ID3v2Tag struct {
	Header  Header
	Frames  []Frame
	Padding uint32
}

// Creates a new v2 tag from given created frames
func NewTAG(frames []Frame) *ID3v2Tag {
	var newtag ID3v2Tag

	header := newHeader(frames)

	newtag.Header = *header
	newtag.Frames = frames

	return &newtag
}

// Searches for frame with the same identifier as id in tag,
// returns &it if found
func (tag *ID3v2Tag) GetFrame(id string) *Frame {
	for _, frame := range tag.Frames {
		if strings.EqualFold(frame.Header.ID(), id) {
			return &frame
		}
	}
	return nil
}

// Checks if a frame with given id exists
func (tag *ID3v2Tag) FrameExists(id string) bool {
	for _, frame := range tag.Frames {
		if strings.EqualFold(frame.Header.ID(), id) {
			return true
		}
	}
	return false
}

// Returns the contents for the title frame
func (tag *ID3v2Tag) Title() string {
	switch tag.Header.Version() {
	case V2_2:
		if !tag.FrameExists("TT2") {
			return ""
		}
		return tag.GetFrame("TT2").Text()
	default:
		if !tag.FrameExists("TIT2") {
			return ""
		}
		return tag.GetFrame("TIT2").Text()
	}
}

// Returns the contents for the album frame
func (tag *ID3v2Tag) Album() string {
	switch tag.Header.Version() {
	case V2_2:
		if !tag.FrameExists("TAL") {
			return ""
		}
		return tag.GetFrame("TAL").Text()
	default:
		if !tag.FrameExists("TALB") {
			return ""
		}
		return tag.GetFrame("TALB").Text()
	}
}

// Returns the contents for the artist frame
func (tag *ID3v2Tag) Artist() string {
	switch tag.Header.Version() {
	case V2_2:
		if !tag.FrameExists("TP1") {
			return ""
		}
		return tag.GetFrame("TP1").Text()
	default:
		if !tag.FrameExists("TPE1") {
			return ""
		}
		return tag.GetFrame("TPE1").Text()
	}
}

// Returns the contents for the year frame
func (tag *ID3v2Tag) Year() string {
	switch tag.Header.Version() {
	case V2_2:
		if !tag.FrameExists("TYE") {
			return ""
		}
		return tag.GetFrame("TYE").Text()
	default:
		if !tag.FrameExists("TYER") {
			return ""
		}
		return tag.GetFrame("TYER").Text()
	}
}

// Returns the contents for the comment frame
func (tag *ID3v2Tag) Comment() string {
	switch tag.Header.Version() {
	case V2_2:
		if !tag.FrameExists("COM") {
			return ""
		}
		return tag.GetFrame("COM").Text()
	default:
		if !tag.FrameExists("COMM") {
			return ""
		}
		return tag.GetFrame("COMM").Text()
	}
}

// Returns raw bytes of embed picture
func (tag *ID3v2Tag) Picture() []byte {
	switch tag.Header.Version() {
	case V2_2:
		if !tag.FrameExists("PIC") {
			return nil
		}
		return tag.GetFrame("PIC").Contents
	default:
		if !tag.FrameExists("APIC") {
			return nil
		}
		return tag.GetFrame("APIC").Contents
	}
}
