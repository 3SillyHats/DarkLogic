package main

import (
	"math"
	"testing"
)

func TestNewFrame(t *testing.T) {
	f := NewFrame()

	if !checkFrame(t, f.Transform,
		1, 0, 0,
		0, 1, 0,
		0, 0, 1,
	) {
		t.Error("NewFrame did not return frame with identity transform")
	}

	for pos := range f.chunks {
		if f.chunks[pos].isEmpty() {
			t.Error("NewFrame stored empty chunk at " + pos.String())
		} else {
			t.Error("NewFrame returned nonempty chunk at " + pos.String())
		}
	}
}

func checkFrame(
	t *testing.T, s *SQT,
	xpx, xpy, xpz float64,
	ypx, ypy, ypz float64,
	zpx, zpy, zpz float64,
) bool {
	xppx, xppy, xppz := s.TransformAbs(1, 0, 0)
	if math.Abs(xppx-xpx) > 1e-15 || math.Abs(xppy-xpy) > 1e-15 || math.Abs(xppz-xpz) > 1e-15 {
		return false
	}

	yppx, yppy, yppz := s.TransformAbs(0, 1, 0)
	if math.Abs(yppx-ypx) > 1e-15 || math.Abs(yppy-ypy) > 1e-15 || math.Abs(yppz-ypz) > 1e-15 {
		return false
	}

	zppx, zppy, zppz := s.TransformAbs(0, 0, 1)
	if math.Abs(zppx-zpx) > 1e-15 || math.Abs(zppy-zpy) > 1e-15 || math.Abs(zppz-zpz) > 1e-15 {
		return false
	}

	return true
}

func TestBlocks(t *testing.T) {
	f := NewFrame()

	for pos := range f.chunks {
		if f.chunks[pos].isEmpty() {
			t.Error("NewFrame stored empty chunk at " + pos.String())
		} else {
			t.Error("NewFrame returned nonempty chunk at " + pos.String())
		}
	}

	if !f.Block(5, 5, 6).IsEmpty() {
		t.Error("NewFrame returned non-empty block at (5,5,6)")
	}
	f.SetBlock(5, 5, 6, Block{1, 0})
	if f.Block(5, 5, 6).IsEmpty() || f.Block(5, 5, 6).Id != 1 {
		t.Error("Block did not return block passed to SetBlock at (5,5,6)")
	}

	chunk := false
	for pos := range f.chunks {
		if f.chunks[pos].isEmpty() {
			t.Error("Empty chunk at " + pos.String())
		} else if pos.x == 0 && pos.y == 0 && pos.z == 0 {
			chunk = true
		}
	}
	if !chunk {
		t.Error("No chunk for block at (5,5,6)")
	}

	f.SetBlock(5, 5, 6, Block{0, 0})
	if !f.Block(5, 5, 6).IsEmpty() {
		t.Error("SetBlock did not clear block at (5,5,6)")
	}

	for pos := range f.chunks {
		if f.chunks[pos].isEmpty() {
			t.Error("SetBlock did not clear chunk at " + pos.String())
		} else if pos.x == 0 && pos.y == 0 && pos.z == 0 {
			t.Error("SetBlock did not empty chunk at " + pos.String())
		}
	}
}
