package gfx

import (
	"image"
	"image/color"
	"image/draw"
)

// SoftScreen is a SoftScreenOf[] using a slice of color.Colors for pixel data.
// This is horribly inefficient, but very flexible.
type SoftScreen = SoftScreenOf[color.Color]

// SoftScreen implements draw.Image and has some helper methods to make using it as a display
// more convenient. It wraps its x and y coordinates so it acts like a ring-buffer: you can continuously
// pan/scroll the viewport.
type SoftScreenOf[PixType color.Color] struct {
	Pix []PixType
	// Canvas is the full image buffer.
	Canvas image.Rectangle
	// Viewport is what sub-section of Canvas is displayed. This image.Rectange is what is returned
	// when SoftScreen.Bounds() is called and shifting it around is how
	Viewport image.Rectangle

	Cell image.Rectangle

	Convert color.Model
}

func NewSoftScreen(cell, viewport, canvas image.Rectangle) *SoftScreen {
	//TODO: add check so that canvas and viewport are a multiple of cell

	pix := make([]color.Color, canvas.Dx()*canvas.Dy())
	for i := range pix {
		pix[i] = color.Alpha{0}
	}
	return &SoftScreen{
		Pix:      pix,
		Canvas:   canvas,
		Viewport: viewport,
		Cell:     cell,
		Convert:  color.RGBAModel,
	}
}

func NewSoftScreenOf[PixType color.Color](cell, viewport, canvas image.Rectangle) *SoftScreenOf[PixType] {
	//TODO: add check so that canvas and viewport are a multiple of cell

	pix := make([]PixType, canvas.Dx()*canvas.Dy())
	// for i := range pix {
	// 	pix[i] =
	// }

	return &SoftScreenOf[PixType]{
		Pix:      pix,
		Canvas:   canvas,
		Viewport: viewport,
		Cell:     cell,
	}
}

// At implements image.Image interface and will wrap so that reads before the beginning and
// after the end always work.
func (s *SoftScreenOf[PixType]) At(x, y int) color.Color {
	pt := image.Pt(x, y).Mod(s.Canvas)
	/*
		y = (y % s.Canvas.Dy())
		if y < 0 {
			y += s.Canvas.Dy()
		}
		x = (x % s.Canvas.Dx())
		if x < 0 {
			x += s.Canvas.Dx()
		}
	*/
	return s.Pix[pt.Y*s.Canvas.Dx()+pt.X]
}

// Set implements draw.Image interface and is set to wrap so that writes before the beginning and
// after the end always work.
func (s *SoftScreenOf[PixType]) Set(x, y int, c color.Color) {
	// Is below better than image.Pt(x,y).Mod(s.Canvas) ?
	y = (y % s.Canvas.Dy())
	if y < 0 {
		y += s.Canvas.Dy()
	}
	x = (x % s.Canvas.Dx())
	if x < 0 {
		x += s.Canvas.Dx()
	}

	if native, ok := c.(PixType); ok {
		s.Pix[y*s.Canvas.Dx()+x] = native
		return
	}

	// otherwise attempt conversion;
	// Silently fail. Is there not a better option?
	if s.Convert == nil {
		return
	}
	if native, ok := s.Convert.Convert(c).(PixType); ok {
		s.Pix[y*s.Canvas.Dx()+x] = native
	}

}

func (s *SoftScreenOf[PixType]) Bounds() image.Rectangle {
	return s.Viewport
}

func (s *SoftScreenOf[PixType]) ColorModel() color.Model {
	return s.Convert
}

// CellAt returns a cell-sized section of SoftScreen at the given row and column
func (s *SoftScreenOf[PixType]) CellAt(c, r int) draw.Image {
	return cell[PixType]{
		SoftScreenOf: s,
		bounds:       s.Cell.Add(image.Point{s.Cell.Dx() * c, s.Cell.Dy() * r}),
	}
}

func (s *SoftScreenOf[PixType]) Scroll(amount int) {
	s.Pan(0, amount)
}

// Pan shifts the viewport around. x and y are in pixels.
func (s *SoftScreenOf[PixType]) Pan(x, y int) {
	s.Viewport = s.Viewport.Add(image.Point{x, y})

	// return if scrolling has not caused Viewport to stop intersecting with Canvas
	if !s.Viewport.Overlaps(s.Canvas) {
		return
	}

	// If we got to here, that means we need to adjust Viewport so it's back to overlapping
	// with Canvas. This doesn't matter for displaying the right pixels, as SoftScreen will wrap
	// any out-of-bound requests, but to keep our infinite pan/scroll infinite, we need to reset
	// the position of Viewport so we don't hit integer wrapping.

	// If Viewport's lower X is past Canvas's Maximum X, shift viewport back to
	// begining of Canvas

	s.SetViewport(s.Viewport.Min.Mod(s.Canvas))

	/*
		vDx, vDy := s.Viewport.Dx(), s.Viewport.Dy()
		s.Viewport.Min = s.Viewport.Min.Mod(s.Canvas)
		//s.Viewport.Min.X = ((s.Viewport.Min.X - s.Canvas.Min.X) % s.Canvas.Dx()) + s.Canvas.Min.X
		//s.Viewport.Min.Y = ((s.Viewport.Min.Y - s.Canvas.Min.Y) % s.Canvas.Dy()) + s.Canvas.Min.Y
		s.Viewport.Max.X = s.Viewport.Min.X + vDx
		s.Viewport.Max.Y = s.Viewport.Min.Y + vDy
	*/
}

func (s *SoftScreenOf[PixType]) SetViewport(pt image.Point) {
	s.Viewport = s.Viewport.Sub(s.Viewport.Min).Add(pt)
}

func (s *SoftScreenOf[PixType]) Blit(src image.Image, at image.Point) {
	blit(s, src, at)
}

func (s *SoftScreenOf[PixType]) Fill(rect image.Rectangle, c color.Color) {
	fill(s, rect, c)
}

type cell[PixType color.Color] struct {
	*SoftScreenOf[PixType]
	bounds image.Rectangle
}

func (c cell[PixType]) Bounds() image.Rectangle {
	return c.bounds
}
