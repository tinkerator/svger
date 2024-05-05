package svger

import (
	"zappem.net/pub/graphics/svger/mtransform"
	mt "zappem.net/pub/graphics/svger/mtransform"
)

// Circle is an SVG circle element
type Circle struct {
	ID        string  `xml:"id,attr"`
	Transform string  `xml:"transform,attr"`
	Style     string  `xml:"style,attr"`
	Cx        float64 `xml:"cx,attr"`
	Cy        float64 `xml:"cy,attr"`
	Radius    float64 `xml:"r,attr"`
	Fill      string  `xml:"fill,attr"`
	Stroke    string  `xml:"stroke,attr"`

	transform mtransform.Transform
	group     *Group
}

// ParseDrawingInstructions implements the DrawingInstructionParser
// interface
func (c *Circle) ParseDrawingInstructions() chan *DrawingInstruction {
	if c.Fill == "" && c.group.Fill != "" {
		c.Fill = c.group.Fill
	}
	if c.Stroke == "" && c.group.Stroke != "" {
		c.Stroke = c.group.Stroke
	}
	scale := 1.0
	if c.group.Owner != nil {
		scale = c.group.Owner.scale
	}
	pdp := newPathDParse()
	circTransform := mt.Identity()
	if c.Transform != "" {
		if ct, err := parseTransform(c.Transform); err == nil {
			circTransform = ct
		}
	}
	pdp.transform = mt.MultiplyTransforms(pdp.transform, *c.group.Transform)
	pdp.transform = mt.MultiplyTransforms(pdp.transform, circTransform)

	draw := make(chan *DrawingInstruction)
	go func() {
		defer close(draw)

		x, y := pdp.transform.Apply(c.Cx, c.Cy)
		r := scale * c.Radius
		s := scale * c.group.StrokeWidth

		draw <- &DrawingInstruction{
			Kind:   CircleInstruction,
			M:      &Tuple{x, y},
			Radius: &r,
		}
		draw <- &DrawingInstruction{
			Kind:        PaintInstruction,
			StrokeWidth: &s,
			Stroke:      &c.Stroke,
			Fill:        &c.Fill,
		}
	}()

	return draw
}
