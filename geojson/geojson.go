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

type Bbox struct {
	xmin, ymin, xmax, ymax float64
}

// Geo represents a GeoJSON entity
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

func (g *Geo) Bbox() (bb *Bbox, err error) {
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

type Boundable interface {
	Bbox() (*Bbox, error)
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

func (p *Point) Bbox() (*Bbox, error) {
	return &Bbox{p.Coordinates[0], p.Coordinates[1],
		p.Coordinates[0], p.Coordinates[1]}, nil
}

func (g *LineString) Bbox() (*Bbox, error) {
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
	return &Bbox{xmin, ymin, xmax, ymax}, nil
}

func (g *Polygon) Bbox() (*Bbox, error) {
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
	return &Bbox{xmin, ymin, xmax, ymax}, nil
}

func (g *MultiPoint) Bbox() (*Bbox, error) {
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
	return &Bbox{xmin, ymin, xmax, ymax}, nil
}

func (g *MultiLineString) Bbox() (*Bbox, error) {
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
	return &Bbox{xmin, ymin, xmax, ymax}, nil
}

func (g *MultiPolygon) Bbox() (*Bbox, error) {
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
	return &Bbox{xmin, ymin, xmax, ymax}, nil

}

func (coll *GeometryCollection) Bbox() (bb *Bbox, err error) {
	bboxes := make([]*Bbox, len(coll.Geometries))
	for _, g := range coll.Geometries {
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

func (f *Feature) Bbox() (bb *Bbox, err error) {
	bb, err = f.Geometry.Bbox()
	return
}

func (coll *FeatureCollection) Bbox() (bb *Bbox, err error) {
	bboxes := make([]*Bbox, len(coll.Features))
	for _, f := range coll.Features {
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
	return fmt.Sprintf("LineString[%d]", len(g.Coordinates))
}

func (g *Polygon) String() string {
	return fmt.Sprintf("Polygon[%d]", len(g.Coordinates))
}

func (g *MultiPoint) String() string {
	return fmt.Sprintf("MultiPoint %.6f", g.Coordinates)
}

func (g *MultiLineString) String() string {
	return fmt.Sprintf("MultiLineString[%d]", len(g.Coordinates))
}

func (g *MultiPolygon) String() string {
	return fmt.Sprintf("MultiPolygon[%d]", len(g.Coordinates))
}

func (coll *GeometryCollection) String() string {
	return fmt.Sprintf("GeometryCollection[%d]", len(coll.Geometries))
}

func (f *Feature) String() string {
	return fmt.Sprintf("Feature(%s)", f.Geometry.Type)
}

func (coll *FeatureCollection) String() string {
	return fmt.Sprintf("FeatureCollection[%d]", len(coll.Features))
}
