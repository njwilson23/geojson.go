package geojson

import (
	"testing"
)

func TestIsLeft(t *testing.T) {
	var q, p0, p1 []float64

	q = []float64{0, 1}
	p0 = []float64{1, 0}
	p1 = []float64{2, 2}
	if !isLeft(q, p0, p1) {
		t.Fail()
	}

	q = []float64{2, 1}
	p0 = []float64{1, 0}
	p1 = []float64{2, 2}
	if isLeft(q, p0, p1) {
		t.Fail()
	}

	q = []float64{0, 1}
	p0 = []float64{1, 0}
	p1 = []float64{2, 0.5}
	if !isLeft(q, p0, p1) {
		t.Fail()
	}
}

func TestRingWinding(t *testing.T) {
	var ring [][]float64

	// clockwise ring
	ring = [][]float64{
		[]float64{0, 0},
		[]float64{0, 1},
		[]float64{1, 1},
		[]float64{1, 0},
		[]float64{0, 0}}
	if isCounterClockwise(&ring) {
		t.Fail()
	}

	// counter clockwise ring
	ring = [][]float64{
		[]float64{0, 0},
		[]float64{1, 0},
		[]float64{1, 1},
		[]float64{0, 1},
		[]float64{0, 0}}
	if !isCounterClockwise(&ring) {
		t.Fail()
	}
}
