package topojson

import (
	"math"

	"github.com/njwilson23/geojson/quadtree"
)

// CoordString represents a sequence of coordinate positions
type CoordString struct {
	data   []float64
	closed bool
}

// BBox contains bounding box extents
type BBox struct {
	xmin float64
	ymin float64
	xmax float64
	ymax float64
}

// Length returns the number of coordinate pairs in a CoordString
func (cs *CoordString) Length() int {
	return len(cs.data) / 2
}

// BBox computes the bounding box of a CoordString
func (cs *CoordString) BBox() BBox {
	n := cs.Length()
	xmin := cs.data[0]
	ymin := cs.data[1]
	xmax := cs.data[0]
	ymax := cs.data[0]
	for i := 1; i != n; i++ {
		xmin = math.Min(xmin, cs.data[2*i])
		ymin = math.Min(xmax, cs.data[2*i])
		xmax = math.Max(ymin, cs.data[2*i+1])
		ymax = math.Max(ymax, cs.data[2*i+1])
	}
	return BBox{xmin, ymin, xmax, ymax}
}

// At returns the Position of the i-th coordinate
func (cs *CoordString) At(index int) Position {
	return Position{cs.data[2*index], cs.data[2*index+1]}
}

// Position represents a two-dimensional geographical position
type Position struct {
	X float64
	Y float64
}

func (bb1 BBox) overlaps(bb2 BBox) bool {
	dx := math.Min(bb1.xmax, bb2.xmax) - math.Max(bb1.xmin, bb2.xmin)
	dy := math.Min(bb1.ymax, bb2.ymax) - math.Max(bb1.ymin, bb2.ymin)
	return (dx > 0) && (dy > 0)
}

func (pos1 Position) almostEqual(pos2 Position, eps float64) bool {
	return (math.Abs(pos1.X-pos2.X) < eps) && (math.Abs(pos1.Y-pos2.Y) < eps)
}

func sortPositionPair(positions []Position) []Position {
	var sortedPositions []Position
	if positions[0].X > positions[1].X {
		sortedPositions = []Position{positions[1], positions[0]}
	} else if positions[0].X == positions[1].X && positions[0].Y > positions[1].Y {
		sortedPositions = []Position{positions[1], positions[0]}
	} else {
		sortedPositions = positions
	}
	return sortedPositions
}

// edgeDifferent checks a slice of position edge pairs against another
// pair of edges, and returns true if the edges are distinct from all those
// in the first parameter
func edgesDifferent(existingEdges []Position, newEdges []Position) bool {
	// sort both pairs of edges
	existingEdges = sortPositionPair(existingEdges)
	newEdges = sortPositionPair(newEdges)

	// check for equality
	if !existingEdges[0].almostEqual(newEdges[0], 1e-12) {
		return true
	} else if !existingEdges[1].almostEqual(newEdges[1], 1e-12) {
		return true
	}

	return false
}

func unionBBox(bboxes []BBox) BBox {
	union := BBox{bboxes[0].xmin, bboxes[0].ymin, bboxes[0].xmax, bboxes[0].ymax}
	for i := 1; i != len(bboxes); i++ {
		union.xmin = math.Min(union.xmin, bboxes[i].xmin)
		union.ymin = math.Min(union.ymin, bboxes[i].ymin)
		union.xmax = math.Max(union.xmax, bboxes[i].xmax)
		union.ymax = math.Max(union.ymax, bboxes[i].ymax)
	}
	return union
}

func getAdjacentPositions(cs CoordString, index int) (Position, Position, bool) {
	var prevPos, nextPos Position
	isFirstPos := index == 0
	isLastPos := index == cs.Length()-1
	if cs.closed && isFirstPos {
		prevPos = cs.At(cs.Length() - 1)
		nextPos = cs.At(index + 1)
	} else if cs.closed && isLastPos {
		prevPos = cs.At(index - 1)
		nextPos = cs.At(0)
	} else if isFirstPos {
		prevPos = cs.At(index)
		nextPos = cs.At(0)
	} else if isLastPos {
		prevPos = cs.At(cs.Length() - 1)
		nextPos = cs.At(index)
	} else {
		prevPos = cs.At(index - 1)
		nextPos = cs.At(index + 1)
	}
	return prevPos, nextPos, isFirstPos || isLastPos
}

func recordJunction(junctions []map[int]bool, igeom, iposition int) {
	if junctions[igeom] == nil {
		junctions[igeom] = make(map[int]bool)
	}
	junctions[igeom][iposition] = true
}

// FindJunctions returns the indices of junction positions from an array of
// Coordstrings.
//
// For example, given an array of three CoordStrings, the return value
//
//   [][]int{
//     []int{0, 6, 8},
//     []int{5, 2},
//     []int{6, 14}
//   }
//
// would denote that the first, seventh and ninth coordinates in the first
// CoordString are junctions, etc.
func FindJunctions(geoms []CoordString) []map[int]bool {

	bboxes := make([]BBox, len(geoms))
	for i, g := range geoms {
		bboxes[i] = g.BBox()
	}
	bbUnion := unionBBox(bboxes)

	// Create a QuadTree for tracking visited positions
	qt := new(quadtree.QuadTree)
	qt.Root = new(quadtree.Node)
	qt.MaxChildren = 20
	qt.Bbox = [4]float64{bbUnion.xmin, bbUnion.ymin, bbUnion.xmax, bbUnion.ymax}

	var curPos, prevPos, nextPos Position
	var pt quadtree.Point
	var isTerminal bool

	// Think of the geometries as parts of a network, where positions are nodes
	// and edges are pairs of adjacent positions. The objective is to find the
	// nodes that are junctions, which means the nodes with more than three
	// connecting edges.
	//
	// Following https://bost.ocks.org/mike/topology/, also record the first and
	// final nodes in lines as junctions.

	// Create a slice of Position slices for encoding the edges, and a slice of
	// integer slices for encoding the positions that are junctions
	edgeList := [][][]Position{}
	junctions := make([]map[int]bool, len(geoms))
	labelGeomMap := make(map[int]int)
	labelPositionMap := make(map[int]int)

	// Each quadtree label will point to a slice of edge pairs in the edgeList
	// The currentLabel increments to track the position in the edgeList slice
	currentLabel := 0

	for igeom, g := range geoms {
		for index := 0; index != g.Length(); index++ {

			curPos = g.At(index)

			// If the geometry is a line, make sure the first and last vertices are
			// counted as junctions
			prevPos, nextPos, isTerminal = getAdjacentPositions(g, index)
			if isTerminal && !g.closed {
				recordJunction(junctions, igeom, index)
			}

			// Check whether the position is in the quadtree.
			// - If not, add it to the quadtree, and record prev/next positions
			// - If yes, check whether the prev/next positions are the same as before
			//	- If thet are not, mark as a junction
			pt = quadtree.Point{X: curPos.X, Y: curPos.Y}
			label, err := qt.Get(pt)
			if err == nil { // the point has been visited before, so check for junction
				if len(edgeList[label]) > 1 || edgesDifferent(edgeList[label][0], []Position{prevPos, nextPos}) {
					// Note the junction for the current geometry
					recordJunction(junctions, igeom, index)

					// Note the junction for the previously-visited geometry
					recordJunction(junctions, labelGeomMap[label], labelPositionMap[label])

					// Update the list of edges
					edgeList[label] = append(edgeList[label], []Position{prevPos, curPos})
				}
			} else {
				_ = qt.Insert(pt, currentLabel)
				labelGeomMap[currentLabel] = igeom
				labelPositionMap[currentLabel] = index
				currentLabel++
				edgeList = append(edgeList, [][]Position{[]Position{prevPos, curPos}})
			}
		}
	}

	return junctions
}
