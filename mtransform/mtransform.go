package mtransform

import (
	"zappem.net/pub/math/geom"
)

// Transform is a co-opted 3D matrix for the purpose of performing 2D
// transformations.
type Transform geom.Matrix

// multiplyWith overwrites *a with (*a)x(*b)
func (a *Transform) multiplyWith(b *Transform) {
	*a = Transform(geom.Matrix(*a).XM(geom.Matrix(*b)))
}

// Translate returns a an (x,y) shift transform.
func Translate(x float64, y float64) Transform {
	return Transform(geom.M(
		1, 0, x,
		0, 1, y,
		0, 0, 1))
}

// Identity is the identity transformation.
func Identity() Transform {
	return Transform(geom.M(1))
}

// NewTransform returns a new transformation that is initialized as
// the identity.
func NewTransform() *Transform {
	t := Identity()
	return &t
}

// Apply transforms an (x,y) coordinate pair into (X,Y) using the
// transformation, t.
func (t *Transform) Apply(x float64, y float64) (X float64, Y float64) {
	v := geom.V(x, y, 1)
	V := geom.Matrix(*t).XV(v)
	X = V[0]
	Y = V[1]
	return
}

// MultiplyTransforms combines two transformations into a newly
// allocated one.
func MultiplyTransforms(a Transform, b Transform) Transform {
	return Transform(geom.Matrix(a).XM(geom.Matrix(b)))
}

// Scale scales a transformation independently in x and y without
// affecting its translation properties.
func (t *Transform) Scale(x float64, y float64) {
	a := geom.M(x, 0, 0,
		0, y, 0,
		0, 0, 1)
	*t = Transform(geom.Matrix(*t).XM(a))
}

// RotatePoint extends the transformation to rotate by an angle around
// a specific (x,y) point.
func (t *Transform) RotatePoint(angle geom.Angle, x float64, y float64) {
	shift := Translate(x, y)
	unshift := Translate(-x, -y)
	t.multiplyWith(&shift)
	*t = Transform(geom.Matrix(*t).XM(geom.RZ(angle)))
	t.multiplyWith(&unshift)
}
