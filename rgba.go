package gfx

import (
	"image"
	"image/color"
	"image/draw"
)

const (
	// rgbaWidth is how many bytes a single bgra pixel / color takes up.
	rgbaWidth = 4
)

// RGBA wraps an *image.RGBA and implements all the gfx interfaces
type RGBA struct {
	*image.RGBA
	doubleBuf *image.RGBA
}

func NewRGBA(base *image.RGBA) *RGBA {
	return &RGBA{RGBA: base}
}

func NewRGBAWithDoubleBuffer(base *image.RGBA) *RGBA {
	rgba := &RGBA{
		RGBA:      image.NewRGBA(base.Bounds()),
		doubleBuf: base,
	}
	draw.Draw(rgba, base.Bounds(), base, base.Bounds().Min, draw.Over)
	return rgba
}

func (rgba *RGBA) Flush() {
	// TODO: partial updates / dirty rect
	if rgba.doubleBuf != nil {
		copy(rgba.doubleBuf.Pix, rgba.RGBA.Pix)
	}
}

// Scroll implements gfx.Scroller
func (rgba *RGBA) Scroll(amount int) {
	switch {
	case amount == 0:
		return
	case amount > 0:
		if amount > rgba.Rect.Dy() {
			amount = rgba.Rect.Dy()
		}
		copy(rgba.Pix, rgba.Pix[rgba.Stride*amount:])
	case amount < 0:
		amount *= -1
		if amount > rgba.Rect.Dy() {
			amount = rgba.Rect.Dy()
		}
		reverseCopy(rgba.Pix[rgba.Stride*amount:], rgba.Pix[:len(rgba.Pix)-rgba.Stride*amount])
	}
}

// RegionScroll implements gfx.RegionScroller
func (rgba *RGBA) RegionScroll(region image.Rectangle, amount int) {
	region = rgba.Rect.Intersect(region)
	if region.Empty() || amount == 0 {
		return
	}
	// if amount is positive or negative, copy lines forwards or backwards

	var start, end int
	if amount > 0 {
		for y := region.Min.Y; y < (region.Max.Y - amount); y++ {
			start = rgba.Stride*y + region.Min.X*rgbaWidth
			end = rgba.Stride*y + region.Max.X*rgbaWidth

			copy(rgba.Pix[start:end], rgba.Pix[start+amount*rgba.Stride:end+amount*rgba.Stride])
		}
		return
	}

	// negative scrolling (scrolling up)
	for y := region.Max.Y - 1; y >= (region.Min.Y - amount); y-- {
		start = rgba.Stride*y + region.Min.X*rgbaWidth
		end = rgba.Stride*y + region.Max.X*rgbaWidth

		copy(rgba.Pix[start:end], rgba.Pix[start+amount*rgba.Stride:end+amount*rgba.Stride])
	}
}

// VectorScroll implements gfx.VectorScroller.
func (rgba *RGBA) VectorScroll(region image.Rectangle, vector image.Point) {
	region = rgba.Rect.Intersect(region)
	if region.Empty() || vector == (image.Point{}) {
		return
	}
	// The below is a bit verbose (could simplify it with a few if statements inside the for-loops).
	// It is kept verbose for the sake of performance.

	// the long, confusing src = lines work like this
	// newOffset = stride * [ (y+y_offset+height) % height ] + rgbaWidth * [ (x+x_offset+width) % width]
	//
	// (d + d_offset + d_total) % d_total is the majority of the magic. Increment (or decrement) d by d_offset modulo
	// the maximum of our working dimention. The extra addition of the maximum work dimension is to take care of negative
	// offsets.
	var dst, src int
	if vector.Y > 0 {
		if vector.X >= 0 {
			// down and (possibly) to the right
			// start at the top work our way down, start on the left, work our way right
			for y := range region.Dy() {
				for x := range region.Dx() {
					dst = rgba.PixOffset(region.Min.X+x, region.Min.Y+y)
					src = (region.Min.Y+((y+vector.Y+region.Dy())%region.Dy()))*rgba.Stride + (region.Min.X+((x+vector.X+region.Dx())%region.Dx()))*rgbaWidth
					copy(rgba.Pix[dst:dst+rgbaWidth:dst+rgbaWidth], rgba.Pix[src:src+rgbaWidth:src+rgbaWidth])
				}
			}
		} else {
			for y := range region.Dy() {
				for x := range region.Dx() {
					x = region.Dx() - 1 - x
					dst = rgba.PixOffset(region.Min.X+x, region.Min.Y+y)

					src = (region.Min.Y+((y+vector.Y+region.Dy())%region.Dy()))*rgba.Stride + (region.Min.X+((x+vector.X+region.Dx())%region.Dx()))*rgbaWidth
					copy(rgba.Pix[dst:dst+rgbaWidth:src+rgbaWidth], rgba.Pix[src:src+rgbaWidth:src+rgbaWidth])
				}
			}
		}
	} else {
		if vector.X >= 0 {
			// down and (possibly) to the right
			// start at the top work our way down, start on the left, work our way right

			for y := range region.Dy() {
				y = region.Dy() - 1 - y
				for x := range region.Dx() {
					dst = rgba.PixOffset(region.Min.X+x, region.Min.Y+y)

					src = (region.Min.Y+((y+vector.Y+region.Dy())%region.Dy()))*rgba.Stride + (region.Min.X+((x+vector.X+region.Dx())%region.Dx()))*rgbaWidth
					copy(rgba.Pix[dst:dst+rgbaWidth:dst+rgbaWidth], rgba.Pix[src:src+rgbaWidth:src+rgbaWidth])
				}
			}
		} else {
			for y := range region.Dy() {
				y = region.Dy() - 1 - y
				for x := range region.Dx() {
					x = region.Dx() - 1 - x
					dst := rgba.PixOffset(region.Min.X+x, region.Min.Y+y)

					src := (region.Min.Y+((y+vector.Y+region.Dy())%region.Dy()))*rgba.Stride + (region.Min.X+((x+vector.X+region.Dx())%region.Dx()))*rgbaWidth
					copy(rgba.Pix[dst:dst+rgbaWidth:dst+rgbaWidth], rgba.Pix[src:src+rgbaWidth:src+rgbaWidth])
				}
			}
		}
	}
}

// Fill implements gfx.Filler. Whereever p overlaps with 'where', set those
// pixels to color c.
func (rgba *RGBA) Fill(where image.Rectangle, c color.Color) {
	// get c as native color
	nc := color.RGBAModel.Convert(c).(color.RGBA)

	where = rgba.Bounds().Intersect(where)

	if where.Empty() {
		return
	}

	// previously, I tried to be clever and used a maximum-run-length buffer and then
	// copied that to the pix buffer and the below code is just as fast without
	// thrashing memory as much. Go figure.
	var pix []uint8
	for y := where.Min.Y; y < where.Max.Y; y++ {
		for x := where.Min.X; x < where.Max.X; x++ {
			pix = rgba.Pix[rgba.PixOffset(x, y) : rgba.PixOffset(x, y)+rgbaWidth : rgba.PixOffset(x, y)+rgbaWidth]
			pix[0] = nc.R
			pix[1] = nc.G
			pix[2] = nc.B
			pix[3] = nc.A
		}
	}
}

func (rgba *RGBA) Blit(src image.Image, where image.Point) {
	// fast path,
	gfxRGBA, okGFX := src.(*RGBA)
	srcRGBA, okRGBA := src.(*image.RGBA)
	if okGFX || okRGBA {
		var src *image.RGBA
		if okGFX {
			src = gfxRGBA.RGBA
		} else {
			src = srcRGBA
		}

		srcBounds := src.Bounds()
		destRect := image.Rect(0, 0, srcBounds.Dx(), srcBounds.Dy())
		destRect = srcBounds.Intersect(destRect.Add(where))

		for yd, ys := destRect.Min.Y, srcBounds.Min.Y; yd < destRect.Max.Y && ys < srcBounds.Max.Y; {
			for xd, xs := destRect.Min.X, srcBounds.Min.X; xd < destRect.Max.X && xs < srcBounds.Max.X; {
				rgba.Pix[rgba.PixOffset(xd, yd)] = src.Pix[src.PixOffset(xs, ys)]
				xd++
				xs++
			}
			yd++
			ys++
		}

		return
	}

	// slow fall back
	blit(rgba, src, where)
}

func reverseCopy[E any](dst, src []E) {
	for i := len(src) - 1; i >= 0; i-- {
		dst[i] = src[i]
	}
}
