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
	src  = flag.String("src", "/dev/stdin", "source SVG file")
	dest = flag.String("dest", "/dev/stdout", "destination SVG file")
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

func main() {
	flag.Parse()

	s := readSVG()

	trimmedGroups := 0
	for i, g := range s.Groups {
		if len(g.Elements) == 0 {
			trimmedGroups++
			continue
		}
		log.Printf("group[%d]: %#v", i, g)
		for j, e := range g.Elements {
			log.Printf("  element[%d]: %#v", j, e)
		}
	}

	log.Printf("program not written yet, but skipped %d empty groups", trimmedGroups)
	log.Printf("parsed: %#v", s)

	if *dest == "" {
		log.Fatal("please provide a --dest=output.svg argument")
	}
}
