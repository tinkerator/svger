package svg

import (
	"zappem.net/pub/graphics/svg/mtransform"
)

// Ellipse is an SVG ellipse XML element
type Ellipse struct {
	ID        string `xml:"id,attr"`
	Transform string `xml:"transform,attr"`
	Style     string `xml:"style,attr"`
	Cx        string `xml:"cx,attr"`
	Cy        string `xml:"cy,attr"`
	Rx        string `xml:"rx,attr"`
	Ry        string `xml:"ry,attr"`

	transform mtransform.Transform
	group     *Group
}
