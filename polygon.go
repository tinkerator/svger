package svg

import (
	mt "zappem.net/pub/graphics/svg/mtransform"
)

// Polygon is a closed shape of straight line segments
type Polygon struct {
	ID        string `xml:"id,attr"`
	Transform string `xml:"transform,attr"`
	Style     string `xml:"style,attr"`
	Points    string `xml:"points,attr"`

	transform mt.Transform
	group     *Group
}
