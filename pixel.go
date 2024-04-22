package gfx

import (
	"image/color"
	"math/bits"
)

// RGB565 as 'borrowed' from github.com/tinygo-org/drivers/pixel/pixel.go
//
// RGB565 as used in many SPI displays. Stored as a big endian value.
//
// The color format in integer form is gggbbbbb_rrrrrggg on little endian
// systems, which is the standard RGB565 format but with the top and bottom
// bytes swapped.
//
// There are a few alternatives to this weird big-endian format, but they're not
// great:
//   - Storing the value in two 8-bit stores (to make the code endian-agnostic)
//     incurs too much of a performance penalty.
//   - Swapping the upper and lower bits just before storing. This is still less
//     efficient than it could be, since colors are usually constructed once and
//     then reused in many store operations. Doing the swap once instead of many
//     times for each store is a performance win.
type RGB565BE uint16

var RGB565BEModel = color.ModelFunc(rgb565BEModelFunc)

func rgb565BEModelFunc(c color.Color) color.Color {
	r, g, b, _ := c.RGBA()

	val := uint16((r/0x101)&0xF8)<<8 +
		uint16((g/0x101)&0xFC)<<3 +
		uint16((b/0x101)&0xF8)>>3
	// Swap endianness (make big endian).
	// This is done using a single instruction on ARM (rev16).
	// TODO: this should only be done on little endian systems, but TinyGo
	// doesn't currently (2023) support big endian systems so it's difficult to
	// test. Also, big endian systems don't seem fasionable these days.
	val = bits.ReverseBytes16(val)
	return RGB565BE(val)
}

func NewRGB565BE(r, g, b uint8) RGB565BE {
	val := uint16(r&0xF8)<<8 +
		uint16(g&0xFC)<<3 +
		uint16(b&0xF8)>>3
	// Swap endianness (make big endian).
	// This is done using a single instruction on ARM (rev16).
	// TODO: this should only be done on little endian systems, but TinyGo
	// doesn't currently (2023) support big endian systems so it's difficult to
	// test. Also, big endian systems don't seem fasionable these days.
	val = bits.ReverseBytes16(val)
	return RGB565BE(val)
}

func (c RGB565BE) BitsPerPixel() int {
	return 16
}

// Convert makes RGB565BE implement color.ColorModel
func (c RGB565BE) Convert(in color.Color) color.Color {
	return rgb565BEModelFunc(in)
}

func (c RGB565BE) RGBA() (r, g, b, a uint32) {
	// Note: on ARM, the compiler uses a rev instruction instead of a rev16
	// instruction. I wonder whether this can be optimized further to use rev16
	// instead?
	c = c<<8 | c>>8

	r = uint32(uint8(c>>11) << 3)
	g = uint32(uint8(c>>5) << 2)
	b = uint32(uint8(c&0xff) << 3)
	a = 255 * 0x101

	// Correct color rounding, so that 0xff roundtrips back to 0xff.

	r |= r >> 5
	g |= g >> 6
	b |= b >> 5

	// RGBA() returns values scalled to 16bits, so have to multiply by 0x101
	r *= 0x101
	g *= 0x101
	b *= 0x101

	return
}
