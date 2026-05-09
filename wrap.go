package gfx

import (
	"errors"
	"fmt"
	"image"
	"image/draw"
)

var (
	ErrUnsupported = errors.New("not supported")
)

// Wrap takes an image.Image and attempts to wrap it
// with as a gfx image wrappers, adding support for
// the interfaces defined in this package.
// If the type is not supported, the original image
// is and ErrUnsuported is returned.
func Wrap(img draw.Image) (draw.Image, error) {
	switch i := img.(type) {
	case *image.RGBA:
		return NewRGBA(i), nil
	default:
		return img, fmt.Errorf("%w: %T", ErrUnsupported, img)
	}
}

// WrapDoubleBuf works the same as Wrap, but returns an
// image.Image that implements the DoubleBufferer interface; pixels
// are only updated on a call to Flush()
func WrapDoubleBuf(img draw.Image) (draw.Image, error) {
	switch i := img.(type) {
	case *image.RGBA:
		return NewRGBAWithDoubleBuffer(i), nil
	default:
		return img, fmt.Errorf("%w: %T", ErrUnsupported, img)
	}
}
