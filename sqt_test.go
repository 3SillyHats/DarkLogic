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
	s.SetScale(-1)
	checkTransform(t, "SQT", "rotation by pi/2 about x, translation by (+2, +1, -3) and scale by *-1", s,
		[16]float32{
			-1, 0, 0, 2,
			0, 0, 1, 1,
			0, -1, 0, -3,
			0, 0, 0, 1,
		},
		1, 1, -3,
		2, 1, -4,
		2, 2, -3,
	)

	s.SetRotation(3*math.Pi/2, 0, 1, 0)
	s.SetTranslation(0, 4, 0)
	s.SetScale(5)
	checkTransform(t, "SQT", "rotation by 3pi/2 about y, translation by (0, +4, +0) and scale by *5", s,
		[16]float32{
			0, 0, -5, 0,
			0, 5, 0, 4,
			5, 0, 0, 0,
			0, 0, 0, 1,
		},
		0, 4, 5,
		0, 9, 0,
		-5, 4, 0,
	)
}

func TestCompose(t *testing.T) {
	s1 := NewSQT()
	s1.SetRotation(math.Pi/2, 1, 0, 0)
	s1.SetTranslation(2, 1, -3)
	s1.SetScale(0.5)

	s2 := NewSQT()
	s2.SetRotation(3*math.Pi/2, 0, 1, 0)
	s2.SetTranslation(0, 4, 0)
	s2.SetScale(-3)

	s3 := s2.Compose(s1)

	x1x, x1y, x1z := s1.Transform(1, 0, 0)
	y1x, y1y, y1z := s1.Transform(0, 1, 0)
	z1x, z1y, z1z := s1.Transform(0, 0, 1)

	x2x, x2y, x2z := s2.Transform(x1x, x1y, x1z)
	y2x, y2y, y2z := s2.Transform(y1x, y1y, y1z)
	z2x, z2y, z2z := s2.Transform(z1x, z1y, z1z)

	x3x, x3y, x3z := s3.Transform(1, 0, 0)
	y3x, y3y, y3z := s3.Transform(0, 1, 0)
	z3x, z3y, z3z := s3.Transform(0, 0, 1)

	if math.Abs(x3x-x2x) > 1e-14 || math.Abs(x3y-x2y) > 1e-14 || math.Abs(x3z-x2z) > 1e-14 {
		t.Errorf("%s %d,%d,%d", "Compose did not transform x the same as sequential transformation", math.Abs(x2x-x3x), math.Abs(x2y-x3y), math.Abs(x2z-x3z))
	}
	if math.Abs(y3x-y2x) > 1e-14 || math.Abs(y3y-y2y) > 1e-14 || math.Abs(y3z-y2z) > 1e-14 {
		t.Errorf("%s %d,%d,%d", "Compose did not transform y the same as sequential transformation", math.Abs(y2x-y3x), math.Abs(y2y-y3y), math.Abs(y2z-y3z))
	}
	if math.Abs(z3x-z2x) > 1e-14 || math.Abs(z3y-z2y) > 1e-14 || math.Abs(z3z-z2z) > 1e-14 {
		t.Errorf("%s %d,%d,%d vs %d %d %d", "Compose did not transform z the same as sequential transformation", z2x, z2y, z2z, z3x, z3y, z3z)
	}
}

func TestInverse(t *testing.T) {
	s := NewSQT()
	s.SetRotation(math.Pi/2, 1, 0, 0)
	s.SetTranslation(2, 1, -3)
	s.SetScale(-1)
	i := s.Compose(s.Inverse())
	checkTransform(t, "Inverse", "rotation by pi/2 about x, translation by (+2, +1, -3) and scale by *-1", i,
		[16]float32{
			1, 0, 0, 0,
			0, 1, 0, 0,
			0, 0, 1, 0,
			0, 0, 0, 1,
		},
		1, 0, 0,
		0, 1, 0,
		0, 0, 1,
	)

	s.SetRotation(3*math.Pi/2, 0, 1, 0)
	s.SetTranslation(0, 4, 0)
	s.SetScale(5)
	i = s.Compose(s.Inverse())
	checkTransform(t, "Inverse", "rotation by 3pi/2 about y, translation by (0, +4, +0) and scale by *5", i,
		[16]float32{
			1, 0, 0, 0,
			0, 1, 0, 0,
			0, 0, 1, 0,
			0, 0, 0, 1,
		},
		1, 0, 0,
		0, 1, 0,
		0, 0, 1,
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
