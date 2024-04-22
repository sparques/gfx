package gfx

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"os"
	"testing"
)

func Test_SoftScreen(t *testing.T) {
	rand.Seed(0)
	screen := NewSoftScreen(image.Rect(0, 0, 8, 16), image.Rect(0, 0, 240, 128), image.Rect(0, 0, 240, 128))
	// screen := NewSoftScreenOf[RGB565BE](image.Rect(0, 0, 8, 16), image.Rect(0, 0, 240, 128), image.Rect(0, 0, 240, 128))
	// screen.Convert = RGB565BEModel

	for y := 0; y < 16*20; y++ {
		for x := 0; x < 240; x++ {
			screen.Set(x, y, color.RGBA{uint8(y * 256 / (135)), uint8(y * 256 / 135), uint8(y * 256 / 135), 255})
		}
	}

	for r := 0; r < 9; r++ {
		for c := 0; c < 40; c++ {
			screen.Fill(screen.CellAt(c, r).Bounds(), randomColor())
			/*
				draw.Draw(
					screen.CellAt(c, r),
					screen.CellAt(c, r).Bounds(),
					image.NewUniform(randomColor()),
					// invertColors{screen.CellAt(c, r)},
					screen.CellAt(c, r).Bounds().Min,
					draw.Src)
			*/

		}
	}

	for i := 0; i < 360; i++ {
		save(fmt.Sprintf("grid-frame-%03d.png", i), screen)
		//screen.Scroll(0, )
		screen.SetViewport(image.Pt(int(240*math.Cos(float64(i)*math.Pi/180)), int(135*math.Sin(float64(i)*math.Pi/180))))
	}
}

func randomColor() color.Color {
	return color.RGBA{uint8(rand.Intn(256)), uint8(rand.Intn(256)), uint8(rand.Intn(256)), 255}
}

func save(fname string, img image.Image) {
	fh, err := os.Create(fname)
	if err != nil {
		panic(err)
	}
	png.Encode(fh, img)
	fh.Close()
}
