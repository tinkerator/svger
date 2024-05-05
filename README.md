# svger

Go package for parsing re-rendering svg files. It includes support for
Bezier Curve rasterization. See history for where the initial version
of the code originated.

## Overview

The `svger` package parses SVG files and generates a series of drawing
instructions to re-render them.

Example, using the not-yet-finished `svsoutline` program:

```
$ go run examples/svgoutline.go --src examples/test-board-F_Cu.svg
```

Automated documentation for the svger package can be found on
[go.dev](https://pkg.go.dev/zappem.net/pub/graphics/svger).

## Planned changes

The package's main purpose is to transform multi-group SVG kicad
generated SVG metal layers into a flattened outline represenation,
suitable for laser cutting on a Snapmaker 2.

## History

This package was evolved from a forked version of:
[`github.com/rustyoz/svg`](https://github.com/rustyoz/svg), trimmed of
superflous dependencies and merged with the two independent packages:
[`github.com/rustyoz/Mtransform`](https://github.com/rustyoz/Mtransform)
and
[`github.com/rustyoz/genericlexer`](https://github.com/rustyoz/genericlexer).
