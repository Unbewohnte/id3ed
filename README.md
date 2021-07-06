# ∙ ID3ED (ID3 - Encoder - Decoder)
## ⚬ Library for encoding/decoding ID3 tags

---
  
# Under construction !

# Project status

Right now it`s capable of reading and writing ID3v1 and ID3v1.1 tags.

ID3v2.x support is still in making, but it can read header and v.3~v.4 frames

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

# ∙ Usage

## ⚬ Decoding ID3v1.1
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

    // extract ID3v1.1 tags 
    mp3tags, err := id3v1.Getv11Tags(mp3file)
    if err != nil {
       panic(err)
    }

    // print all tags
    fmt.Printf("%+v",mp3tags)

    songname := mp3tags.SongName

    genre := mp3tags.Genre

    // etc.
}
```

## ⚬ Encoding ID3v1.1
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
	tags := &id3v1.ID3v11Tags{
            SongName: "mysong",
            Artist:   "me",
            Album:    "my album",
            Year:     2021,
            Comment:  "Cool song",
            Track:    1,
            Genre:    "Christian Gangsta Rap", // for list of genres see: "./v1/genres.go"
	}

    // write tags to file
	err = tags.WriteToFile(f)
	if err != nil {
		panic(err)
	}
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