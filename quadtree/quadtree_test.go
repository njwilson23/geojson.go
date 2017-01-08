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
}

func TestOverlaps(t *testing.T) {
	var tf bool
	tf = overlaps(&[4]float64{0, 0, 1, 1}, &[4]float64{0.9, 0.8, 1.9, 1.8})
	if !tf {
		t.Fail()
	}

	tf = overlaps(&[4]float64{0.5, -0.5, 0.6, 0.5}, &[4]float64{0, 0, 1, 1})
	if !tf {
		t.Fail()
	}

	tf = overlaps(&[4]float64{0, 0, 1, 1}, &[4]float64{1.1, 0.8, 1.9, 1.8})
	if tf {
		t.Fail()
	}

	tf = overlaps(&[4]float64{0, 0, 1, 1}, &[4]float64{0.9, 1.2, 1.9, 1.8})
	if tf {
		t.Fail()
	}
}

func TestGet(t *testing.T) {
	// build a tree, insert a point, extend tree, and then retrieve point
	quadtree := new(QuadTree)
	quadtree.MaxChildren = 10

	// temporary because Bbox growing not implemented
	quadtree.Root = new(Node)
	quadtree.Bbox = [4]float64{0, 0, 1, 1}

	r := rand.New(rand.NewSource(49))

	var pt Point
	var err error
	for i := 0; i != 500; i++ {
		pt = Point{r.Float64(), r.Float64()}
		err = quadtree.Insert(pt, i)
		if err != nil {
			fmt.Println(err)
			t.Error()
		}
	}
	err = quadtree.Insert(Point{0.75, 0.3}, 501)
	if err != nil {
		fmt.Println(err)
		t.Error()
	}
	for i := 502; i != 1000; i++ {
		pt = Point{r.Float64(), r.Float64()}
		err = quadtree.Insert(pt, i)
		if err != nil {
			fmt.Println(err)
			t.Error()
		}
	}

	label, err := quadtree.Get(Point{0.75, 0.3})
	if err != nil {
		fmt.Println(err)
		t.Error()
	}
	if label != 501 {
		fmt.Printf("label was %d, but expected %d\n", label, 501)
		t.Fail()
	}
}

func TestWithin(t *testing.T) {
	var pt Point

	pt = Point{0.5, 0.5}
	if pt.Within(&[4]float64{0, 0, 1, 1}) != true {
		t.Fail()
	}

	pt = Point{0.0, 0.0}
	if pt.Within(&[4]float64{0, 0, 1, 1}) != true {
		t.Fail()
	}

	pt = Point{1.0, 1.0}
	if pt.Within(&[4]float64{0, 0, 1, 1}) == true {
		t.Fail()
	}

	pt = Point{0.0, 1.5}
	if pt.Within(&[4]float64{0, 0, 1, 1}) == true {
		t.Fail()
	}
}

func TestSelect(t *testing.T) {

	quadtree := new(QuadTree)
	quadtree.MaxChildren = 100
	quadtree.Root = new(Node)
	quadtree.Bbox = [4]float64{0, 0, 1, 1}

	i := 0
	for x := 0.0; x < 1.0; x = x + 0.02 {
		for y := 0.0; y < 1.0; y = y + 0.02 {
			quadtree.Insert(Point{x, y}, i)
			i++
		}
	}

	labels, err := quadtree.Select(&[4]float64{0.1999, 0.3999, 0.3, 0.6})
	if err != nil {
		fmt.Println(err)
		t.Error()
	}

	if len(labels) != 50 {
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
