package gfx

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"
)

func Test_ScaledRGBA(t *testing.T) {
	fh, err := os.Open("in.png")
	if err != nil {
		panic(err)
	}

	img, _, err := image.Decode(fh)
	if err != nil {
		panic(err)
	}

	fh.Close()

	scaled := NewScaledRGBA(NewRGBA(img.(*image.RGBA)), 2)

	scaled.Set(1, 1, color.White)

	scaled.Flush()

	fh, _ = os.Create("scaled.png")
	err = png.Encode(fh, img)
	if err != nil {
		panic(err)
	}
}
