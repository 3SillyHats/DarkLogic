package main

import (
	"math"
	"testing"
)

func TestNew(t *testing.T) {
	s := NewSQT()
	m := s.Matrix()
	identity := [16]float32{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
	if m != identity {
		t.Error("NewSQT did not return identity transformation")
	}
}

func TestRotate(t *testing.T) {
	s := NewSQT()
	s.SetRotation(math.Pi/2, 1, 0, 0)
	checkTransform(t, "SetRotation", "pi/2 about x", s,
		[16]float32{
			1, 0, 0, 0,
			0, 0, -1, 0,
			0, 1, 0, 0,
			0, 0, 0, 1,
		},
		1, 0, 0,
		0, 0, 1,
		0, -1, 0,
	)

	s.SetRotation(math.Pi, 0, 1, 0)
	checkTransform(t, "SetRotation", "pi about y", s,
		[16]float32{
			-1, 0, 0, 0,
			0, 1, 0, 0,
			0, 0, -1, 0,
			0, 0, 0, 1,
		},
		-1, 0, 0,
		0, 1, 0,
		0, 0, -1,
	)

	s.SetRotation(3*math.Pi/2, 0, 0, -1)
	checkTransform(t, "SetRotation", "3pi/2 about -z", s,
		[16]float32{
			0, -1, 0, 0,
			1, 0, 0, 0,
			0, 0, 1, 0,
			0, 0, 0, 1,
		},
		0, 1, 0,
		-1, 0, 0,
		0, 0, 1,
	)

	s.Rotate(math.Pi/2, -1, 0, 0)
	checkTransform(t, "Rotate", "3pi/2 about -z then pi/2 about x", s,
		[16]float32{
			0, -1, 0, 0,
			0, 0, 1, 0,
			-1, 0, 0, 0,
			0, 0, 0, 1,
		},
		0, 0, -1,
		-1, 0, 0,
		0, 1, 0,
	)
}

func TestTranslate(t *testing.T) {
	s := NewSQT()
	s.SetTranslation(5, 0, 0)
	checkTransform(t, "SetTranslation", "(+5, 0, 0)", s,
		[16]float32{
			1, 0, 0, 5,
			0, 1, 0, 0,
			0, 0, 1, 0,
			0, 0, 0, 1,
		},
		6, 0, 0,
		5, 1, 0,
		5, 0, 1,
	)

	s.SetTranslation(0, 1, -1)
	checkTransform(t, "SetTranslation", "(0, +1, -1)", s,
		[16]float32{
			1, 0, 0, 0,
			0, 1, 0, 1,
			0, 0, 1, -1,
			0, 0, 0, 1,
		},
		1, 1, -1,
		0, 2, -1,
		0, 1, 0,
	)

	s.Translate(5, -1, 6)
	checkTransform(t, "Translate", "(0, +1, -1) then (+5, -1, +6)", s,
		[16]float32{
			1, 0, 0, 5,
			0, 1, 0, 0,
			0, 0, 1, 5,
			0, 0, 0, 1,
		},
		6, 0, 5,
		5, 1, 5,
		5, 0, 6,
	)
}

func TestScale(t *testing.T) {
	s := NewSQT()
	s.SetScale(2)
	checkTransform(t, "SetScale", "*2", s,
		[16]float32{
			2, 0, 0, 0,
			0, 2, 0, 0,
			0, 0, 2, 0,
			0, 0, 0, 1,
		},
		2, 0, 0,
		0, 2, 0,
		0, 0, 2,
	)

	s.SetScale(0.5)
	checkTransform(t, "SetScale", "/2", s,
		[16]float32{
			0.5, 0, 0, 0,
			0, 0.5, 0, 0,
			0, 0, 0.5, 0,
			0, 0, 0, 1,
		},
		0.5, 0, 0,
		0, 0.5, 0,
		0, 0, 0.5,
	)

	s.Scale(3)
	checkTransform(t, "Scale", "/2 then *3", s,
		[16]float32{
			1.5, 0, 0, 0,
			0, 1.5, 0, 0,
			0, 0, 1.5, 0,
			0, 0, 0, 1,
		},
		1.5, 0, 0,
		0, 1.5, 0,
		0, 0, 1.5,
	)
}

func TestOrder(t *testing.T) {
	s := NewSQT()
	s.SetRotation(math.Pi/2, 1, 0, 0)
	s.SetTranslation(2, 1, -3)
	s.SetScale(0.5, 3, -1)
	checkTransform(t, "SQT", "rotation by pi/2 about x, translation by (+2, +1, -3) and scale by (/2, *3, *-1)", s,
		[16]float32{
			0.5, 0, 0, 2,
			0, 0, 1, 1,
			0, 3, 0, -3,
			0, 0, 0, 1,
		},
		2.5, 1, -3,
		2, 1, 0,
		2, 2, -3,
	)

	s.SetRotation(3*math.Pi/2, 0, 1, 0)
	s.SetTranslation(0, 4, 0)
	s.SetScale(1, -2, 5)
	checkTransform(t, "SQT", "rotation by 3pi/2 about y, translation by (0, +4, +0) and scale by (*1, *-2, *5)", s,
		[16]float32{
			0, 0, -5, 0,
			0, -2, 0, 4,
			1, 0, 0, 0,
			0, 0, 0, 1,
		},
		0, 4, 1,
		0, 2, 0,
		-5, 4, 0,
	)
}

func checkTransform(
	t *testing.T, operation, desc string, s *SQT,
	matrix [16]float32,
	xpx, xpy, xpz float64,
	ypx, ypy, ypz float64,
	zpx, zpy, zpz float64,
) {
	m := s.Matrix()

	x := [4]float32{1, 0, 0, 1}
	y := [4]float32{0, 1, 0, 1}
	z := [4]float32{0, 0, 1, 1}

	if diff(m[:], matrix[:]) > 1e-16 {
		t.Error(operation + " did not return correct transformation for " + desc)
	}

	mx := matmul(m, x)
	xp := [4]float32{float32(xpx), float32(xpy), float32(xpz), 1}
	if diff(xp[:], mx[:]) > 1e-16 {
		t.Error(operation + " did not return expected transformation for x for " + desc)
	}

	my := matmul(m, y)
	yp := [4]float32{float32(ypx), float32(ypy), float32(ypz), 1}
	if diff(yp[:], my[:]) > 1e-16 {
		t.Error(operation + " did not return expected transformation for y for " + desc)
	}

	mz := matmul(m, z)
	zp := [4]float32{float32(zpx), float32(zpy), float32(zpz), 1}
	if diff(zp[:], mz[:]) > 1e-16 {
		t.Error(operation + " did not return expected transformation for z for " + desc)
	}

	xppx, xppy, xppz := s.Transform(1, 0, 0)
	if math.Abs(xppx-xpx) > 1e-15 || math.Abs(xppy-xpy) > 1e-15 || math.Abs(xppz-xpz) > 1e-15 {
		t.Error(operation + " did not transform x as expected for " + desc)
	}

	yppx, yppy, yppz := s.Transform(0, 1, 0)
	if math.Abs(yppx-ypx) > 1e-15 || math.Abs(yppy-ypy) > 1e-15 || math.Abs(yppz-ypz) > 1e-15 {
		t.Error(operation + " did not transform y as expected for " + desc)
	}

	zppx, zppy, zppz := s.Transform(0, 0, 1)
	if math.Abs(zppx-zpx) > 1e-15 || math.Abs(zppy-zpy) > 1e-15 || math.Abs(zppz-zpz) > 1e-15 {
		t.Error(operation + " did not transform z as expected for " + desc)
	}
}

func diff(a, b []float32) float32 {
	var diff float32
	for i := 0; i < len(a) && i < len(b); i++ {
		diff += (a[i] - b[i]) * (a[i] - b[i])
	}
	return diff
}

func matmul(M [16]float32, v [4]float32) [4]float32 {
	return [4]float32{
		M[0]*v[0] + M[1]*v[1] + M[2]*v[2] + M[3]*v[3],
		M[4]*v[0] + M[5]*v[1] + M[6]*v[2] + M[7]*v[3],
		M[8]*v[0] + M[9]*v[1] + M[10]*v[2] + M[11]*v[3],
		M[12]*v[0] + M[13]*v[1] + M[14]*v[2] + M[15]*v[3],
	}
}
