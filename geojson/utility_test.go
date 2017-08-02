package geojson

import "testing"

func TestUnionBbox(t *testing.T) {
	bb1 := Bbox{0, 0, 2, 4}
	bb2 := Bbox{1, 1, 3, 2}
	bb3 := Bbox{-1, 1, 1, 3}
	bboxes := []*Bbox{&bb1, &bb2, &bb3}
	ubb, err := unionBbox(bboxes)
	if err != nil {
		t.Error()
	}
	if ubb.xmin != -1 || ubb.ymin != 0 || ubb.xmax != 3 || ubb.ymax != 4 {
		t.Fail()
	}
}
