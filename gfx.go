package gfx // import "github.com/sparques/gfx"

import (
	"image"
	"image/color"
	"image/draw"
)

// Scroller is an interface for implementing scrolling.
type Scroller interface {
	// Scroll scrolls the object by amount pixels. Positve amount
	// scrolls the "screen" down / moves the image up
	Scroll(amount int)
}

// RegionScroller is an interface for scrolling a subsection.
type RegionScroller interface {
	RegionScroll(region image.Rectangle, amount int)
}

// VectorScroller is an interface for scrolling an image based on a vector.
// If you can efficient implement VectorScroller, you can then implement
// RegionScroller and Scroller via VectorScroller.
type VectorScroller interface {
	// VectorScroll scrolls an image by vector.X pixels in x, with postive
	// values shifting the view port to the right, or the image to left, depending
	// on how you want to look at it. Likewise, vector.Y scrolls in y, with
	// postive values moving the viewport down or the image up.
	// How the newly scrolled area is populated is implementation dependent,
	// though simply wrapping the edges is perhaps the easiest way to do this.
	VectorScroll(region image.Rectangle, vector image.Point)
}

// Blitter is an interface for writing a rectangle of pixels at a given offset.
// A display driver might implement the blitter interface using low-level, efficient
// instructions for transfering pixels. The idea behind this is it is faster than
// addressing and updating individual pixels.
type Blitter interface {
	Blit(image.Image, image.Point)
	Bounds() image.Rectangle
}

// Filler is an interface for filling an area with a single color. Some display
// drivers may have dedicated instructions for this, making it very fast. Others
// may simply do the equivalent of Blitting (specifying an area, then streaming
// pixel data). Both of which are faster than addressing and setting individual
// pixels.
type Filler interface {
	Fill(image.Rectangle, color.Color)
	Bounds() image.Rectangle
}

// Drawer not to be confused with draw.Drawer, but rather, draw.Image
type Drawer = draw.Image

// Image is same as image.Image, except we don't care about the ColorModel()
type Image interface {
	At(x, y int) color.Color
	Bounds() image.Rectangle
}

// DoubleBufferer defines an interface such that writes/changes are bufferred
// and Flush() must be called to actually write those changes.
type DoubleBufferer interface {
	Drawer
	Flush()
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

// getRGBAPixels uses a slice of
func getRGBAPixels(img image.Image) []color.RGBA {
	return getRGBAPixelsIn(img, img.Bounds())
}

// getRGBAPixels uses a slice of
func getRGBAPixelsIn(img image.Image, rect image.Rectangle) []color.RGBA {
	pix := make([]color.RGBA, rect.Dx()*rect.Dy())
	i := 0
	forAllPix(rect, func(x, y int) {
		pix[i] = colorToRGBA(img.At(x, y))
		i++
	})

	return pix
}

// software implementation of fill,
// same deal as blit, presumably hardware implementations are faster.
// fill sets all pixels in dst, where they intersect with r, to c.
func fill(dst Drawer, r image.Rectangle, c color.Color) {
	drawArea := dst.Bounds().Intersect(r)
	native := dst.ColorModel().Convert(c)
	forAllPix(drawArea, func(x, y int) {
		dst.Set(x, y, native)
	})
}

// I keep making this set of for-loops
// TODO: Investigate if using a lambda like this hurts performance
func forAllPix(rect image.Rectangle, do func(x, y int)) {
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			do(x, y)
		}
	}
}

func colorToRGBA(c color.Color) color.RGBA {
	if rgba, ok := c.(color.RGBA); ok {
		return rgba
	}
	r, g, b, a := c.RGBA()
	return color.RGBA{uint8(r / 0x101), uint8(g / 0x101), uint8(b / 0x101), uint8(a / 0x101)}
}
