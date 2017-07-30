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
	var p struct {
		CRSReferencable
		Type        string    `json:"type"`
		Coordinates []float64 `json:"coordinates"`
	}
	p.Type = "Point"
	p.Crs = pt.Crs
	p.Coordinates = pt.Coordinates
	b, err := json.Marshal(p)
	return b, err
}

func (ls *LineString) MarshalJSON() ([]byte, error) {
	var l struct {
		CRSReferencable
		Type        string      `json:"type"`
		Coordinates [][]float64 `json:"coordinates"`
	}
	l.Type = "LineString"
	l.Crs = ls.Crs
	l.Coordinates = ls.Coordinates
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

	var p struct {
		CRSReferencable
		Type        string        `json:"type"`
		Coordinates [][][]float64 `json:"coordinates"`
	}
	p.Type = "Polygon"
	p.Crs = poly.Crs
	p.Coordinates = coordinates
	b, err := json.Marshal(p)
	return b, err
}

func (mpt *MultiPoint) MarshalJSON() ([]byte, error) {
	var p struct {
		CRSReferencable
		Type        string      `json:"type"`
		Coordinates [][]float64 `json:"coordinates"`
	}
	p.Type = "MultiPoint"
	p.Crs = mpt.Crs
	p.Coordinates = mpt.Coordinates
	b, err := json.Marshal(p)
	return b, err
}

func (mls *MultiLineString) MarshalJSON() ([]byte, error) {
	var l struct {
		CRSReferencable
		Type        string        `json:"type"`
		Coordinates [][][]float64 `json:"coordinates"`
	}
	l.Type = "MultiMineString"
	l.Crs = mls.Crs
	l.Coordinates = mls.Coordinates
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

	var p struct {
		CRSReferencable
		Type        string          `json:"type"`
		Coordinates [][][][]float64 `json:"coordinates"`
	}
	p.Type = "MultiPolygon"
	p.Crs = mpoly.Crs
	p.Coordinates = coordinates
	b, err := json.Marshal(p)
	return b, err
}

func (gc *GeometryCollection) MarshalJSON() ([]byte, error) {
	var collection struct {
		CRSReferencable
		Type       string     `json:"type"`
		Geometries []Geometry `json:"geometry"`
	}
	collection.Type = "GeometryCollection"
	collection.Crs = gc.Crs
	collection.Geometries = gc.Geometries
	b, err := json.Marshal(collection)
	return b, err
}

func (f *Feature) MarshalJSON() ([]byte, error) {
	var feature struct {
		CRSReferencable
		Type       string      `json:"type"`
		Id         string      `json:"id,omitempty"`
		Geometry   Geometry    `json:"geometry"`
		Properties interface{} `json:"properties"`
	}
	feature.Type = "Feature"
	feature.Id = f.Id
	feature.Crs = f.Crs
	feature.Geometry = f.Geometry
	feature.Properties = f.Properties
	b, err := json.Marshal(feature)
	return b, err
}

func (fc *FeatureCollection) MarshalJSON() ([]byte, error) {
	var collection struct {
		CRSReferencable
		Type     string    `json:"type"`
		Features []Feature `json:"geometry"`
	}
	collection.Type = "FeatureCollection"
	collection.Crs = fc.Crs
	collection.Features = fc.Features
	b, err := json.Marshal(collection)
	return b, err
}
