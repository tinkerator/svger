package svger

import (
	"zappem.net/pub/graphics/svger/mtransform"
	mt "zappem.net/pub/graphics/svger/mtransform"
)

// Rect is an SVG XML rect element
type Rect struct {
	ID          string  `xml:"id,attr"`
	Width       float64 `xml:"width,attr"`
	Height      float64 `xml:"height,attr"`
	Transform   string  `xml:"transform,attr"`
	Style       string  `xml:"style,attr"`
	X           float64 `xml:"x,attr"`
	Y           float64 `xml:"y,attr"`
	Fill        string  `xml:"fill,attr"`
	Stroke      string  `xml:"stroke,attr"`
	StrokeWidth float64 `xml:"stroke-width,attr"`

	transform mtransform.Transform
	group     *Group
}

// ParseDrawingInstructions implements the DrawingInstructionParser
// interface
func (r *Rect) ParseDrawingInstructions() chan *DrawingInstruction {
	scale := 1.0
	if r.group == nil {
		r.group = new(Group)
		temp := mt.Identity()
		r.group.Transform = &temp
	} else {
		if r.Fill == "" && r.group.Fill != "" {
			r.Fill = r.group.Fill
		}
		if r.Stroke == "" && r.group.Stroke != "" {
			r.Stroke = r.group.Stroke
		}
		if r.StrokeWidth == 0 && r.group.StrokeWidth != 0 {
			r.StrokeWidth = r.group.StrokeWidth
		}
		if r.group.Owner != nil {
			scale = r.group.Owner.scale
		}
	}
	pdp := newPathDParse()
	rectTransform := mt.Identity()
	if r.Transform != "" {
		if rt, err := parseTransform(r.Transform); err == nil {
			rectTransform = rt
		}
	}
	pdp.transform = mt.MultiplyTransforms(pdp.transform, *r.group.Transform)
	pdp.transform = mt.MultiplyTransforms(pdp.transform, rectTransform)

	draw := make(chan *DrawingInstruction)
	go func() {
		defer close(draw)

		for i, pt := range []struct{ x, y float64 }{
			{r.X, r.Y},
			{r.X + r.Width, r.Y},
			{r.X + r.Width, r.Y + r.Height},
			{r.X, r.Y + r.Height},
		} {
			k := LineInstruction
			if i == 0 {
				k = MoveInstruction
			}
			x, y := pdp.transform.Apply(pt.x, pt.y)
			draw <- &DrawingInstruction{
				Kind: k,
				M:    &Tuple{x, y},
			}
		}
		draw <- &DrawingInstruction{
			Kind: CloseInstruction,
		}

		s := scale * r.StrokeWidth
		draw <- &DrawingInstruction{
			Kind:        PaintInstruction,
			StrokeWidth: &s,
			Stroke:      &r.Stroke,
			Fill:        &r.Fill,
		}
	}()
	return draw
}
