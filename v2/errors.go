package v2

// Exported ID3v2-specific errors

import "fmt"

var ErrDoesNotUseID3v2 error = fmt.Errorf("does not use ID3v2")
var ErrGotPadding error = fmt.Errorf("got padding")
var ErrReadMoreThanSize error = fmt.Errorf("read more bytes than size of the whole tag")
var ErrInvalidID error = fmt.Errorf("invalid identifier")
