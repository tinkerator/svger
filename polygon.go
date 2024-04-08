package svger

import (
	"zappem.net/pub/graphics/svger/mtransform"
)

// Polygon is a closed shape of straight line segments
type Polygon struct {
	ID        string `xml:"id,attr"`
	Transform string `xml:"transform,attr"`
	Style     string `xml:"style,attr"`
	Points    string `xml:"points,attr"`

	transform mtransform.Transform
	group     *Group
}
