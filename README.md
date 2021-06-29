# ∙ ID3ED (ID3 - Encoder - Decoder)
## ⚬ Library for encoding/decoding ID3 tags

---
  
# Under construction !


---

# ∙ Installation 

To download package source:
```
go get github.com/Unbewohnte/id3ed
```

To compile package:
```
go install github.com/Unbewohnte/id3ed
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
    mp3tags, err := id3v1.GetID3v11Tags(mp3file)
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
            Genre:    "Christian Gangsta Rap", // list of genres see "id3v1genres.go"
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