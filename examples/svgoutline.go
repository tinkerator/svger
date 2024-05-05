// Program svgoutline transforms an SVG collection of lines into a
// flattened outline SVG. This program is intended primarily to
// facilitate turning Kicad generated SVG metal masks into a laser
// cut-able outline on a PCB.
package main

import (
	"flag"
	"log"
	"os"

	"zappem.net/pub/graphics/svger"
)

var (
	src   = flag.String("src", "/dev/stdin", "source SVG file")
	dest  = flag.String("dest", "/dev/stdout", "destination SVG file")
	debug = flag.Bool("debug", false, "extra debugging output")
)

// read an SVG or fail the program.
func readSVG() *svger.Svg {
	if *src == "" {
		log.Fatal("please provide a --src=file.svg argument")
	}

	f, err := os.Open(*src)
	if err != nil {
		log.Fatalf("failed to open %q: %v", *src, err)
	}
	defer f.Close()

	s, err := svger.ParseSvgFromReader(f, *src, 1)
	if err != nil {
		log.Fatalf("failed to parse %q: %v", *src, err)
	}
	return s
}

// examinePath investigates the elements of a path.
func examinePath(p *svger.Path) {
	ins, errs := p.ParseDrawingInstructions()
	for {
		err := <-errs
		if err != nil {
			log.Fatalf("examinePath encountered an error: %v", err)
		}
		i, ok := <-ins
		if !ok {
			// End of path
			return
		}
		log.Printf("  %v:", i.Kind)
		if i.M != nil {
			log.Printf("    PathPoint=%v", *i.M)
		}
		if i.CurvePoints != nil {
			log.Printf("    CurvePoints=%v", *i.CurvePoints)
		}
		if i.Radius != nil {
			log.Printf("    Radius=%v", *i.Radius)
		}
		if i.StrokeWidth != nil {
			log.Printf("    StrokeWidth=%v", *i.StrokeWidth)
		}
		if i.Fill != nil {
			log.Printf("    Fill=%v", *i.Fill)
		}
		if i.Stroke != nil {
			log.Printf("    Stroke=%v", *i.Stroke)
		}
		if i.StrokeLineCap != nil {
			log.Printf("    StrokeLineCap=%v", *i.StrokeLineCap)
		}
		if i.StrokeLineJoin != nil {
			log.Printf("    StrokeLineJoin=%v", *i.StrokeLineJoin)
		}
	}
}

func examineCircle(c *svger.Circle) {
	ins, errs := c.ParseDrawingInstructions()
	for {
		err := <-errs
		if err != nil {
			log.Fatalf("examineCircle encountered an error: %v", err)
		}
		i, ok := <-ins
		if !ok {
			// End of path
			return
		}
		log.Printf("  %v:", i.Kind)
		if i.M != nil {
			log.Printf("    PathPoint=%v", *i.M)
		}
		if i.CurvePoints != nil {
			log.Printf("    CurvePoints=%v", *i.CurvePoints)
		}
		if i.Radius != nil {
			log.Printf("    Radius=%v", *i.Radius)
		}
		if i.StrokeWidth != nil {
			log.Printf("    StrokeWidth=%v", *i.StrokeWidth)
		}
		if i.Fill != nil {
			log.Printf("    Fill=%v", *i.Fill)
		}
		if i.Stroke != nil {
			log.Printf("    Stroke=%v", *i.Stroke)
		}
		if i.StrokeLineCap != nil {
			log.Printf("    StrokeLineCap=%v", *i.StrokeLineCap)
		}
		if i.StrokeLineJoin != nil {
			log.Printf("    StrokeLineJoin=%v", *i.StrokeLineJoin)
		}
	}
}

// examineGroup investigates the members of a group.
func examineGroup(g *svger.Group) {
	for j, e := range g.Elements {
		switch e.(type) {
		case *svger.Path:
			if *debug {
				log.Printf("  path[%d]: %#v", j, e.(*svger.Path))
			}
			examinePath(e.(*svger.Path))
		case *svger.Circle:
			examineCircle(e.(*svger.Circle))
		case *svger.Group:
			// nested groups.
			examineGroup(e.(*svger.Group))
		default:
			log.Printf("  element[%d]: %#v", j, e)
		}
	}
}

// parseSVG works through the parsed SVG data structures, group by
// group.
func parseSVG(s *svger.Svg) {
	trimmedGroups := 0
	for i, g := range s.Groups {
		if len(g.Elements) == 0 {
			trimmedGroups++
			continue
		}
		if *debug {
			log.Printf("group[%d]:", i)
		}
		examineGroup(&g)
	}
	log.Printf("program not written yet, but skipped %d empty groups", trimmedGroups)
}

func main() {
	flag.Parse()
	svger.Debug = *debug

	s := readSVG()

	if *debug {
		log.Printf("SVG: %#v", s)
	}

	parseSVG(s)

	if *dest == "" {
		log.Fatal("please provide a --dest=output.svg argument")
	}
}
