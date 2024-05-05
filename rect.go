package svger

import (
	"errors"

	"zappem.net/pub/graphics/svger/mtransform"
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

	transform mtransform.Transform
	group     *Group
}

// ErrNoSupportForRect indicates <rect> is not supported yet.
var ErrNoSupportForRect = errors.New("no support for <rect> yet")

// ParseDrawingInstructions implements the DrawingInstructionParser
// interface
func (r *Rect) ParseDrawingInstructions() chan *DrawingInstruction {
	draw := make(chan *DrawingInstruction)
	go func() {
		defer close(draw)
		draw <- &DrawingInstruction{
			Kind:  ErrorInstruction,
			Error: ErrNoSupportForRect,
		}

	}()
	return draw
}
