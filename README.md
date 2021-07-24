# ∙ ID3ED (ID3 - Encoder - Decoder)
## ⚬ Library for encoding/decoding ID3 tags

---

# Status of the package

**ID3v1**. can:

- read
- write

**ID3v1 Enhanced**. can:

- read
- write

**ID3v2**. can: 

- read
- write

---

# ∙ Installation 

To download package source:
```
go get github.com/Unbewohnte/id3ed/...
```
will be deprecated in 1.17 as far as I know

To donwload and compile package:
```
go install github.com/Unbewohnte/id3ed/...
```

---

# ∙ Example

## ⚬ Read/Write v1&&v2 tags

```
package main

import(
    "fmt"
    "github.com/Unbewohnte/id3ed"
    id3v1 "github.com/Unbewohnte/id3ed/v1"
    id3v2 "github.com/Unbewohnte/id3ed/v2"
)

func main() {
    f, err := id3ed.Open("/path/to/mp3/myMP3.mp3")
    if err != nil {
        panic(err)
    }
    
    // reading and writing tags is easy and done by one struct !

    // print all v2 frames
    if f.ContainsID3v2 {
        for_, frame := range f.ID3v2Tag.Frames{
           fmt.Printf("%s\n", frame.Header.ID)
        }
    }

    // create your own v2 tag from custom frames
    myframe, _ := id3v2.NewFrame("TXXX", []byte("Very important data (ᗜˬᗜ)") ,true)
    v2tag := id3v2.NewTAG([]id3v2.Frame{*myframe})

    // replace v2 tag of the file
    f.WriteID3v2(v2tag)
    


    // create your own ID3v1 tag 
    tag := &id3v1.ID3v1Tag{
            SongName: "mysong",
            Artist:   "me",
            Album:    "my album",
            Year:     2021,
            Comment:  "Cool song",
            Track:    1,
            Genre:    "Christian Gangsta Rap", // for list of genres see: "./v1/genres.go"
	}

    // write your v1 tag
    err = f.WriteID3v1(tag)
    if err != nil {
        panic(err)
    }
    
    // etc. 
}
```

---

# ∙ Testing

To test everything
```
go test ./...
```
or
```
go test -v ./...
```
to get a verbose output

OR

```
go test ./package_name_here
```
to test a specific package

---

# ∙ License

[MIT LICENSE](https://github.com/Unbewohnte/id3ed/blob/main/LICENSE)

# ∙ Note

This is **NOT** a fully tested and it is **NOT** a flawlessly working and edge-cases-covered package.

I work on it alone and I am **NOT** a professional who knows what he does.

Please, use with caution !