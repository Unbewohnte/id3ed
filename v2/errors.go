package v2

// Exported ID3v2-specific errors

import "fmt"

var ErrDoesNotUseID3v2 error = fmt.Errorf("does not use ID3v2")
