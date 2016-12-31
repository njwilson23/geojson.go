package geojson

import (
	"encoding/json"
)

func (pt Point) MarshalJSON() ([]byte, error) {
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

func (ls LineString) MarshalJSON() ([]byte, error) {
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

func (poly Polygon) MarshalJSON() ([]byte, error) {
	var p struct {
		CRSReferencable
		Type        string        `json:"type"`
		Coordinates [][][]float64 `json:"coordinates"`
	}
	p.Type = "Polygon"
	p.Crs = poly.Crs
	p.Coordinates = poly.Coordinates
	b, err := json.Marshal(p)
	return b, err
}

func (mpt MultiPoint) MarshalJSON() ([]byte, error) {
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

func (mls MultiLineString) MarshalJSON() ([]byte, error) {
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

func (mpoly MultiPolygon) MarshalJSON() ([]byte, error) {
	var p struct {
		CRSReferencable
		Type        string          `json:"type"`
		Coordinates [][][][]float64 `json:"coordinates"`
	}
	p.Type = "MultiPolygon"
	p.Crs = mpoly.Crs
	p.Coordinates = mpoly.Coordinates
	b, err := json.Marshal(p)
	return b, err
}

func (gc GeometryCollection) MarshalJSON() ([]byte, error) {
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

func (f Feature) MarshalJSON() ([]byte, error) {
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

func (fc FeatureCollection) MarshalJSON() ([]byte, error) {
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

func AsGeoJSON(g Geometry) ([]byte, error) {
	b, err := json.Marshal(g)
	if err != nil {
		return []byte{}, err
	}
	return b, err
}

func MarshalFeature(f Feature) ([]byte, error) {
	b, err := json.Marshal(f)
	if err != nil {
		return []byte{}, err
	}
	return b, err
}
