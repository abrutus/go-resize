package main

import (
	_ "fmt"
	"strconv"
	"image"
	"image/jpeg"
	_ "image/png"
	"net/http"
	"log"
	_ "io"
	"bytes"
	"io/ioutil"
	"resize"
)
type Img struct {
	Filename string
	Url string
	Content []byte
}

func (p *Img) save() error {
	filename := p.Filename
	return ioutil.WriteFile(filename, p.Content, 0600)
}

func handler(w http.ResponseWriter, r *http.Request) {
	get := func(n string) int {
		i, _  :=  strconv.Atoi(r.FormValue(n))
		return i
	}
	myurl := "https://s3.amazonaws.com/lenses_s3/uploads/IMG_1762.JPG"
	res, err := http.Get(myurl)
	if err != nil {
		log.Fatal(err)
	}
	img, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)

	}
	res.Body.Close()
	rdr := bytes.NewReader(img)
	i, _, err := image.Decode(rdr)
	check(err)
	// Resize if too large, for more efficient moustachioing.
	// We aim for less than 1200 pixels in any dimension; if the
	// picture is larger than that, we squeeze it down to 600.
	max := 0
	if x := get("x"); x != 0 {
		max = x
	} else {
		max = 1200
	}
	if b := i.Bounds(); b.Dx() > max || b.Dy() > max {
		// If it's gigantic, it's more efficient to downsample first
		// and then resize; resizing will smooth out the roughness.
		if b.Dx() > 2*max || b.Dy() > 2*max {
			w, h := max, max
			if b.Dx() > b.Dy() {
				h = b.Dy() * h / b.Dx()
			} else {
				w = b.Dx() * w / b.Dy()
			}
			i = resize.Resample(i, i.Bounds(), w, h)
			b = i.Bounds()
		}
		w, h := max/2, max/2
		if b.Dx() > b.Dy() {
			h = b.Dy() * h / b.Dx()
		} else {
			w = b.Dx() * w / b.Dy()
		}
		i = resize.Resize(i, i.Bounds(), w, h)
	}
	// Deliver image
	w.Header().Set("Content-type", "image/jpeg")
	jpeg.Encode(w, i, nil)
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
	// get file
}
func check(err error) {
	if err != nil {
		panic(err)
	}
}
