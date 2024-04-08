package svger

import (
	"fmt"
	"strconv"

	gl "zappem.net/pub/graphics/svger/genericlexer"
	"zappem.net/pub/graphics/svger/mtransform"
	"zappem.net/pub/math/geom"
)

func parseNumber(i gl.Item) (float64, error) {
	var n float64
	var ok error
	if i.Type == gl.ItemNumber {
		n, ok = strconv.ParseFloat(i.Value, 64)
		if ok != nil {
			return n, fmt.Errorf("Error passing number %s", ok)
		}
	}
	return n, nil
}

func parseTuple(l *gl.Lexer) (Tuple, error) {
	t := Tuple{}

	l.ConsumeWhiteSpace()

	ni := l.NextItem()
	if ni.Type == gl.ItemNumber {
		n, ok := strconv.ParseFloat(ni.Value, 64)
		if ok != nil {
			return t, fmt.Errorf("Error parsing number %s", ok)
		}
		t[0] = n
	} else {
		return t, fmt.Errorf("Error parsing Tuple expected Number got: %s", ni.Value)
	}

	if l.PeekItem().Type == gl.ItemWSP || l.PeekItem().Type == gl.ItemComma {
		l.NextItem()
	}
	ni = l.NextItem()
	if ni.Type == gl.ItemNumber {
		n, ok := strconv.ParseFloat(ni.Value, 64)
		if ok != nil {
			return t, fmt.Errorf("Error passing Number %s", ok)
		}
		t[1] = n
	} else {
		return t, fmt.Errorf("Error passing Tuple expected Number got: %v", ni)
	}

	return t, nil
}

func parseTransform(tstring string) (mtransform.Transform, error) {
	lexer, _ := gl.Lex("tlexer", tstring)
	for {
		i := lexer.NextItem()
		switch i.Type {
		case gl.ItemEOS:
			return mtransform.Identity(),
				fmt.Errorf("transform parse failed")
		case gl.ItemWord:
			switch i.Value {
			case "matrix":
				return parseMatrix(lexer)
			case "translate":
				return parseTranslate(lexer)
			case "rotate":
				return parseRotate(lexer)
			case "scale":
				return parseScale(lexer)
			}
		}
	}
}

func parseMatrix(l *gl.Lexer) (mtransform.Transform, error) {
	nums, err := parseParenNumList(l)
	if err != nil || len(nums) != 6 {
		return mtransform.Identity(),
			fmt.Errorf("Error Parsing Transform Matrix: %v", err)
	}
	tm := mtransform.Transform(geom.M(
		nums[0], nums[2], nums[4],
		nums[1], nums[3], nums[5],
		0, 0, 1))
	return tm, nil
}

func parseTranslate(l *gl.Lexer) (mtransform.Transform, error) {
	nums, err := parseParenNumList(l)
	if err != nil || len(nums) != 2 {
		return mtransform.Identity(), fmt.Errorf("Error Parsing Translate: %v", err)
	}
	tm := mtransform.Translate(nums[0], nums[1])
	return tm, nil
}

func parseRotate(l *gl.Lexer) (mtransform.Transform, error) {
	nums, err := parseParenNumList(l)
	if err != nil || (len(nums) != 1 && len(nums) != 3) {
		return mtransform.Identity(), fmt.Errorf("Error Parsing Rotate: %v", err)
	}
	a, px, py := nums[0], 0.0, 0.0
	if len(nums) == 3 {
		px, py = nums[1], nums[2]
	}

	tm := mtransform.Identity()
	(&tm).RotatePoint(geom.Radians(a), px, py)
	return tm, nil
}

func parseScale(l *gl.Lexer) (mtransform.Transform, error) {
	nums, err := parseParenNumList(l)
	if err != nil || (len(nums) != 1 && len(nums) != 2) {
		return mtransform.Identity(), fmt.Errorf("Error Parsing Scale: %v", err)
	}
	x, y := nums[0], nums[0]
	if len(nums) == 2 {
		y = nums[1]
	}
	tm := mtransform.Identity()
	(&tm).Scale(x, y)
	return tm, nil
}

// Parse a parenthesized list of ncount numbers.
func parseParenNumList(l *gl.Lexer) ([]float64, error) {
	i := l.NextItem()
	if i.Type != gl.ItemParan {
		return nil, fmt.Errorf("Expected Opening Parantheses")
	}
	var nums []float64
	for {
		if len(nums) > 0 {
			for l.PeekItem().Type == gl.ItemComma || l.PeekItem().Type == gl.ItemWSP {
				l.NextItem()
			}
		}
		if l.PeekItem().Type == gl.ItemParan {
			l.NextItem() // consume Parantheses
			break
		}
		if l.PeekItem().Type != gl.ItemNumber {
			return nil, fmt.Errorf("Expected Number got %v", l.PeekItem().String())
		}
		n, err := parseNumber(l.NextItem())
		if err != nil {
			return nil, err
		}
		nums = append(nums, n)

	}
	return nums, nil
}
