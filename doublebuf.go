package gfx

import "image"

// DoubleBuf implements a generic doublebuffer for anything implementing the Blitter interface.
// This is not especially efficient, as the pixels are not necesarily the same format and conversions
// must be done every flush.
type DoubleBuf struct {
	*RGBA
	backing Blitter
}

func NewDoubleBuf(img Blitter) *DoubleBuf {
	db := &DoubleBuf{
		RGBA:    NewRGBA(image.NewRGBA(img.Bounds())),
		backing: img,
	}

	return db
}

func (db *DoubleBuf) Flush() {
	db.backing.Blit(db.RGBA.SubImage(db.RGBA.dirty), db.RGBA.dirty.Min)
	db.RGBA.dirty = image.Rectangle{}
}
