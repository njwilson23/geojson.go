/*
 * Package that models GeoJSON types
 */
package geojson

import "math"

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

func (g *CRSReferencable) GetCrs() *CRS {
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
	Type        string   `json:"type"`
	Coordinates Position `json:"coordinates"`
}

func (p *Point) Bbox() *bbox {
	return &bbox{p.Coordinates.X(), p.Coordinates.Y(),
		p.Coordinates.X(), p.Coordinates.Y()}
}

type LineString struct {
	CRSReferencable
	Type        string     `json:"type"`
	Coordinates []Position `json:"coordinates"`
}

func (g *LineString) Bbox() *bbox {
	xmin := g.Coordinates[0].X()
	ymin := g.Coordinates[0].Y()
	xmax := g.Coordinates[0].X()
	ymax := g.Coordinates[0].Y()
	var i int
	for i = 1; i != len(g.Coordinates); i++ {
		xmin = math.Min(g.Coordinates[i].X(), xmin)
		ymin = math.Min(g.Coordinates[i].Y(), ymin)
		xmax = math.Max(g.Coordinates[i].X(), xmax)
		ymax = math.Max(g.Coordinates[i].Y(), ymax)
	}
	return &bbox{xmin, ymin, xmax, ymax}
}

type Polygon struct {
	CRSReferencable
	Type        string       `json:"type"`
	Coordinates [][]Position `json:"coordinates"`
}

func (g *Polygon) Bbox() *bbox {
	xmin := g.Coordinates[0][0].X()
	ymin := g.Coordinates[0][0].Y()
	xmax := g.Coordinates[0][0].X()
	ymax := g.Coordinates[0][0].Y()
	var i int
	for i = 1; i != len(g.Coordinates); i++ {
		xmin = math.Min(g.Coordinates[0][i].X(), xmin)
		ymin = math.Min(g.Coordinates[0][i].Y(), ymin)
		xmax = math.Max(g.Coordinates[0][i].X(), xmax)
		ymax = math.Max(g.Coordinates[0][i].Y(), ymax)
	}
	return &bbox{xmin, ymin, xmax, ymax}
}

type MultiLineString struct {
	CRSReferencable
	Type        string       `json:"type"`
	Coordinates [][]Position `json:"coordinates"`
}

func (g *MultiLineString) Bbox() *bbox {
	xmin := g.Coordinates[0][0].X()
	ymin := g.Coordinates[0][0].Y()
	xmax := g.Coordinates[0][0].X()
	ymax := g.Coordinates[0][0].Y()
	var i, j int
	var position Position
	for i = 0; i != len(g.Coordinates); i++ {
		for j = 0; j != len(g.Coordinates[i]); j++ {
			position = g.Coordinates[i][j]
			xmin = math.Min(position.X(), xmin)
			ymin = math.Min(position.Y(), ymin)
			xmax = math.Max(position.X(), xmax)
			ymax = math.Max(position.Y(), ymax)
		}
	}
	return &bbox{xmin, ymin, xmax, ymax}
}

type MultiPolygon struct {
	CRSReferencable
	Type     string `json:"type"`
	polygons [][][]Position
}

type Feature struct {
	CRSReferencable
	Type       string `json:"type"`
	geometry   Geometry
	properties map[string]string
}

type GeometryCollection struct {
	CRSReferencable
	Type       string `json:"type"`
	geometries []Geometry
}

type FeatureCollection struct {
	CRSReferencable
	Type     string `json:"type"`
	features []Feature
}
