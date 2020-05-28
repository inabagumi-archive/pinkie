package thumbnail

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"net/http"

	"golang.org/x/image/draw"
)

func Fetch(url string) (*http.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		resp.Body.Close()

		return nil, fmt.Errorf("thumbnail: invalid status %d", resp.StatusCode)
	}

	return resp, nil
}

func Scale(img image.Image) image.Image {
	rect := img.Bounds()
	width := rect.Dx()
	height := rect.Dy()

	newRect := image.Rect(0, 0, 10, 10*height/width)
	newImg := image.NewRGBA(newRect)

	draw.BiLinear.Scale(newImg, newImg.Bounds(), img, rect, draw.Over, nil)

	return newImg
}

func EncodeToDataURL(img image.Image) (string, error) {
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

func New(id string, size string) (*Thumbnail, error) {
	src := fmt.Sprintf("https://i.ytimg.com/vi/%s/%sdefault.jpg", id, size)

	resp, err := Fetch(src)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, err
	}

	lqipImg := Scale(img)

	preSrc, err := EncodeToDataURL(lqipImg)
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
