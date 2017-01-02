/* functions to check the orientation of polygon rings */
package geojson

func isLeft(q []float64, p0 []float64, p1 []float64) bool {
	return ((p1[0]-p0[0])*(q[1]-p0[1]) - (q[0]-p0[0])*(p1[1]-p0[1])) > 0
}

func isCounterClockwise(ring [][]float64) bool {
	var xmin, ymin float64
	trimmedRing := ring[:len(ring)-1]
	xmin = trimmedRing[len(trimmedRing)-1][0]
	ymin = trimmedRing[len(trimmedRing)-1][1]
	imin := len(trimmedRing) - 1

	var pt []float64
	for i := 0; i != len(trimmedRing)-1; i++ {
		pt = trimmedRing[i]
		if pt[1] < ymin || (pt[1] == ymin && pt[0] < xmin) {
			imin = i
			xmin = pt[0]
			ymin = pt[1]
		}
	}

	var prevVertex, nextVertex []float64
	if imin == 0 {
		prevVertex = trimmedRing[len(trimmedRing)-1]
	} else {
		prevVertex = trimmedRing[imin-1]
	}
	if imin == len(trimmedRing)-1 {
		nextVertex = trimmedRing[0]
	} else {
		nextVertex = trimmedRing[imin+1]
	}
	return isLeft(prevVertex, trimmedRing[imin], nextVertex)
}
