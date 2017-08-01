/*
 * Package that models GeoJSON types and makes marshaling and unmarshaling between
 * Go structs and GeoJSON bytes simpler
 */
package geojson

import (
	"fmt"
	"math"
)

type CRS struct {
	Type       string            `json:"type"`
	Properties map[string]string `json:"properties"`
}

type CRSReferencable struct {
	CRS *CRS `json:"crs,omitempty"`
}

func (g *CRSReferencable) GetCRS() *CRS {
	return g.CRS
}

type bbox struct {
	xmin, ymin, xmax, ymax float64
}

type Geo struct {
	Type               string
	Point              *Point
	LineString         *LineString
	Polygon            *Polygon
	MultiPoint         *MultiPoint
	MultiLineString    *MultiLineString
	MultiPolygon       *MultiPolygon
	Feature            *Feature
	GeometryCollection *GeometryCollection
	FeatureCollection  *FeatureCollection
}

func (g *Geo) Bbox() (bb *bbox, err error) {
	switch g.Type {
	case "Point":
		bb, err = g.Point.Bbox()
	case "LineString":
		bb, err = g.LineString.Bbox()
	case "Polygon":
		bb, err = g.Polygon.Bbox()
	case "MultiPoint":
		bb, err = g.MultiPoint.Bbox()
	case "MultiLineString":
		bb, err = g.MultiLineString.Bbox()
	case "MultiPolygon":
		bb, err = g.MultiPolygon.Bbox()
	case "GeometryCollection":
		bb, err = g.GeometryCollection.Bbox()
	case "Feature":
		bb, err = g.Feature.Bbox()
	case "FeatureCollection":
		bb, err = g.FeatureCollection.Bbox()
	default:
		err = fmt.Errorf("unhandled type: '%s'", g.Type)
	}
	return
}

type Geometry interface {
	Bbox() (*bbox, error)
	GetCRS() *CRS
}

type Point struct {
	CRSReferencable
	Coordinates []float64 `json:"coordinates"`
}

type LineString struct {
	CRSReferencable
	Coordinates [][]float64 `json:"coordinates"`
}

type Polygon struct {
	CRSReferencable
	Coordinates [][][]float64 `json:"coordinates"`
}

type MultiPoint struct {
	CRSReferencable
	Coordinates [][]float64 `json:"coordinates"`
}

type MultiLineString struct {
	CRSReferencable
	Coordinates [][][]float64 `json:"coordinates"`
}

type MultiPolygon struct {
	CRSReferencable
	Coordinates [][][][]float64 `json:"coordinates"`
}

type GeometryCollection struct {
	CRSReferencable
	Geometries []*Geo `json:"geometries"`
}

type Feature struct {
	CRSReferencable
	ID         string                 `json:"id,omitempty"`
	Geometry   Geo                    `json:"geometry"`
	Properties map[string]interface{} `json:"properties"`
}

type FeatureCollection struct {
	CRSReferencable
	Features []Feature `json:"features"`
}

/* Bbox methods */

func (p *Point) Bbox() (*bbox, error) {
	return &bbox{p.Coordinates[0], p.Coordinates[1],
		p.Coordinates[0], p.Coordinates[1]}, nil
}

func (g *LineString) Bbox() (*bbox, error) {
	xmin := g.Coordinates[0][0]
	ymin := g.Coordinates[0][1]
	xmax := g.Coordinates[0][0]
	ymax := g.Coordinates[0][1]
	var i int
	for i = 1; i != len(g.Coordinates); i++ {
		xmin = math.Min(g.Coordinates[i][0], xmin)
		ymin = math.Min(g.Coordinates[i][1], ymin)
		xmax = math.Max(g.Coordinates[i][0], xmax)
		ymax = math.Max(g.Coordinates[i][1], ymax)
	}
	return &bbox{xmin, ymin, xmax, ymax}, nil
}

func (g *Polygon) Bbox() (*bbox, error) {
	xmin := g.Coordinates[0][0][0]
	ymin := g.Coordinates[0][0][1]
	xmax := g.Coordinates[0][0][0]
	ymax := g.Coordinates[0][0][1]
	var i int
	for i = 1; i != len(g.Coordinates); i++ {
		xmin = math.Min(g.Coordinates[0][i][0], xmin)
		ymin = math.Min(g.Coordinates[0][i][1], ymin)
		xmax = math.Max(g.Coordinates[0][i][0], xmax)
		ymax = math.Max(g.Coordinates[0][i][1], ymax)
	}
	return &bbox{xmin, ymin, xmax, ymax}, nil
}

func (g *MultiPoint) Bbox() (*bbox, error) {
	xmin := g.Coordinates[0][0]
	ymin := g.Coordinates[0][1]
	xmax := g.Coordinates[0][0]
	ymax := g.Coordinates[0][1]
	var i int
	var position []float64
	for i = 0; i != len(g.Coordinates); i++ {
		position = g.Coordinates[i]
		xmin = math.Min(position[0], xmin)
		ymin = math.Min(position[1], ymin)
		xmax = math.Max(position[0], xmax)
		ymax = math.Max(position[1], ymax)
	}
	return &bbox{xmin, ymin, xmax, ymax}, nil
}

func (g *MultiLineString) Bbox() (*bbox, error) {
	xmin := g.Coordinates[0][0][0]
	ymin := g.Coordinates[0][0][1]
	xmax := g.Coordinates[0][0][0]
	ymax := g.Coordinates[0][0][1]
	var i, j int
	var position []float64
	for i = 0; i != len(g.Coordinates); i++ {
		for j = 0; j != len(g.Coordinates[i]); j++ {
			position = g.Coordinates[i][j]
			xmin = math.Min(position[0], xmin)
			ymin = math.Min(position[1], ymin)
			xmax = math.Max(position[0], xmax)
			ymax = math.Max(position[1], ymax)
		}
	}
	return &bbox{xmin, ymin, xmax, ymax}, nil
}

func (g *MultiPolygon) Bbox() (*bbox, error) {
	xmin := g.Coordinates[0][0][0][0]
	ymin := g.Coordinates[0][0][0][1]
	xmax := g.Coordinates[0][0][0][0]
	ymax := g.Coordinates[0][0][0][1]
	var i, j int
	var position []float64
	for i = 0; i != len(g.Coordinates); i++ {
		for j = 0; j != len(g.Coordinates[i]); j++ {
			position = g.Coordinates[i][0][j]
			xmin = math.Min(position[0], xmin)
			ymin = math.Min(position[1], ymin)
			xmax = math.Max(position[0], xmax)
			ymax = math.Max(position[1], ymax)
		}
	}
	return &bbox{xmin, ymin, xmax, ymax}, nil

}

func (collection *GeometryCollection) Bbox() (bb *bbox, err error) {
	bboxes := make([]*bbox, len(collection.Geometries))
	for _, g := range collection.Geometries {
		bb, err := g.Bbox()
		if err != nil {
			break
		}
		bboxes = append(bboxes, bb)
	}
	if err == nil {
		bb, err = unionBbox(bboxes)
	}
	return
}

func (f *Feature) Bbox() (bb *bbox, err error) {
	bb, err = f.Geometry.Bbox()
	return
}

func (fc *FeatureCollection) Bbox() (bb *bbox, err error) {
	bboxes := make([]*bbox, len(fc.Features))
	for _, f := range fc.Features {
		bb, err := f.Bbox()
		if err != nil {
			break
		}
		bboxes = append(bboxes, bb)
	}
	if err == nil {
		bb, err = unionBbox(bboxes)
	}
	return
}

/* String methods */

func (g *Point) String() string {
	return fmt.Sprintf("Point %.6f", g.Coordinates)
}

func (g *LineString) String() string {
	if len(g.Coordinates) <= 8 {
		return fmt.Sprintf("LineString %.6f", g.Coordinates)
	} else {
		return fmt.Sprintf("LineString %.6f...", g.Coordinates[0:8])
	}
}

func (g *Polygon) String() string {
	if len(g.Coordinates[0]) <= 8 {
		return fmt.Sprintf("Polygon %.6f", g.Coordinates[0])
	} else {
		return fmt.Sprintf("Polygon %.6f...", g.Coordinates[0][0:8])
	}
}

func (g *GeometryCollection) String() string {
	return fmt.Sprintf("GeometryCollection (%d members)", len(g.Geometries))
}
