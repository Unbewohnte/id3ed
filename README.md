# ∙ ID3ED (ID3 - Encoder - Decoder)
## ⚬ Library for encoding/decoding ID3 tags

---
  
# Under construction !

# Project status

Right now it`s capable of reading and writing ID3v1 and ID3v1.1 tags.

ID3v2 support is still in making, but it can read header and frames

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

# ∙ Examples

## ⚬ Reading ID3v1
```
package main

import(
    "fmt"
    id3v1 "github.com/Unbewohnte/id3ed/v1"
)

func main() {
    mp3file, err := os.Open("/path/to/mp3/myMP3.mp3")
    if err != nil {
        panic(err)
    }

    // extract ID3v1.1 tag 
    mp3tag, err := id3v1.Readv1Tag(mp3file)
    if err != nil {
       panic(err)
    }

    // print the whole tag
    fmt.Printf("%+v",mp3tag)

    // get certain fields
    songname := mp3tag.SongName

    genre := mp3tag.Genre

    // etc.
}
```

## ⚬ Writing ID3v1
```
package main

import(
    "fmt"
    id3v1 "github.com/Unbewohnte/id3ed/v1"
)

func main() {
	f, err := os.OpenFile("/path/to/file/myfile.mp3",os.O_WRONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer f.Close()

    // create your tags struct
	tag := &id3v1.ID3v1Tag{
            SongName: "mysong",
            Artist:   "me",
            Album:    "my album",
            Year:     2021,
            Comment:  "Cool song",
            Track:    1,
            Genre:    "Christian Gangsta Rap", // for list of genres see: "./v1/genres.go"
	}

    // write tags to file
	err = tag.WriteToFile(f)
	if err != nil {
		panic(err)
	}
}
```

## ⚬ Reading ID3v2
```
package main

import(
    "fmt"
    id3v2 "github.com/Unbewohnte/id3ed/v2"
)

func main() {
    mp3file, err := os.Open("/path/to/mp3/myMP3.mp3")
    if err != nil {
        panic(err)
    }

    // extract ID3v2 tag 
    mp3tag, err := id3v2.Readv2Tag(mp3file)
    if err != nil {
       panic(err)
    }

    // you can get some essential frames` contents from getters
    title := tag.Title() // string
    
    // or get direct frame by id
    commentFrame := tag.GetFrame("COMM") // *Frame 
    // commentFrame.(Header.(...)|Contents)

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