package gfx

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"testing"

	"github.com/sparques/gfx/xform"
)

func Test_Conversion(t *testing.T) {
	// rgb24 := color.RGBA{255, 255, 255, 255}
	// rgb16 := rgb565BEModelFunc(rgb24)
	rand.Seed(0)
	for i := 0; i < 1024; i++ {
		rgb24 := color.RGBA{uint8(rand.Intn(256)), uint8(rand.Intn(256)), uint8(rand.Intn(256)), 255}
		// k, convert it to rgb565be
		rgb16 := rgb565BEModelFunc(rgb24)
		// convert it back and compare
		or, og, ob, _ := rgb24.RGBA()
		nr, ng, nb, _ := rgb16.RGBA()

		if or != nr || og != ng || ob != nb {
			fmt.Println("Mismatch!")
			fmt.Println(rgb24.RGBA())
			fmt.Println(rgb16.RGBA())
		}
	}
}

func Test_Gradient565(t *testing.T) {
	screen := NewSoftScreenOf[RGB565BE](image.Rect(0, 0, 8, 16), image.Rect(0, 0, 240, 256), image.Rect(0, 0, 240, 256))
	screen.Convert = RGB565BEModel

	for y := screen.Bounds().Min.Y; y < screen.Bounds().Max.Y; y++ {
		for x := screen.Bounds().Min.X; x < screen.Bounds().Max.X; x++ {
			screen.Set(x, y, rgb565BEModelFunc(color.RGBA{R: uint8(y), G: 0, B: uint8(y), A: 255}))
		}
	}

	save("gradient-565.png", xform.Blur(screen))
}

func Test_Gradient888(t *testing.T) {
	screen := NewSoftScreen(image.Rect(0, 0, 8, 16), image.Rect(0, 0, 240, 256), image.Rect(0, 0, 240, 256))

	for y := screen.Bounds().Min.Y; y < screen.Bounds().Max.Y; y++ {
		for x := screen.Bounds().Min.X; x < screen.Bounds().Max.X; x++ {
			screen.Set(x, y, color.RGBA{R: uint8(y), G: 0, B: uint8(y), A: 255})
		}
	}

	save("gradient-888.png", screen)
}
