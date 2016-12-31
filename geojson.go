/*
 * Package that models GeoJSON types
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

type Position interface {
	Vertex() []float64
	X() float64
	Y() float64
}

type Position2 struct {
	x, y float64
}

func (p *Position2) Vertex() []float64 {
	vertex := make([]float64, 2)
	vertex[0] = p.x
	vertex[1] = p.y
	return vertex
}

func (p *Position2) X() float64 {
	return p.x
}

func (p *Position2) Y() float64 {
	return p.y
}

type Position3 struct {
	x, y, z float64
}

func (p *Position3) Vertex() []float64 {
	vertex := make([]float64, 3)
	vertex[0] = p.x
	vertex[1] = p.y
	vertex[2] = p.z
	return vertex
}

func (p *Position3) X() float64 {
	return p.x
}

func (p *Position3) Y() float64 {
	return p.y
}

type Point struct {
	CRSReferencable
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

func (p Point) Bbox() *bbox {
	return &bbox{p.Coordinates[0], p.Coordinates[1],
		p.Coordinates[0], p.Coordinates[1]}
}

type LineString struct {
	CRSReferencable
	Type        string      `json:"type"`
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
	Type        string        `json:"type"`
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
	Type        string      `json:"type"`
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
	Type        string        `json:"type"`
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
	Type        string          `json:"type"`
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

type Feature struct {
	CRSReferencable
	Type       string `json:"type"`
	Geometry   Geometry
	Properties map[string]string
}

type GeometryCollection struct {
	CRSReferencable
	Type       string `json:"type"`
	Geometries []Geometry
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

type FeatureCollection struct {
	CRSReferencable
	Type     string `json:"type"`
	Features []Feature
}
