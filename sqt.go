package main

import (
	"math"
)

// SQT represents a transformation consisting of a scaling factor,
// a rotation and a translation.
type SQT struct {
	scale          float64 // scaling factor
	qx, qy, qz, qw float64 // rotation
	tx, ty, tz     float64 // translation
}

// NewSQT creates a new identity transformation.
func NewSQT() *SQT {
	return &SQT{
		1,
		0, 0, 0, 1,
		0, 0, 0,
	}
}

// SetRotation sets the rotation of the SQT transformation to the
// angle theta about the axis (x,y,z) (which must be normalised).
func (s *SQT) SetRotation(theta, x, y, z float64) {
	sin := math.Sin(theta / 2)
	s.qx = x * sin
	s.qy = y * sin
	s.qz = z * sin
	s.qw = math.Cos(theta / 2)
}

// SetTranslation sets the translation of the SQT transformation.
func (s *SQT) SetTranslation(x, y, z float64) {
	s.tx = x
	s.ty = y
	s.tz = z
}

// SetScale sets the scale of the SQT transformation in each axis.
func (s *SQT) SetScale(scale float64) {
	s.scale = scale
}

// qmult calculates the Grassman product of two quaternions.
func qmult(ax, ay, az, aw, bx, by, bz, bw float64) (cx, cy, cz, cw float64) {
	cx = aw*bx + ax*bw + ay*bz - az*by
	cy = aw*by + ay*bw + az*bx - ax*bz
	cz = aw*bz + az*bw + ax*by - ay*bx
	cw = aw*bw - (ax*bx + ay*by + az*bz)
	return
}

// Rotate adds another rotation by the angle theta about the axis (x,y,z)
// (which must be normalised) to the SQT transformation.
func (s *SQT) Rotate(theta, x, y, z float64) {
	sin := math.Sin(theta / 2)
	qx := x * sin
	qy := y * sin
	qz := z * sin
	qw := math.Cos(theta / 2)
	s.qx, s.qy, s.qz, s.qw = qmult(qx, qy, qz, qw, s.qx, s.qy, s.qz, s.qw)
}

// Translate adds another translation to the SQT transformation.
func (s *SQT) Translate(x, y, z float64) {
	s.tx += x
	s.ty += y
	s.tz += z
}

// Scale adds another scaling factor to each axis in the SQT transformation.
func (s *SQT) Scale(scale float64) {
	s.scale *= scale
}

// TransformRel applys the SQT transformation without the tranlation to the input vector (ix, iy, iz),
// producing the output vector (ox, oy, oz) (equivalent to premultiplying (ix, iy, iz, 0) by s.Matrix()).
func (s *SQT) TransformRel(ix, iy, iz float64) (ox, oy, oz float64) {
	sx := ix * s.scale
	sy := iy * s.scale
	sz := iz * s.scale

	x, y, z, w := qmult(sx, sy, sz, 0, -s.qx, -s.qy, -s.qz, s.qw)
	ox, oy, oz, _ = qmult(s.qx, s.qy, s.qz, s.qw, x, y, z, w)
	return
}

// TransformRel applys the SQT transformation to the input vector (ix, iy, iz),
// producing the output vector (ox, oy, oz) (equivalent to premultiplying (ix, iy, iz, 1) by s.Matrix()).
func (s *SQT) TransformAbs(ix, iy, iz float64) (ox, oy, oz float64) {
	rx, ry, rz := s.TransformRel(ix, iy, iz)

	ox = rx + s.tx
	oy = ry + s.ty
	oz = rz + s.tz
	return
}

// Commpose returns the SQT transformation representing the application of first transform, then s.
func (s *SQT) Compose(transform *SQT) (o *SQT) {
	o = NewSQT()
	o.scale = transform.scale * s.scale
	o.qx, o.qy, o.qz, o.qw = qmult(s.qx, s.qy, s.qz, s.qw, transform.qx, transform.qy, transform.qz, transform.qw)
	o.tx, o.ty, o.tz = s.TransformAbs(transform.tx, transform.ty, transform.tz)
	return
}

// Inverse returns the SQT transformation that undoes this SQT transformation, such that s.Compose(s.Inverse)
// is the identity transformation.
func (s *SQT) Inverse() *SQT {
	sx := s.tx / s.scale
	sy := s.ty / s.scale
	sz := s.tz / s.scale

	x, y, z, w := qmult(sx, sy, sz, 0, s.qx, s.qy, s.qz, s.qw)
	tx, ty, tz, _ := qmult(-s.qx, -s.qy, -s.qz, s.qw, x, y, z, w)

	return &SQT{
		1 / s.scale,
		-s.qx, -s.qy, -s.qz, s.qw,
		-tx, -ty, -tz,
	}
}

// Matrix returns a representation of the SQT transformation as an affine
// transformation matrix suitable for OpenGL rendering.
func (s *SQT) Matrix() [16]float32 {
	return [16]float32{
		float32((1 - 2*s.qy*s.qy - 2*s.qz*s.qz) * s.scale), float32((2*s.qx*s.qy - 2*s.qz*s.qw) * s.scale), float32((2*s.qx*s.qz + 2*s.qy*s.qw) * s.scale), float32(s.tx),
		float32((2*s.qx*s.qy + 2*s.qz*s.qw) * s.scale), float32((1 - 2*s.qx*s.qx - 2*s.qz*s.qz) * s.scale), float32((2*s.qy*s.qz - 2*s.qx*s.qw) * s.scale), float32(s.ty),
		float32((2*s.qx*s.qz - 2*s.qy*s.qw) * s.scale), float32((2*s.qy*s.qz + 2*s.qx*s.qw) * s.scale), float32((1 - 2*s.qx*s.qx - 2*s.qy*s.qy) * s.scale), float32(s.tz),
		0, 0, 0, 1,
	}
}
