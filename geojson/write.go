/*
 * Implements functions for marshaling structs to GeoJSON bytes
 */
package geojson

import "encoding/json"

// MarshalGeometry returns a byte array encoding a GeoJSON geometry
func MarshalGeometry(g Geometry) ([]byte, error) {
	b, err := json.Marshal(g)
	if err != nil {
		return []byte{}, err
	}
	return b, err
}

// MarshalFeature returns a byte array encoding a GeoJSON Feature
func MarshalFeature(f *Feature) ([]byte, error) {
	b, err := json.Marshal(f)
	if err != nil {
		return []byte{}, err
	}
	return b, err
}

// MarshalFeatureCollection returns a byte array encoding a GeoJSON FeatureCollection
func MarshalFeatureCollection(fc *FeatureCollection) ([]byte, error) {
	b, err := json.Marshal(fc)
	if err != nil {
		return []byte{}, err
	}
	return b, err
}

/* MarshalJSON methods for all GeoJSON types */
func (pt *Point) MarshalJSON() ([]byte, error) {
	p := struct {
		Point
		Type string `json:"type"`
	}{*pt, "Point"}

	b, err := json.Marshal(p)
	return b, err
}

func (ls *LineString) MarshalJSON() ([]byte, error) {
	l := struct {
		LineString
		Type string `json:"type"`
	}{*ls, "LineString"}

	b, err := json.Marshal(l)
	return b, err
}

// Ring wraps a coordinate ring with its hierarchical level
type Ring struct {
	i  int
	cx [][]float64
}

// ringSender seeds an array of Rings into a channel
func ringSender(ringSlice []Ring, ch chan<- Ring) {
	for _, ring := range ringSlice {
		ch <- ring
	}
	close(ch)
}

// closedEnforcer ensures that rings are closed
func closedEnforcer(input <-chan Ring, output chan<- Ring) {
	// ensure ring closed
	for ring := range input {
		last := ring.cx[len(ring.cx)-1]
		for i, v := range ring.cx[0] {
			if last[i] != v {
				ring.cx = append(ring.cx, ring.cx[0])
				break
			}
		}
		output <- ring
	}
	close(output)
}

// windingEnforcer ensures that rings have the correct winding order
func windingEnforcer(input <-chan Ring, output chan<- Ring) {
	var ccw bool
	for ring := range input {
		// ensure ring winds correctly
		ccw = isCounterClockwise(ring.cx)
		if (ring.i != 0 && ccw) || (ring.i == 0 && !ccw) {
			for j, k := 0, len(ring.cx)-1; j < k; j, k = j+1, k-1 {
				ring.cx[j], ring.cx[k] = ring.cx[k], ring.cx[j]
			}
		}
		output <- ring
	}
	close(output)
}

func (poly *Polygon) MarshalJSON() ([]byte, error) {
	// enforce CCW winding on external rings and CW winding on internal rings
	chWinding := make(chan Ring)
	chClosed := make(chan Ring)
	chDone := make(chan Ring)

	ringSlice := []Ring{}
	for i, slc := range poly.Coordinates {
		ringSlice = append(ringSlice, Ring{i, slc})
	}

	go ringSender(ringSlice, chWinding)
	go windingEnforcer(chWinding, chClosed)
	go closedEnforcer(chClosed, chDone)

	coordinates := [][][]float64{}
	for ring := range chDone {
		coordinates = append(coordinates, ring.cx)
	}

	p := struct {
		Polygon
		Type string `json:"type"`
	}{Polygon{poly.CRSReferencable, coordinates}, "Polygon"}

	b, err := json.Marshal(p)
	return b, err
}

func (mpt *MultiPoint) MarshalJSON() ([]byte, error) {
	mp := struct {
		MultiPoint
		Type string `json:"type"`
	}{*mpt, "MultiPoint"}

	b, err := json.Marshal(mp)
	return b, err
}

func (mls *MultiLineString) MarshalJSON() ([]byte, error) {
	l := struct {
		MultiLineString
		Type string `json:"type"`
	}{*mls, "MultiLineString"}

	b, err := json.Marshal(l)
	return b, err
}

func (mpoly *MultiPolygon) MarshalJSON() ([]byte, error) {
	// enforce CCW winding on external rings and CW winding on internal rings

	chWinding := make(chan Ring)
	chClosed := make(chan Ring)
	chDone := make(chan Ring)

	ringSlice := []Ring{}
	for _, polycx := range mpoly.Coordinates {
		for i, slc := range polycx {
			ringSlice = append(ringSlice, Ring{i, slc})
		}
	}

	go ringSender(ringSlice, chWinding)
	go windingEnforcer(chWinding, chClosed)
	go closedEnforcer(chClosed, chDone)

	coordinates := [][][][]float64{}
	for ring := range chDone {
		if ring.i == 0 {
			coordinates = append(coordinates, [][][]float64{ring.cx})
		} else {
			coordinates[len(coordinates)-1] = append(coordinates[len(coordinates)-1], ring.cx)
		}
	}

	p := struct {
		MultiPolygon
		Type string `json:"type"`
	}{MultiPolygon{mpoly.CRSReferencable, coordinates}, "MultiPolygon"}

	b, err := json.Marshal(p)
	return b, err
}

func (gc *GeometryCollection) MarshalJSON() ([]byte, error) {
	collection := struct {
		GeometryCollection
		Type string `json:"type"`
	}{*gc, "GeometryCollection"}

	b, err := json.Marshal(collection)
	return b, err
}

func (f *Feature) MarshalJSON() ([]byte, error) {
	feature := struct {
		Feature
		Type string `json:"type"`
	}{*f, "Feature"}

	b, err := json.Marshal(feature)
	return b, err
}

func (fc *FeatureCollection) MarshalJSON() ([]byte, error) {
	collection := struct {
		FeatureCollection
		Type string `json:"type"`
	}{*fc, "FeatureCollection"}

	b, err := json.Marshal(collection)
	return b, err
}
