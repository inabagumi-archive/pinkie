package scraper

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/jpeg"

	"golang.org/x/image/draw"
)

func scaleImage(img image.Image) image.Image {
	rect := img.Bounds()
	width := rect.Dx()
	height := rect.Dy()

	newRect := image.Rect(0, 0, 10, 10*height/width)
	newImg := image.NewRGBA(newRect)

	draw.BiLinear.Scale(newImg, newImg.Bounds(), img, rect, draw.Over, nil)

	return newImg
}

func encodeToDataURL(img image.Image) (string, error) {
	buf := bytes.NewBuffer(nil)

	err := jpeg.Encode(buf, img, &jpeg.Options{Quality: 70})
	if err != nil {
		return "", err
	}

	data := base64.StdEncoding.EncodeToString(buf.Bytes())

	return "data:image/jpeg;base64," + data, nil
}

type Thumbnail struct {
	Height int    `json:"height"`
	PreSrc string `json:"preSrc"`
	Src    string `json:"src"`
	Width  int    `json:"width"`
}

func NewThumbnail(src string) (*Thumbnail, error) {
	resp, err := Fetch(src)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, err
	}

	lqipImg := scaleImage(img)

	preSrc, err := encodeToDataURL(lqipImg)
	if err != nil {
		return nil, err
	}

	t := &Thumbnail{
		Height: img.Bounds().Dy(),
		PreSrc: preSrc,
		Src:    src,
		Width:  img.Bounds().Dx(),
	}

	return t, nil
}
