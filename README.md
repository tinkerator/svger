# svger

Go package for parsing and modifying svg files. It includes support
for Bezier Curve rasterization.

The package's main purpose is to transform multi-group SVG kicad
generated SVG metal layers into a flattened outline represenation,
suitable for laser cutting on a Snapmaker 2.

This operation will eventually be performed by the included
`example/svgoutline.go` program. A work in progress, at present.

## History

This package was evolved from a forked version of:
[`github.com/rustyoz/svg`](https://github.com/rustyoz/svg), trimmed of
superflous dependencies and merged with the two independent packages:
[`github.com/rustyoz/Mtransform`](https://github.com/rustyoz/Mtransform)
and
[`github.com/rustyoz/genericlexer`](https://github.com/rustyoz/genericlexer).
