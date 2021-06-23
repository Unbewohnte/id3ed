# ID3ED (ID3 - Encoder - Decoder)
## Library for encoding/decoding ID3 tags

# Installation 
`go get https://github.com/Unbewohnte/id3ed`

# Usage

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

    mp3tags, err := GetID3v11Tags(mp3file)
    if err != nil {
       panic(err)
    }

    // printing all tags
    fmt.Printf("%+v",mp3tags)

    // getting certain tag
    songname := mp3tags.GetSongName()

    // etc.
}

```

# Under construction !
## Bugs are a possibility rn in this state