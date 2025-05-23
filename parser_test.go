package svger

import (
	"strings"
	"testing"
)

func TestParseSvg(t *testing.T) {
	vs := []string{
		`<?xml version="1.0" encoding="utf-8"?>
<!-- Generator: Adobe Illustrator 15.0.2, SVG Export Plug-In . SVG Version: 6.00 Build 0)  -->
<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN" "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd">
<svg version="1.1" id="Layer_1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" x="0px" y="0px"
	 width="595.201px" height="841.922px" viewBox="0 0 595.201 841.922" enable-background="new 0 0 595.201 841.922"
	 xml:space="preserve">
<rect x="207" y="53" fill="#009FE3" width="181.667" height="85.333"/>
<text transform="matrix(1 0 0 1 232.3306 107.5952)" fill="#FFFFFF" font-family="'ArialMT'" font-size="31.9752">PODIUM</text>
<g><text transform="matrix(1 0 0 1 232.3306 107.5952)" fill="#FFFFFF" font-family="'ArialMT'" font-size="31.9752">PODIUM</text></g>
</svg>`,
		`<?xml version="1.0" standalone="no"?>
 <!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN"
 "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd">
<svg
  xmlns:svg="http://www.w3.org/2000/svg"
  xmlns="http://www.w3.org/2000/svg"
  xmlns:xlink="http://www.w3.org/1999/xlink"
  version="1.1"
  width="297.0022mm" height="210.0072mm" viewBox="0.0000 0.0000 297.0022 210.0072">
<title>SVG Image created as test-board-F_Silkscreen.svg date 2025/04/22 20:13:19 </title>
  <desc>Image generated by PCBNEW </desc>
<g style="fill:none;
stroke:#000000; stroke-width:0.1500; stroke-opacity:1;
stroke-linecap:round; stroke-linejoin:round;">
<g class="stroked-text"><desc>J4</desc>
<path d="M107.7667 67.3098
L107.7667 68.0241
" />
<path d="M107.7667 68.0241
L107.7190 68.1670
" />
<path d="M107.7190 68.1670
L107.6238 68.2622
" />
<path d="M107.6238 68.2622
L107.4810 68.3098
" />
<path d="M107.4810 68.3098
L107.3857 68.3098
" />
<path d="M108.6714 67.6432
L108.6714 68.3098
" />
<path d="M108.4333 67.2622
L108.1952 67.9765
" />
<path d="M108.1952 67.9765
L108.8143 67.9765
" />
</g></g>
</svg>
`,
	}

	for i, sText := range vs {
		if svg, err := ParseSvg(sText, "test", 0); err != nil {
			t.Errorf("[%d] parsing failed: %v", i, err)
		} else if svg == nil {
			t.Errorf("[%d] ParseSvg returned nil without error", i)
		}
		svg, err := ParseSvgFromReader(strings.NewReader(sText), "test", 0)
		if err != nil {
			t.Errorf("[%d] ParseSvgFromReader parsing failed: %v", i, err)
		} else if svg == nil {
			t.Errorf("[%d] ParseSvgFromReader returned nil without error", i)
		}
	}
}
