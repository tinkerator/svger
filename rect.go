package svg

import (
	mt "zappem.net/pub/graphics/svg/mtransform"
)

// Rect is an SVG XML rect element
type Rect struct {
	ID        string `xml:"id,attr"`
	Width     string `xml:"width,attr"`
	Height    string `xml:"height,attr"`
	Transform string `xml:"transform,attr"`
	Style     string `xml:"style,attr"`
	Rx        string `xml:"rx,attr"`
	Ry        string `xml:"ry,attr"`

	transform mt.Transform
	group     *Group
}

// ParseDrawingInstructions implements the DrawingInstructionParser
// interface
func (r *Rect) ParseDrawingInstructions() (chan *DrawingInstruction, chan error) {
	draw := make(chan *DrawingInstruction)
	errs := make(chan error)

	defer close(draw)
	defer close(errs)

	return draw, errs
}
