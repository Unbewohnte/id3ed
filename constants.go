package id3ed

// ID3v1
const ID3v1IDENTIFIER string = "TAG"
const ID3v1SIZE int = 128 // bytes
const ID3v1INVALIDGENRE int = 255

//ID3v2
const ID3v2IDENTIFIER string = "ID3"
const ID3v2HEADERSIZE int = 10     // bytes
const ID3v2MAXSIZE int = 268435456 // bytes (256 MB)
