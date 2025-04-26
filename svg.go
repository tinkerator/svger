package svger

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"zappem.net/pub/graphics/svger/mtransform"
)

// Debug causes the package to log.Print*() debugging information.
var Debug = false

// Tuple is an X,Y coordinate
type Tuple [2]float64

// Svg represents an SVG file containing at least a top level group or a
// number of Paths
type Svg struct {
	// Title is the title string for the SVG image
	Title string `xml:"title"`
	// Groups lists the top level groups
	Groups []Group `xml:"g"`
	// Width is the width of the SVG image
	Width string `xml:"width,attr"`
	// Height is the height of the SVG image
	Height string `xml:"height,attr"`
	// ViewBox holds the unparsed view box description
	ViewBox string `xml:"viewBox,attr"`
	// Elements lists all of the top level elements in this SVG image
	Elements []DrawingInstructionParser
	// Name names the SVG - typically the filename
	Name string
	// Transform holds the base frame information for the image
	// groups and elements of the image are relative to this
	// base frame.
	Transform *mtransform.Transform
	// scale holds the scaling factor applied to descendant
	// coordinates.
	scale float64
	// instructions is a common channel for emitting the sequence
	// of drawing instructions
	instructions chan *DrawingInstruction
	// errors is a common channel for emitting decoding errors
	errors chan error
}

// Group represents an SVG group (usually located in a 'g' XML element)
type Group struct {
	ID              string
	Stroke          string
	StrokeLineCap   string
	StrokeLineJoin  string
	StrokeWidth     float64
	Fill            string
	FillRule        string
	Elements        []DrawingInstructionParser
	TransformString string
	Transform       *mtransform.Transform // row, column
	Parent          *Group
	Owner           *Svg
	instructions    chan *DrawingInstruction
	errors          chan error
}

// ParseDrawingInstructions implements the DrawingInstructionParser interface
//
// This method makes it easier to get all the drawing instructions.
func (g *Group) ParseDrawingInstructions() chan *DrawingInstruction {
	g.instructions = make(chan *DrawingInstruction, 100)
	go func() {
		defer close(g.instructions)
		for _, e := range g.Elements {
			instrs := e.ParseDrawingInstructions()
			for is := range instrs {
				g.instructions <- is
				if is.Error != nil {
					return
				}
			}
		}
	}()
	return g.instructions
}

// UnmarshalXML implements the encoding.xml.Unmarshaler interface
func (g *Group) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "id":
			g.ID = attr.Value
		case "stroke":
			g.Stroke = attr.Value
		case "stroke-width":
			floatValue, err := strconv.ParseFloat(attr.Value, 64)
			if err != nil {
				return err
			}
			g.StrokeWidth = floatValue
		case "fill":
			g.Fill = attr.Value
		case "fill-rule":
			g.FillRule = attr.Value
		case "transform":
			g.TransformString = attr.Value
			t, err := parseTransform(g.TransformString)
			if err != nil {
				fmt.Println(err)
			}
			g.Transform = &t
		case "style":
			// another way to get some of the above
			suppressStroke := false
			suppressFill := false
			for a, val := range splitStyle(attr.Value) {
				switch a {
				case "fill":
					g.Fill = val
				case "fill-opacity":
					if v := parseDecimal(val); v == 0 {
						suppressFill = true
					}
				case "stroke":
					g.Stroke = val
				case "stroke-linecap":
					g.StrokeLineCap = val
				case "stroke-linejoin":
					g.StrokeLineJoin = val
				case "stroke-opacity":
					if v := parseDecimal(val); v == 0 {
						suppressStroke = true
					}
				case "stroke-width":
					g.StrokeWidth = parseDecimal(val)
				default:
					if Debug {
						log.Printf("TODO ingest style attr %q = %q", a, val)
					}
				}
			}
			if suppressFill {
				g.Fill = "none"
			}
			if suppressStroke {
				g.Stroke = "none"
			}
		default:
			if Debug {
				log.Printf("TODO unable to parse %#v", attr)
			}
		}
	}

	for {
		token, err := decoder.Token()
		if err != nil {
			return err
		}

		switch tok := token.(type) {
		case xml.StartElement:
			var elementStruct DrawingInstructionParser

			switch tok.Name.Local {
			case "g":
				sub := &Group{
					Parent:         g,
					Owner:          g.Owner,
					StrokeLineCap:  g.StrokeLineCap,
					StrokeLineJoin: g.StrokeLineJoin,
					StrokeWidth:    g.StrokeWidth,
					Stroke:         g.Stroke,
					Fill:           g.Fill,
				}
				x := mtransform.MultiplyTransforms(*mtransform.NewTransform(), *g.Transform)
				sub.Transform = &x
				elementStruct = sub
			case "rect":
				rect := &Rect{group: g}
				elementStruct = rect
			case "circle":
				circ := &Circle{group: g}
				elementStruct = circ
			case "path":
				path := &Path{
					group:          g,
					StrokeWidth:    g.StrokeWidth,
					StrokeLineCap:  &g.StrokeLineCap,
					StrokeLineJoin: &g.StrokeLineJoin,
					Stroke:         &g.Stroke,
					Fill:           &g.Fill,
				}
				elementStruct = path
			default:
				if Debug {
					log.Printf("TODO support for %q elements", tok.Name.Local)
				}
				continue
			}
			if err = decoder.DecodeElement(elementStruct, &tok); err != nil {
				return fmt.Errorf("error decoding element of Group: %v", err)
			}
			g.Elements = append(g.Elements, elementStruct)
		case xml.EndElement:
			if tok.Name.Local == "g" {
				return nil
			}
		}
	}
}

// ParseDrawingInstructions implements the DrawingInstructionParser interface
//
// This method makes it easier to get all the drawing instructions.
func (s *Svg) ParseDrawingInstructions() chan *DrawingInstruction {
	s.instructions = make(chan *DrawingInstruction, 100)
	go func() {
		var elecount int
		defer close(s.instructions)
		for _, e := range s.Elements {
			elecount++
			instrs := e.ParseDrawingInstructions()
			for is := range instrs {
				s.instructions <- is
				if is.Error != nil {
					return
				}
			}
		}
		for _, g := range s.Groups {
			instrs := g.ParseDrawingInstructions()
			for is := range instrs {
				s.instructions <- is
				if is.Error != nil {
					return
				}
			}
		}
	}()
	return s.instructions
}

// UnmarshalXML implements the encoding.xml.Unmarshaler interface
func (s *Svg) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	for {
		for _, attr := range start.Attr {
			if attr.Name.Local == "viewBox" {
				s.ViewBox = attr.Value
			}
			if attr.Name.Local == "width" {
				s.Width = attr.Value
			}
			if attr.Name.Local == "height" {
				s.Height = attr.Value
			}
		}

		token, err := decoder.Token()
		if err != nil {
			return err
		}

		switch tok := token.(type) {
		case xml.StartElement:
			var dip DrawingInstructionParser

			switch tok.Name.Local {
			case "g":
				g := &Group{Owner: s, Transform: mtransform.NewTransform()}
				if err = decoder.DecodeElement(g, &tok); err != nil {
					return fmt.Errorf("error decoding group element within SVG struct: %s", err)
				}
				s.Groups = append(s.Groups, *g)
				continue
			case "rect":
				dip = &Rect{}
			case "circle":
				dip = &Circle{}
			case "path":
				dip = &Path{}

			default:
				continue
			}

			if err = decoder.DecodeElement(dip, &tok); err != nil {
				return fmt.Errorf("error decoding element of SVG struct: %s", err)
			}

			s.Elements = append(s.Elements, dip)

		case xml.EndElement:
			if tok.Name.Local == "svg" {
				return nil
			}
		}
	}
}

// ParseSvg parses an SVG string into an SVG struct
func ParseSvg(str string, name string, scale float64) (*Svg, error) {
	var svg Svg
	svg.Name = name
	svg.Transform = mtransform.NewTransform()
	if scale > 0 {
		svg.Transform.Scale(scale, scale)
		svg.scale = scale
	}
	if scale < 0 {
		svg.Transform.Scale(1.0/-scale, 1.0/-scale)
		svg.scale = 1.0 / -scale
	}

	err := xml.Unmarshal([]byte(str), &svg)
	if err != nil {
		return nil, fmt.Errorf("ParseSvg Error: %v", err)
	}

	for i := range svg.Groups {
		svg.Groups[i].SetOwner(&svg)
		if svg.Groups[i].Transform == nil {
			svg.Groups[i].Transform = mtransform.NewTransform()
		}
	}
	return &svg, nil
}

// ParseSvgFromReader parses an SVG struct from an io.Reader
func ParseSvgFromReader(r io.Reader, name string, scale float64) (*Svg, error) {
	var svg Svg
	svg.Name = name
	svg.Transform = mtransform.NewTransform()
	if scale > 0 {
		svg.Transform.Scale(scale, scale)
		svg.scale = scale
	}
	if scale < 0 {
		svg.Transform.Scale(1.0/-scale, 1.0/-scale)
		svg.scale = 1.0 / -scale
	}

	if err := xml.NewDecoder(r).Decode(&svg); err != nil {
		return nil, fmt.Errorf("ParseSvg Error: %v", err)
	}

	for i := range svg.Groups {
		svg.Groups[i].SetOwner(&svg)
		if svg.Groups[i].Transform == nil {
			svg.Groups[i].Transform = mtransform.NewTransform()
		}
	}
	return &svg, nil
}

// ViewBoxValues returns all the numerical values in the viewBox
// attribute.
func (s *Svg) ViewBoxValues() ([]float64, error) {
	var vals []float64

	if s.ViewBox == "" {
		return vals, errors.New("viewBox attribute is empty")
	}

	split := strings.Split(s.ViewBox, " ")

	for _, val := range split {
		ival, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return vals, err
		}
		vals = append(vals, ival)
	}

	return vals, nil
}

// SetOwner sets the owner of a SVG Group
func (g *Group) SetOwner(svg *Svg) {
	g.Owner = svg
	for _, gn := range g.Elements {
		switch gn.(type) {
		case *Group:
			gn.(*Group).Owner = g.Owner
			gn.(*Group).SetOwner(svg)
		case *Path:
			gn.(*Path).group = g
		}
	}
}
