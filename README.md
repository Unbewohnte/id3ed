# ID3ED (ID3 - Encoder - Decoder)
## Library for encoding/decoding ID3 tags

# Installation 

```
go get github.com/Unbewohnte/id3ed
```

# Usage

## Decoding ID3v1.1
```
package main

import(
    "fmt"
    "github.com/Unbewohnte/id3ed"
)

func main() {
    mp3file, err := os.Open("/path/to/mp3/myMP3.mp3")
    if err != nil {
        panic(err)
    }

    // extract ID3v1.1 tags 
    mp3tags, err := GetID3v11Tags(mp3file)
    if err != nil {
       panic(err)
    }

    // print all tags
    fmt.Printf("%+v",mp3tags)

    // get a certain tag from "getter" function
    songname := mp3tags.GetSongName()
    
    // get a certain tag from struct field
    genre := mp3tags.Genre

    // etc.
}
```

## Encoding ID3v1.1
```
```

# Testing

```
go test
```
or
```
go test -v
```
to get a verbose output

# Under construction !
## Bugs are a possibility rn in this state, the package is still not tested properly 