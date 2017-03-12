package topojson

import (
	"fmt"
	"testing"
)

func TestJunctionsCrossedLines(t *testing.T) {
	line1 := CoordString{[]float64{1, 1, 2, 4, 3, 9}, false}
	line2 := CoordString{[]float64{1, 7, 2, 4, 3, -1, 4, -2}, false}

	junctions := FindJunctions([]CoordString{line1, line2})
	fmt.Println(junctions)
}

func TestJunctionsAdjacentSquares(t *testing.T) {
	square1 := CoordString{[]float64{0, 0, 1, 0, 1, 1, 0, 1}, true}
	square2 := CoordString{[]float64{0, 1, 1, 1, 1, 2, 0, 2}, true}

	junctions := FindJunctions([]CoordString{square1, square2})
	fmt.Println(junctions)
}
