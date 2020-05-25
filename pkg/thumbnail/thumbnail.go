package thumbnail

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/jpeg"
	"io/ioutil"
	"net/http"

	"golang.org/x/image/draw"
)

type Thumbnail struct {
	Height int    `json:"height"`
	PreSrc string `json:"preSrc"`
	Src    string `json:"src"`
	Width  int    `json:"width"`
}

func New(id string, size string) (*Thumbnail, error) {
	src := "https://i.ytimg.com/vi/" + id + "/" + size + "default.jpg"
	res, err := http.Get(src)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	img, _, err := image.Decode(res.Body)
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	newRect := image.Rect(0, 0, 10, 10*height/width)
	newImg := image.NewRGBA(newRect)
	draw.BiLinear.Scale(newImg, newImg.Bounds(), img, bounds, draw.Over, nil)

	buf := bytes.NewBuffer(nil)
	encoder := base64.NewEncoder(base64.StdEncoding, buf)
	defer encoder.Close()

	err = jpeg.Encode(encoder, newImg, &jpeg.Options{Quality: 70})
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(buf)
	if err != nil {
		return nil, err
	}
	preSrc := "data:image/jpeg;base64," + string(bytes)

	t := &Thumbnail{
		Height: height,
		PreSrc: preSrc,
		Src:    src,
		Width:  width,
	}

	return t, nil
}
