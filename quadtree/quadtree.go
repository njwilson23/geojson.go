package quadtree

import (
	"errors"
	"fmt"
)

// QuadTree represents a quadtree structure and contains a pointer to the trunk node
type QuadTree struct {
	Root        *Node
	MaxChildren int
	Bbox        [4]float64
}

// Node represents a quadtree node, containing either an array of positions and labels, or child nodes
type Node struct {
	Points []Point
	Labels []int
	LL     *Node
	LR     *Node
	UL     *Node
	UR     *Node
}

// Point represents a two-dimensional position
type Point struct {
	X float64
	Y float64
}

// Within returns true if the Point is inside a bounding box
func (p *Point) Within(bbox *[4]float64) bool {
	return (*bbox)[0] <= p.X && p.X < (*bbox)[2] && (*bbox)[1] <= p.Y && p.Y < (*bbox)[3]
}

// Insert adds a Point and a label. Returns a non-nil error when the tree
// bounding box has not yet been initialized
func (q *QuadTree) Insert(pt Point, label int) error {
	var node *Node
	var err error

	if q.Bbox == [4]float64{0, 0, 0, 0} {
		return errors.New("cannot insert without initializing QuadTree.Bbox")
	}

	if q.Root == nil {
		node = new(Node)
		q.Root = node
	} else {
		node = q.Root
	}

	var tmpPt Point
	var tmpLabel int

	nodeBbox := q.Bbox

	// Grow if necessary
	//if n.Bbox[0] > pt.X || n.Bbox[1] > pt.Y || n.Bbox[2] < pt.X || n.Bbox[3] < pt.Y {
	//	n.GrowBbox(pt)
	//}

	for {

		if node.LL != nil {
			// Insert into a child node
			xmid := 0.5 * (nodeBbox[0] + nodeBbox[2])
			ymid := 0.5 * (nodeBbox[1] + nodeBbox[3])
			if pt.X < xmid {
				if pt.Y < ymid {
					node = node.LL
					nodeBbox = [4]float64{nodeBbox[0], nodeBbox[1], xmid, ymid}
				} else {
					node = node.UL
					nodeBbox = [4]float64{nodeBbox[0], ymid, xmid, nodeBbox[3]}
				}
			} else {
				if pt.Y < ymid {
					node = node.LR
					nodeBbox = [4]float64{xmid, nodeBbox[1], nodeBbox[2], ymid}
				} else {
					node = node.UR
					nodeBbox = [4]float64{xmid, ymid, nodeBbox[2], nodeBbox[3]}
				}
			}
		} else if len(node.Labels) == q.MaxChildren {
			// Split
			xmid := 0.5 * (nodeBbox[0] + nodeBbox[2])
			ymid := 0.5 * (nodeBbox[1] + nodeBbox[3])

			node.LL = new(Node)
			node.LR = new(Node)
			node.UL = new(Node)
			node.UR = new(Node)

			for i := 0; i != q.MaxChildren; i++ {

				tmpPt = node.Points[0]
				tmpLabel = node.Labels[0]
				if tmpPt.X < xmid {
					if tmpPt.Y < ymid {
						node.LL.Points = append(node.LL.Points, tmpPt)
						node.LL.Labels = append(node.LL.Labels, tmpLabel)
					} else {
						node.UL.Points = append(node.UL.Points, tmpPt)
						node.UL.Labels = append(node.UL.Labels, tmpLabel)
					}
				} else {
					if tmpPt.Y < ymid {
						node.LR.Points = append(node.LR.Points, tmpPt)
						node.LR.Labels = append(node.LR.Labels, tmpLabel)
					} else {
						node.UR.Points = append(node.UR.Points, tmpPt)
						node.UR.Labels = append(node.UR.Labels, tmpLabel)
					}
				}
				node.Points = node.Points[1:]
				node.Labels = node.Labels[1:]
			}
		} else {
			// Append to leaf
			node.Points = append(node.Points, pt)
			node.Labels = append(node.Labels, label)
			break
		}

	}
	return err
}

// Get returns the label of a matching point if present, and returns a non-nil error if the point is missing
func (q *QuadTree) Get(pt Point) (int, error) {

	node := q.Root
	nodeBbox := q.Bbox

	if pt.X < nodeBbox[0] || pt.X > nodeBbox[2] || pt.Y < nodeBbox[1] || pt.Y > nodeBbox[3] {
		return 0, errors.New("missing")
	}

	for {
		if node.LL != nil {
			xmid := 0.5 * (nodeBbox[0] + nodeBbox[2])
			ymid := 0.5 * (nodeBbox[1] + nodeBbox[3])
			if pt.X < xmid {
				if pt.Y < ymid {
					node = node.LL
					nodeBbox = [4]float64{nodeBbox[0], nodeBbox[1], xmid, ymid}
				} else {
					node = node.UL
					nodeBbox = [4]float64{nodeBbox[0], ymid, xmid, nodeBbox[3]}
				}
			} else {
				if pt.Y < ymid {
					node = node.LR
					nodeBbox = [4]float64{xmid, nodeBbox[1], nodeBbox[2], ymid}
				} else {
					node = node.UR
					nodeBbox = [4]float64{xmid, ymid, nodeBbox[2], nodeBbox[3]}
				}
			}
		} else {
			for i := 0; i != len(node.Labels); i++ {
				if pt == node.Points[i] {
					return node.Labels[i], nil
				}
			}
			return 0, errors.New("missing")
		}
	}
}

type nodebox struct {
	node *Node
	bbox *[4]float64
}

// Select returns the labels within a bounding box
func (q *QuadTree) Select(bbox *[4]float64) ([]int, error) {
	var node *Node
	var err error
	var nodeBbox [4]float64
	var tmpBbox *[4]float64
	var xmid, ymid float64

	var labels []int
	//labels := make([]int, 0)
	nodeStack := make([]nodebox, 1)
	nodeStack[0] = nodebox{q.Root, &q.Bbox}

	for len(nodeStack) != 0 {
		node = nodeStack[0].node
		nodeBbox = *nodeStack[0].bbox
		nodeStack = nodeStack[1:]
		if node.LL == nil {
			for i := 0; i != len(node.Labels); i++ {
				if node.Points[i].Within(bbox) {
					labels = append(labels, node.Labels[i])
				}
			}
		} else {
			xmid = 0.5 * (nodeBbox[0] + nodeBbox[2])
			ymid = 0.5 * (nodeBbox[1] + nodeBbox[3])

			tmpBbox = &[4]float64{nodeBbox[0], nodeBbox[1], xmid, ymid}
			if overlaps(bbox, tmpBbox) {
				nodeStack = append(nodeStack, nodebox{node.LL, tmpBbox})
			}

			tmpBbox = &[4]float64{nodeBbox[0], ymid, xmid, nodeBbox[3]}
			if overlaps(bbox, tmpBbox) {
				nodeStack = append(nodeStack, nodebox{node.UL, tmpBbox})
			}

			tmpBbox = &[4]float64{xmid, nodeBbox[1], nodeBbox[2], ymid}
			if overlaps(bbox, tmpBbox) {
				nodeStack = append(nodeStack, nodebox{node.LR, tmpBbox})
			}

			tmpBbox = &[4]float64{xmid, ymid, nodeBbox[2], nodeBbox[3]}
			if overlaps(bbox, tmpBbox) {
				nodeStack = append(nodeStack, nodebox{node.UR, tmpBbox})
			}
		}
	}
	return labels, err
}

// Depth returns the total depth of a QuadTree
func (q *QuadTree) Depth() int {
	var node *Node
	var depth, maxdepth int

	nodeStack := make([]*Node, 1)
	depthStack := make([]int, 1)
	nodeStack[0] = q.Root
	depthStack[0] = 1
	for len(nodeStack) != 0 {
		node = nodeStack[0]
		depth = depthStack[0]
		nodeStack = nodeStack[1:]
		depthStack = depthStack[1:]
		if node.LL == nil {
			if maxdepth < depth {
				maxdepth = depth
			}
		} else {
			nodeStack = append(nodeStack, node.LL, node.LR, node.UL, node.UR)
			depthStack = append(depthStack, depth+1, depth+1, depth+1, depth+1)
		}
	}
	return maxdepth
}

func overlaps(bb0, bb1 *[4]float64) bool {
	return !(bb0[0] > bb1[2] || bb0[2] < bb1[0] || bb0[1] > bb1[3] || bb0[3] < bb1[1])
}

//func (q *QuadTree) Delete(pt Point) (int, error) {
//}

// String returns a position in the format '(x,y)'
func (p Point) String() string {
	return fmt.Sprintf("(%.4f,%.4f)", p.X, p.Y)
}
