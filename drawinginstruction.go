package svger

import "fmt"

// InstructionType tells our path drawing library which function it has
// to call
type InstructionType int

// These are instruction types that we use with our path drawing library
const (
	ErrorInstruction InstructionType = iota
	MoveInstruction
	LineInstruction
	CloseInstruction
	PaintInstruction
	CircleInstruction
	CurveInstruction
)

// CurvePoints are the points needed by a bezier curve.
type CurvePoints struct {
	C1 *Tuple
	C2 *Tuple
	T  *Tuple
}

// DrawingInstruction contains enough information that a simple drawing
// library can draw the shapes contained in an SVG file.
//
// The struct contains all necessary fields but only the ones needed (as
// indicated byt the InstructionType) will be non-nil.
type DrawingInstruction struct {
	Kind           InstructionType
	Error          error
	M              *Tuple
	CurvePoints    *CurvePoints
	Radius         *float64
	StrokeWidth    *float64
	Fill           *string
	Stroke         *string
	StrokeLineCap  *string
	StrokeLineJoin *string
}

// DrawingInstructionParser allow getting segments and drawing
// instructions from them. All SVG elements should implement this
// interface.
type DrawingInstructionParser interface {
	ParseDrawingInstructions() chan *DrawingInstruction
}

// String describes the kind of instruction type.
func (kind InstructionType) String() string {
	switch kind {
	case ErrorInstruction:
		return "Error"
	case MoveInstruction:
		return "Move"
	case CircleInstruction:
		return "Circle"
	case CurveInstruction:
		return "Curve"
	case LineInstruction:
		return "Line"
	case CloseInstruction:
		return "Close"
	case PaintInstruction:
		return "Paint"
	default:
		return fmt.Sprintf("unknown InstructionType[%d]", kind)
	}
}
