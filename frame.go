package main

import (
	"fmt"
)

// Block represents the physical state of a single voxel. The zero value
// is empty space.
type Block struct {
	Id   uint
	Data uint
}

const ncx, ncy, ncz = 16, 16, 16

type chunk [ncx][ncy][ncz]Block

type pos struct {
	x, y, z int
}

func (p pos) String() string {
	return fmt.Sprintf("(%d, %d, %d)", p.x, p.y, p.z)
}

// Frame represents a reference frame, specifying a coordinate system
// defined by an SQT transformation and storing the voxel data associated
// with that frame.
type Frame struct {
	Transform *SQT
	chunks    map[pos]chunk
}

func NewFrame() *Frame {
	return &Frame{
		NewSQT(),
		make(map[pos]chunk),
	}
}

// Block returns the Block at local voxel coordinates (x, y, z)
func (f *Frame) Block(x, y, z int) Block {
	p := pos{x / ncx, y / ncy, z / ncz}
	c := f.chunks[p]
	return c[x%ncx][y%ncy][z%ncz]
}

// Block changes the Block at local voxel coordinates (x, y, z)
func (f *Frame) SetBlock(x, y, z int, b Block) {
	p := pos{x / ncx, y / ncy, z / ncz}
	c := f.chunks[p]
	c[x%ncx][y%ncy][z%ncz] = b
	if b.IsEmpty() && c.isEmpty() {
		delete(f.chunks, p)
	} else {
		f.chunks[p] = c
	}
}

// IsEmpty returns true if the Block represents empty space, and
// false otherwise.
func (b Block) IsEmpty() bool {
	return b.Id == 0
}

func (c chunk) isEmpty() bool {
	for _, p := range c {
		for _, r := range p {
			for _, b := range r {
				if !b.IsEmpty() {
					return false
				}
			}
		}
	}
	return true
}
