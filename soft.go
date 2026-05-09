package gfx

import (
	"image"
	"image/draw"
)

func SoftVectorScroll(img draw.Image, region image.Rectangle, vector image.Point) {
	region = img.Bounds().Intersect(region)
	var (
		dst, src     image.Point
		xstep, ystep int
	)

	if vector.Y >= 0 {
		dst.Y = region.Min.Y
		ystep = 1
	} else {
		dst.Y = region.Max.Y - 1
		ystep = -1
	}

	if vector.X >= 0 {
		xstep = 1
	} else {
		xstep = -1
	}

	for range region.Dy() {
		if vector.X >= 0 {
			dst.X = region.Min.X
		} else {
			dst.X = region.Max.X - 1
		}
		for range region.Dx() {
			src = dst.Add(vector).Mod(region)
			img.Set(dst.X, dst.Y, img.At(src.X, src.Y))
			dst.X += xstep
		}
		dst.Y += ystep
	}

	return
}

// software implementation of blit
// presumably hardware implementations are faster

// check if images are the same type and have PixOffset()
func blit(dst Drawer, src Image, at image.Point) {

	// Option 1
	/*
		forAllPix(src.Bounds(), func(x, y int) {
			dst.Set(
				x-src.Bounds().Min.X+at.X,
				y-src.Bounds().Min.Y+at.Y,
				src.At(x, y))
		})
	*/

	// Option 2 - Generate a rectangle the size of src, positioned at 'at' and ensure we only
	// operate on valid bits of both (via Intersect())
	// rect := image.Rect(0, 0, src.Bounds().Dx(), src.Bounds().Dy()).Add(at).Intersect(dst.Bounds())
	// srcXOffset := -at.X + src.Bounds().Min.X
	// srcYOffset := -at.Y + src.Bounds().Min.Y
	// forAllPix(rect, func(x, y int) {
	// 	dst.Set(x, y, src.At(x+srcXOffset, y+srcYOffset))
	// })

	// Option 3 - calculate the offsets of xDst and yDst; iterate over all the pixels in src and using
	// xDst and yDst, update pixels in dst.
	xDst := -src.Bounds().Min.X + at.X
	yDst := -src.Bounds().Min.Y + at.Y
	forAllPix(src.Bounds(), func(x, y int) {
		dst.Set(
			x+xDst,
			y+yDst,
			src.At(x, y))
	})
}
