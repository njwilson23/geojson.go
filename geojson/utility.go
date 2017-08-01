package geojson

import (
	"fmt"
	"math"
)

func unionBbox(bboxes []*Bbox) (*Bbox, error) {
	bb := new(Bbox)
	if len(bboxes) == 0 {
		return bb, fmt.Errorf("union of empty set")
	}
	*bb = *bboxes[0]
	for _, tmp := range bboxes {
		bb.xmin = math.Min(bb.xmin, tmp.xmin)
		bb.ymin = math.Min(bb.ymin, tmp.ymin)
		bb.xmax = math.Max(bb.xmax, tmp.xmax)
		bb.ymax = math.Max(bb.ymax, tmp.ymax)
	}
	return bb, nil
}
