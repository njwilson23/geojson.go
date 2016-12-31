/*
 * Package that models GeoJSON types and makes marshaling and unmarshaling between
 * Go structs and GeoJSON bytes simpler
 */
package geojson

import (
	"fmt"
	"math"
)

type bbox struct {
	xmin, ymin, xmax, ymax float64
}

type CRS struct {
	Type       string            `json:"type"`
	Properties map[string]string `json:"properties"`
}

func NameCRS(name string) *CRS {
	prop := make(map[string]string)
	prop["name"] = name
	return &CRS{"name", prop}
}

func LinkCRS(link string) *CRS {
	prop := make(map[string]string)
	prop["link"] = link
	return &CRS{"link", prop}
}

type CRSReferencable struct {
	Crs CRS `json:"crs,omitempty"`
}

func (g CRSReferencable) GetCrs() *CRS {
	return &g.Crs
}

type Geometry interface {
	Bbox() *bbox
	GetCrs() *CRS
}

type Point struct {
	CRSReferencable
	Coordinates []float64 `json:"coordinates"`
}

func (p Point) Bbox() *bbox {
	return &bbox{p.Coordinates[0], p.Coordinates[1],
		p.Coordinates[0], p.Coordinates[1]}
}

type LineString struct {
	CRSReferencable
	Coordinates [][]float64 `json:"coordinates"`
}

func (g LineString) Bbox() *bbox {
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
	return &bbox{xmin, ymin, xmax, ymax}
}

func (g Point) String() string {
	return fmt.Sprintf("Point %.2f", g.Coordinates)
}

func (g LineString) String() string {
	if len(g.Coordinates) <= 8 {
		return fmt.Sprintf("LineString %.2f", g.Coordinates)
	} else {
		return fmt.Sprintf("LineString %.2f...", g.Coordinates[0:8])
	}
}

func (g Polygon) String() string {
	if len(g.Coordinates[0]) <= 8 {
		return fmt.Sprintf("Polygon %.2f", g.Coordinates[0])
	} else {
		return fmt.Sprintf("Polygon %.2f...", g.Coordinates[0][0:8])
	}
}

type Polygon struct {
	CRSReferencable
	Coordinates [][][]float64 `json:"coordinates"`
}

func (g Polygon) Bbox() *bbox {
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
	return &bbox{xmin, ymin, xmax, ymax}
}

type MultiPoint struct {
	CRSReferencable
	Coordinates [][]float64 `json:"coordinates"`
}

func (g MultiPoint) Bbox() *bbox {
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
	return &bbox{xmin, ymin, xmax, ymax}
}

type MultiLineString struct {
	CRSReferencable
	Coordinates [][][]float64 `json:"coordinates"`
}

func (g MultiLineString) Bbox() *bbox {
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
	return &bbox{xmin, ymin, xmax, ymax}
}

type MultiPolygon struct {
	CRSReferencable
	Coordinates [][][][]float64 `json:"coordinates"`
}

func (g MultiPolygon) Bbox() *bbox {
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
	return &bbox{xmin, ymin, xmax, ymax}

}

type GeometryCollection struct {
	CRSReferencable
	Geometries []Geometry `json:"geometries"`
}

func (collection GeometryCollection) Bbox() *bbox {
	bb := collection.Geometries[0].Bbox()
	xmin := bb.xmin
	ymin := bb.ymin
	xmax := bb.xmax
	ymax := bb.ymax
	if len(collection.Geometries) > 1 {
		for i := 1; i != len(collection.Geometries); i++ {
			bb = collection.Geometries[i].Bbox()
			xmin = math.Min(xmin, bb.xmin)
			ymin = math.Min(ymin, bb.ymin)
			xmax = math.Max(xmax, bb.xmin)
			ymax = math.Max(ymax, bb.ymax)
		}
	}
	return &bbox{xmin, ymin, xmax, ymax}
}

type Feature struct {
	CRSReferencable
	Id         string      `json:"id,omitempty"`
	Geometry   Geometry    `json:"geometry"`
	Properties interface{} `json:"properties"`
}

type FeatureCollection struct {
	CRSReferencable
	Features []Feature `json:"features"`
}
