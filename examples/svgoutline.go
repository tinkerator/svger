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
	src     = flag.String("src", "/dev/stdin", "source SVG file")
	dest    = flag.String("dest", "/dev/stdout", "destination SVG file")
	debug   = flag.Bool("debug", false, "extra debugging output")
	scriber = flag.Float64("scriber", 0.1, "mm width of the thinnest line")
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

// decodeSVG unravels the content of the SVG into sequences of svger.DrawingInstructions.
func decodeSVG(s *svger.Svg) (dis []*svger.DrawingInstruction, err error) {
	ins := s.ParseDrawingInstructions()
	for {
		i, ok := <-ins
		if !ok {
			// End of svg
			return
		}
		dis = append(dis, i)
		if i.Error != nil {
			err = i.Error
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

func main() {
	flag.Parse()
	svger.Debug = *debug

	s := readSVG()

	if *debug {
		log.Printf("SVG: %#v", s)
	}

	dis, err := decodeSVG(s)
	if err != nil {
		log.Fatalf("failed to fully decode SVG: %v", err)
	}

	if *debug {
		for i, di := range dis {
			log.Printf("%4d: %#v", i, *di)
		}
	}

	if *dest == "" {
		log.Fatal("please provide a --dest=output.svg argument")
	}
}
