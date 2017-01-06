package quadtree

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestAddPoints(t *testing.T) {
	quadtree := new(QuadTree)
	quadtree.MaxChildren = 5

	// temporary because Bbox growing not implemented
	quadtree.Root = new(Node)
	quadtree.Bbox = [4]float64{0, 0, 50, 50}

	r := rand.New(rand.NewSource(49))

	var pt Point
	var err error
	for i := 0; i != 1000; i++ {
		pt = Point{50.0 * r.Float64(), 50.0 * r.Float64()}
		err = quadtree.Insert(pt, i)
		if err != nil {
			fmt.Println(err)
			t.Error()
		}
	}
	fmt.Println("tree depth:", quadtree.Depth())
}

func TestOverlaps1(t *testing.T) {
	tf := overlaps(&[4]float64{0, 0, 1, 1}, &[4]float64{0.9, 0.8, 1.9, 1.8})
	if !tf {
		t.Fail()
	}
}
func TestOverlaps2(t *testing.T) {
	tf := overlaps(&[4]float64{0.5, -0.5, 0.6, 0.5}, &[4]float64{0, 0, 1, 1})
	if !tf {
		t.Fail()
	}
}
func TestOverlaps3(t *testing.T) {
	tf := overlaps(&[4]float64{0, 0, 1, 1}, &[4]float64{1.1, 0.8, 1.9, 1.8})
	if tf {
		t.Fail()
	}
}
func TestOverlaps4(t *testing.T) {
	tf := overlaps(&[4]float64{0, 0, 1, 1}, &[4]float64{0.9, 1.2, 1.9, 1.8})
	if tf {
		t.Fail()
	}
}

func BenchmarkBuildRandom(b *testing.B) {
	quadtree := new(QuadTree)
	quadtree.MaxChildren = 50

	// temporary because Bbox growing not implemented
	quadtree.Root = new(Node)
	quadtree.Bbox = [4]float64{0, 0, 50, 50}

	r := rand.New(rand.NewSource(49))

	var pt Point
	var err error

	b.ResetTimer()
	for i := 0; i != 1000000; i++ {
		pt = Point{50.0 * r.Float64(), 50.0 * r.Float64()}
		err = quadtree.Insert(pt, i)
		if err != nil {
			b.Error()
		}
	}
}
