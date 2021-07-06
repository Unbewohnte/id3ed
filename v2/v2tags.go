package v2

// type ID3v2Tag struct {
// 	Header Header
// 	Frames []Frame
// }

// type V2TagReader interface {
// 	ReadFrames(io.ReadSeeker) ([]*Frame, error)
// 	GetHeader(io.ReadSeeker) (*Header, error)
// 	HasPadding(io.ReadSeeker) (bool, error)
// }

// type V2TagWriter interface {
// 	Write(*os.File) error
// }

// func Get(f *os.File) (*ID3v2Tag, error) {
// 	var tag ID3v2Tag

// 	header, err := GetHeader(f)
// 	if err != nil {
// 		return nil, err
// 	}
// 	frames, err := GetFrames(f)
// 	if err != nil {
// 		return nil, err
// 	}

// 	tag.Header = header
// 	tag.Frames = frames

// 	return &tag, nil
// }
