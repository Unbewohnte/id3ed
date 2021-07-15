package v1

// Exported ID3v1-specific errors

import "fmt"

var ErrDoesNotUseID3v1 error = fmt.Errorf("does not use ID3v1")
