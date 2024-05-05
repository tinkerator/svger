package svger

import (
	"fmt"
	"log"
	"strconv"

	gl "zappem.net/pub/graphics/svger/genericlexer"
	mt "zappem.net/pub/graphics/svger/mtransform"
)

// Path is an SVG XML path element
type Path struct {
	ID              string `xml:"id,attr"`
	D               string `xml:"d,attr"`
	Style           string `xml:"style,attr"`
	TransformString string `xml:"transform,attr"`
	properties      map[string]string
	StrokeWidth     float64 `xml:"stroke-width,attr"`
	Fill            *string `xml:"fill,attr"`
	Stroke          *string `xml:"stroke,attr"`
	StrokeLineCap   *string `xml:"stroke-linecap,attr"`
	StrokeLineJoin  *string `xml:"stroke-linejoin,attr"`
	Segments        chan Segment
	instructions    chan *DrawingInstruction
	errors          chan error
	group           *Group
}

// A Segment of a path that contains a list of connected points, its
// stroke Width and if the segment forms a closed loop.  Points are
// defined in world space after any matrix transformation is applied.
type Segment struct {
	Width  float64
	Closed bool
	Points [][2]float64
}

func (p Path) newSegment(start [2]float64) *Segment {
	var s Segment
	s.Width = p.StrokeWidth * p.group.Owner.scale
	s.Points = append(s.Points, start)
	return &s
}

func (s *Segment) addPoint(p [2]float64) {
	s.Points = append(s.Points, p)
}

type pathDescriptionParser struct {
	p              *Path
	lex            *gl.Lexer
	x, y           float64
	currentcommand int
	tokbuf         [4]gl.Item
	peekcount      int
	lasttuple      Tuple
	transform      mt.Transform
	svg            *Svg
	currentsegment *Segment
}

func newPathDParse() *pathDescriptionParser {
	pdp := &pathDescriptionParser{}
	pdp.transform = mt.Identity()
	return pdp
}

// ParseDrawingInstructions returns two channels of:
// DrawingInstructions and errors. The former can be used to pass to a
// path drawing library.
func (p *Path) ParseDrawingInstructions() (chan *DrawingInstruction, chan error) {
	p.parseStyle()
	pdp := newPathDParse()
	pdp.p = p
	if p.group == nil {
		p.group = new(Group)
		temp := mt.Identity()
		p.group.Transform = &temp
	}
	pdp.svg = p.group.Owner
	pathTransform := mt.Identity()
	if p.TransformString != "" {
		if pt, err := parseTransform(p.TransformString); err == nil {
			pathTransform = pt
		}
	}
	pdp.transform = mt.MultiplyTransforms(pdp.transform, *p.group.Transform)
	pdp.transform = mt.MultiplyTransforms(pdp.transform, pathTransform)

	p.instructions = make(chan *DrawingInstruction, 100)
	p.errors = make(chan error, 100)
	l, _ := gl.Lex(fmt.Sprint(p.ID), p.D)

	pdp.lex = l
	go func() {
		defer close(p.instructions)
		defer close(p.errors)
		var count int
		for {
			i := pdp.lex.NextItem()
			count++
			switch {
			case i.Type == gl.ItemError:
				return
			case i.Type == gl.ItemEOS:
				scaledStrokeWidth := p.StrokeWidth * pdp.p.group.Owner.scale
				pdp.p.instructions <- &DrawingInstruction{
					Kind:           PaintInstruction,
					StrokeWidth:    &scaledStrokeWidth,
					Stroke:         p.Stroke,
					StrokeLineCap:  p.StrokeLineCap,
					StrokeLineJoin: p.StrokeLineJoin,
					Fill:           p.Fill,
				}
				return
			case i.Type == gl.ItemLetter:
				err := pdp.parseCommandDrawingInstructions(l, i)
				if err != nil {
					p.errors <- fmt.Errorf("error when parsing instruction number %d: %s", count, err)
					return
				}
			default:
				fmt.Printf("Default invoked: %d item %v\n", count, i)
			}
		}
	}()

	return p.instructions, p.errors
}

// parseCommandDrawingInstructions keys off a command letter and
// performs specific further parsing of a path.
func (pdp *pathDescriptionParser) parseCommandDrawingInstructions(l *gl.Lexer, i gl.Item) error {
	switch i.Value {
	case "M":
		return pdp.parseMoveToAbsDI()
	case "m":
		return pdp.parseMoveToRelDI()
	case "c":
		return pdp.parseCurveToRelDI()
	case "C":
		return pdp.parseCurveToAbsDI()
	case "l":
		return pdp.parseLineToRelDI()
	case "L":
		return pdp.parseLineToAbsDI()
	case "H", "h":
		return pdp.parseHLineToDI(i.Value == "H")
	case "V", "v":
		return pdp.parseVLineToDI(i.Value == "V")
	case "z", "Z":
		return pdp.parseCloseDI()
	}

	return fmt.Errorf("unknown command found in SVG: %s", i.Value)
}

func (pdp *pathDescriptionParser) parseMoveToAbsDI() error {
	var tuples []Tuple

	t, err := parseTuple(pdp.lex)
	if err != nil {
		return fmt.Errorf("error parsing MoveToAbs. Expected tuple: %s", err)
	}

	// Update current cursor location with initial point.
	pdp.x = t[0]
	pdp.y = t[1]

	if pdp.p.group.Owner == nil {
		pdp.p.group.Owner = &Svg{scale: 1}
	}
	if pdp.p.StrokeWidth == -1 {
		pdp.p.StrokeWidth = 1
	}

	pdp.lex.ConsumeWhiteSpace()
	for pdp.lex.PeekItem().Type == gl.ItemNumber {
		t, err := parseTuple(pdp.lex)
		if err != nil {
			return fmt.Errorf("Error Passing MoveToAbs\n%s", err)
		}
		tuples = append(tuples, t)
		pdp.lex.ConsumeWhiteSpace()
	}

	x, y := pdp.transform.Apply(pdp.x, pdp.y)
	pdp.p.instructions <- &DrawingInstruction{Kind: MoveInstruction, M: &Tuple{x, y}}

	for _, nt := range tuples {
		pdp.x = nt[0]
		pdp.y = nt[1]
		x, y = pdp.transform.Apply(pdp.x, pdp.y)
		pdp.p.instructions <- &DrawingInstruction{Kind: LineInstruction, M: &Tuple{x, y}}
	}

	return nil
}

func (pdp *pathDescriptionParser) parseLineToAbsDI() error {
	var tuples []Tuple
	pdp.lex.ConsumeWhiteSpace()
	for pdp.lex.PeekItem().Type == gl.ItemNumber {
		t, err := parseTuple(pdp.lex)
		if err != nil {
			return fmt.Errorf("Error Passing LineToAbs\n%s", err)
		}
		tuples = append(tuples, t)
		pdp.lex.ConsumeWhiteSpace()
	}
	if len(tuples) > 0 {
		x, y := pdp.transform.Apply(pdp.x, pdp.y)

		for _, nt := range tuples {
			pdp.x = nt[0]
			pdp.y = nt[1]
			x, y = pdp.transform.Apply(pdp.x, pdp.y)
			pdp.p.instructions <- &DrawingInstruction{Kind: LineInstruction, M: &Tuple{x, y}}
		}
	}

	return nil
}

func (pdp *pathDescriptionParser) parseMoveToRelDI() error {
	pdp.lex.ConsumeWhiteSpace()
	t, err := parseTuple(pdp.lex)
	if err != nil {
		return fmt.Errorf("Error Passing MoveToRel Expected First Tuple %s", err)
	}

	pdp.x += t[0]
	pdp.y += t[1]

	var tuples []Tuple
	pdp.lex.ConsumeWhiteSpace()
	for pdp.lex.PeekItem().Type == gl.ItemNumber {
		t, err := parseTuple(pdp.lex)
		if err != nil {
			return fmt.Errorf("Error Passing MoveToRel\n%s", err)
		}
		tuples = append(tuples, t)
		pdp.lex.ConsumeWhiteSpace()
	}

	x, y := pdp.transform.Apply(pdp.x, pdp.y)
	pdp.p.instructions <- &DrawingInstruction{Kind: MoveInstruction, M: &Tuple{x, y}}

	for _, nt := range tuples {
		pdp.x += nt[0]
		pdp.y += nt[1]
		x, y = pdp.transform.Apply(pdp.x, pdp.y)
		pdp.p.instructions <- &DrawingInstruction{Kind: LineInstruction, M: &Tuple{x, y}}
	}

	return nil
}

func (pdp *pathDescriptionParser) parseHLineToDI(abs bool) error {
	coords := []float64{}
	pdp.lex.ConsumeWhiteSpace()
	for pdp.lex.PeekItem().Type == gl.ItemNumber {
		item := pdp.lex.NextItem()
		c, err := strconv.ParseFloat(item.Value, 64)
		if err != nil {
			return fmt.Errorf("parsing %q: %s", item.Value, err)
		}
		coords = append(coords, c)
		pdp.lex.ConsumeWhiteSpace()
	}
	if len(coords) > 0 {
		for _, c := range coords {
			if abs {
				pdp.x = c
			} else {
				pdp.x += c
			}
			x, y := pdp.transform.Apply(pdp.x, pdp.y)
			pdp.p.instructions <- &DrawingInstruction{Kind: LineInstruction, M: &Tuple{x, y}}
		}
	}
	return nil
}

func (pdp *pathDescriptionParser) parseLineToRelDI() error {
	var tuples []Tuple
	pdp.lex.ConsumeWhiteSpace()
	for pdp.lex.PeekItem().Type == gl.ItemNumber {
		t, err := parseTuple(pdp.lex)
		if err != nil {
			return fmt.Errorf("Error Passing LineToRel\n%s", err)
		}
		tuples = append(tuples, t)
		pdp.lex.ConsumeWhiteSpace()
	}

	if len(tuples) > 0 {
		x, y := pdp.transform.Apply(pdp.x, pdp.y)

		for _, nt := range tuples {
			pdp.x += nt[0]
			pdp.y += nt[1]
			x, y = pdp.transform.Apply(pdp.x, pdp.y)
			pdp.p.instructions <- &DrawingInstruction{Kind: LineInstruction, M: &Tuple{x, y}}
		}
	}

	return nil
}

func (pdp *pathDescriptionParser) parseVLineToDI(abs bool) error {
	pdp.lex.ConsumeWhiteSpace()
	coords := []float64{}
	for pdp.lex.PeekItem().Type == gl.ItemNumber {
		n, err := parseNumber(pdp.lex.NextItem())
		if err != nil {
			return fmt.Errorf("Error Passing VLineToRel\n%s", err)
		}
		coords = append(coords, n)
		pdp.lex.ConsumeWhiteSpace()
	}

	for _, n := range coords {
		if abs {
			pdp.y = n
		} else {
			pdp.y += n
		}
		x, y := pdp.transform.Apply(pdp.x, pdp.y)
		pdp.p.instructions <- &DrawingInstruction{Kind: LineInstruction, M: &Tuple{x, y}}
	}

	return nil
}

func (pdp *pathDescriptionParser) parseCloseDI() error {
	pdp.lex.ConsumeWhiteSpace()
	pdp.p.instructions <- &DrawingInstruction{Kind: CloseInstruction}
	return nil
}

func (pdp *pathDescriptionParser) parseCurveToRelDI() error {
	var tuples []Tuple
	pdp.lex.ConsumeWhiteSpace()
	for pdp.lex.PeekItem().Type == gl.ItemNumber {
		t, err := parseTuple(pdp.lex)
		if err != nil {
			return fmt.Errorf("Error Passing CurveToRel\n%s", err)
		}
		tuples = append(tuples, t)
		pdp.lex.ConsumeWhiteSpace()
	}
	x, y := pdp.transform.Apply(pdp.x, pdp.y)

	for j := 0; j < len(tuples)/3; j++ {
		c1x, c1y := pdp.transform.Apply(x+tuples[j*3][0], y+tuples[j*3][1])
		c2x, c2y := pdp.transform.Apply(x+tuples[j*3+1][0], y+tuples[j*3+1][1])
		tx, ty := pdp.transform.Apply(x+tuples[j*3+2][0], y+tuples[j*3+2][1])

		pdp.p.instructions <- &DrawingInstruction{
			Kind: CurveInstruction,
			CurvePoints: &CurvePoints{C1: &Tuple{c1x, c1y},
				C2: &Tuple{c2x, c2y},
				T:  &Tuple{tx, ty},
			},
		}

		pdp.x += tuples[j*3+2][0]
		pdp.y += tuples[j*3+2][1]
		x, y = pdp.transform.Apply(pdp.x, pdp.y)
	}

	return nil
}

func (pdp *pathDescriptionParser) parseCurveToAbsDI() error {
	var (
		tuples      []Tuple
		instrTuples []Tuple
	)

	pdp.lex.ConsumeWhiteSpace()
	for pdp.lex.PeekItem().Type == gl.ItemNumber {
		t, err := parseTuple(pdp.lex)
		if err != nil {
			return fmt.Errorf("Error parsing CurveToRel\n%s", err)
		}
		tuples = append(tuples, t)
		pdp.lex.ConsumeWhiteSpace()
		pdp.lex.ConsumeComma()
	}

	for j := 0; j < len(tuples)/3; j++ {
		for _, nt := range tuples[j*3 : (j+1)*3] {
			pdp.x = nt[0]
			pdp.y = nt[1]

			tx, ty := pdp.transform.Apply(pdp.x, pdp.y)
			instrTuples = append(instrTuples, Tuple{tx, ty})
		}

		pdp.p.instructions <- &DrawingInstruction{
			Kind: CurveInstruction,
			CurvePoints: &CurvePoints{
				C1: &instrTuples[0],
				C2: &instrTuples[1],
				T:  &instrTuples[2],
			},
		}
	}

	return nil
}

// pathStyle parses the path style string and converts it into
// path properties.
func (p *Path) parseStyle() {
	p.properties = splitStyle(p.Style)
	for key, val := range p.properties {
		switch key {
		case "stroke-width":
			p.StrokeWidth = parseWidth(val)
		case "fill":
			p.Fill = refString(val)
		case "stroke":
			p.Stroke = refString(val)
		default:
			if Debug {
				log.Printf("TODO parse %q property %q", key, val)
			}
		}
	}
}
