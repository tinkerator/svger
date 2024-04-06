// Program svgoutline transforms an SVG collection of lines into a
// flattened outline SVG. This program is intended primarily to
// facilitate turning Kicad generated SVG metal masks into a laser
// cut-able outline on a PCB.
package main

import (
	"flag"
	"log"
)

var (
	src  = flag.String("src", "/dev/stdin", "source SVG file")
	dest = flag.String("dest", "/dev/stdout", "destination SVG file")
)

func main() {
	flag.Parse()

	if *src == "" {
		log.Fatal("please provide a --src=file.svg argument")
	}

	log.Print("program not written yet")

	if *dest == "" {
		log.Fatal("please provide a --dest=output.svg argument")
	}
}
