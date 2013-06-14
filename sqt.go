package main

import (
	"math"
)

// SQT represents a transformation consisting of a scaling factor in each axis,
// a rotation and a translation.
type SQT struct {
	sx, sy, sz     float64
	qx, qy, qz, qw float64
	tx, ty, tz     float64
}

// NewSQT creates a new identity transformation.
func NewSQT() *SQT {
	return &SQT{
		1, 1, 1,
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
func (s *SQT) SetScale(x, y, z float64) {
	s.sx = x
	s.sy = y
	s.sz = z
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
	//s.tx, s.ty, s.tz, _ = qmult(qx, qy, qz, qw, qmult(s.tx, s.ty, s.tz, 0, -qx, -qy, -qz, qw))
}

// Translate adds another translation to the SQT transformation.
func (s *SQT) Translate(x, y, z float64) {
	s.tx += x
	s.ty += y
	s.tz += z
}

// Scale adds another scaling factor to each axis in the SQT transformation.
func (s *SQT) Scale(x, y, z float64) {
	s.sx *= x
	s.sy *= y
	s.sz *= z
}

// Matrix returns a representation of the SQT transformation as an affine
// transformation matrix suitable for OpenGL rendering.
func (s *SQT) Matrix() [16]float32 {
	return [16]float32{
		float32((1 - 2*s.qy*s.qy - 2*s.qz*s.qz) * s.sx), float32(2*s.qx*s.qy - 2*s.qz*s.qw), float32(2*s.qx*s.qz + 2*s.qy*s.qw), float32(s.tx),
		float32(2*s.qx*s.qy + 2*s.qz*s.qw), float32((1 - 2*s.qx*s.qx - 2*s.qz*s.qz) * s.sy), float32(2*s.qy*s.qz - 2*s.qx*s.qw), float32(s.ty),
		float32(2*s.qx*s.qz - 2*s.qy*s.qw), float32(2*s.qy*s.qz + 2*s.qx*s.qw), float32((1 - 2*s.qx*s.qx - 2*s.qy*s.qy) * s.sz), float32(s.tz),
		0, 0, 0, 1,
	}
}
