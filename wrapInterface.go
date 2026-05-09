package gfx

import (
	"image"
	"image/color"
	"image/draw"
)

// IntefaceWrap is a generic way to ensure a draw.Image implements all GFX interfaces.
// If a double buffer is needed as well, use DoubleBuf instead of InterfaceWrap.
type InterfaceWrap struct {
	draw.Image

	vectorScroll func(region image.Rectangle, vector image.Point)
	regionScroll func(region image.Rectangle, amount int)
	scroll       func(amount int)
	fill         func(where image.Rectangle, c color.Color)
	blit         func(image.Image, image.Point)
}

func NewInterfaceWrap(img draw.Image) *InterfaceWrap {
	iw := &InterfaceWrap{
		Image: img,
	}

	if vs, vsok := img.(VectorScroller); vsok {
		iw.vectorScroll = vs.VectorScroll
	} else {
		iw.vectorScroll = func(region image.Rectangle, vector image.Point) {
			SoftVectorScroll(img, region, vector)
		}
	}

	if rs, rsok := img.(RegionScroller); rsok {
		iw.regionScroll = rs.RegionScroll
	} else {
		iw.regionScroll = func(region image.Rectangle, pixAmt int) {
			iw.vectorScroll(region, image.Pt(0, pixAmt))
		}
	}

	if s, sok := img.(Scroller); sok {
		iw.scroll = s.Scroll
	} else {
		iw.scroll = func(amt int) {
			iw.regionScroll(img.Bounds(), amt)
		}
	}

	if fill, ok := img.(Filler); ok {
		iw.fill = fill.Fill
	} else {
		iw.fill = func(r image.Rectangle, c color.Color) {
			draw.Draw(img, r, image.NewUniform(c), r.Min, draw.Src)
		}
	}

	if blitter, ok := img.(Blitter); ok {
		iw.blit = blitter.Blit
	} else {
		iw.blit = func(src image.Image, at image.Point) {
			blit(iw, src, at)
		}
	}

	return iw
}
