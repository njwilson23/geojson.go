package geojson

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

// NameCRS returns a named CRS object
// Note that for RFC 7946 compliance, WGS84 may be used
func NameCRS(name string) *CRS {
	prop := make(map[string]string)
	prop["name"] = name
	return &CRS{"name", prop}
}

var WGS84 = NameCRS("urn:ogc:def:crs:OGC::CRS84")

func TestMarshalPointNoCRS(t *testing.T) {
	point := new(Point)
	point.Coordinates = []float64{3, 4}
	b, err := json.Marshal(point)
	if err != nil {
		fmt.Println("error", err)
		t.Fail()
	}
	ref := `{"coordinates":[3,4],"type":"Point"}`
	if strings.Compare(string(b), ref) != 0 {
		fmt.Println("recieved    ", string(b))
		fmt.Println("but expected", ref)
		t.Fail()
	}
}

func TestMarshalPoint(t *testing.T) {
	point := &Point{CRSReferencable{WGS84}, []float64{3.0, 4.0}}
	b, err := json.Marshal(point)
	if err != nil {
		fmt.Println("error", err)
		t.Fail()
	}
	ref := `{"crs":{"type":"name","properties":{"name":"urn:ogc:def:crs:OGC::CRS84"}},"coordinates":[3,4],"type":"Point"}`
	if strings.Compare(string(b), ref) != 0 {
		fmt.Println("recieved    ", string(b))
		fmt.Println("but expected", ref)
		t.Fail()
	}
}

func TestMarshalLineString(t *testing.T) {
	lineString := new(LineString)
	lineString.Coordinates = [][]float64{
		[]float64{2.0, 1.0}, []float64{3.0, -2.0}, []float64{4.0, -1.0},
	}
	lineString.CRS = WGS84
	b, err := json.Marshal(lineString)
	if err != nil {
		fmt.Println("error", err)
		t.Error()
	}
	ref := `{"crs":{"type":"name","properties":{"name":"urn:ogc:def:crs:OGC::CRS84"}},"coordinates":[[2,1],[3,-2],[4,-1]],"type":"LineString"}`
	if strings.Compare(string(b), ref) != 0 {
		t.Fail()
	}
}

func TestMarshalIncompleteBackwardsPolygon(t *testing.T) {
	poly := new(Polygon)
	poly.Coordinates = [][][]float64{
		[][]float64{
			[]float64{0.0, 1.0}, []float64{2.0, 0.0}, []float64{0.0, -1.0},
		},
	}
	poly.CRS = WGS84
	b, err := json.Marshal(poly)
	if err != nil {
		fmt.Println("error", err)
		t.Error()
	}
	ref := `{"crs":{"type":"name","properties":{"name":"urn:ogc:def:crs:OGC::CRS84"}},"coordinates":[[[0,1],[0,-1],[2,0],[0,1]]],"type":"Polygon"}`
	if strings.Compare(string(b), ref) != 0 {
		fmt.Println("recieved    ", string(b))
		fmt.Println("but expected", ref)
		t.Fail()
	}
}

func TestMarshalPolygon(t *testing.T) {
	poly := new(Polygon)
	poly.Coordinates = [][][]float64{
		[][]float64{
			[]float64{2.0, 1.0}, []float64{3.0, -2.0}, []float64{4.0, -1.0}, []float64{2.0, 1.0},
		},
	}
	poly.CRS = WGS84
	b, err := json.Marshal(poly)
	if err != nil {
		fmt.Println("error", err)
		t.Error()
	}
	ref := `{"crs":{"type":"name","properties":{"name":"urn:ogc:def:crs:OGC::CRS84"}},"coordinates":[[[2,1],[3,-2],[4,-1],[2,1]]],"type":"Polygon"}`
	if strings.Compare(string(b), ref) != 0 {
		fmt.Println("recieved    ", string(b))
		fmt.Println("but expected", ref)
		t.Fail()
	}
}

func TestMarshalMultiPolygon(t *testing.T) {
	// creates a two-part multipolygon, with a hole in the second part
	mpoly := new(MultiPolygon)
	mpoly.CRS = WGS84
	mpoly.Coordinates = [][][][]float64{
		[][][]float64{
			[][]float64{
				[]float64{102, 2}, []float64{103, 2}, []float64{103, 3}, []float64{102, 3}, []float64{102, 2},
			},
		},
		[][][]float64{
			[][]float64{
				[]float64{100, 0}, []float64{101, 0}, []float64{101, 1}, []float64{100, 1}, []float64{100, 0},
			},
			[][]float64{
				[]float64{100.2, 0.2}, []float64{100.8, 0.2}, []float64{100.8, 0.8}, []float64{100.2, 0.8}, []float64{100.2, 0.2},
			},
		},
	}
	b, err := json.Marshal(mpoly)
	ref := `{"crs":{"type":"name","properties":{"name":"urn:ogc:def:crs:OGC::CRS84"}},"coordinates":[[[[102,2],[103,2],[103,3],[102,3],[102,2]]],[[[100,0],[101,0],[101,1],[100,1],[100,0]],[[100.2,0.2],[100.2,0.8],[100.8,0.8],[100.8,0.2],[100.2,0.2]]]],"type":"MultiPolygon"}`
	if err != nil {
		fmt.Println("error", err)
		t.Error()
	}
	if strings.Compare(string(b), ref) != 0 {
		fmt.Println("recieved    ", string(b))
		fmt.Println("but expected", ref)
		t.Fail()
	}
}

func TestMarshalGeo(t *testing.T) {
	geo := Geo{Type: "Point", Point: &Point{CRSReferencable: CRSReferencable{},
		Coordinates: []float64{3, 4}}}

	b, err := json.Marshal(geo)
	if err != nil {
		fmt.Println(err)
		t.Error()
	}
	ref := `{"coordinates":[3,4],"type":"Point"}`
	if strings.Compare(string(b), ref) != 0 {
		fmt.Println("recieved    ", string(b))
		fmt.Println("but expected", ref)
		t.Fail()
	}
}

func TestMarshalFeature(t *testing.T) {
	prop := make(map[string]interface{})
	prop["a"] = 49
	prop["b"] = 17

	f := &Feature{CRSReferencable: CRSReferencable{WGS84},
		Geometry:   Geo{Type: "Point", Point: &Point{CRSReferencable{WGS84}, []float64{3.0, 4.0}}},
		Properties: prop}

	b, err := json.Marshal(f)
	if err != nil {
		fmt.Println("error", err)
		t.Error()
	}
	ref := `{"crs":{"type":"name","properties":{"name":"urn:ogc:def:crs:OGC::CRS84"}},"geometry":{"crs":{"type":"name","properties":{"name":"urn:ogc:def:crs:OGC::CRS84"}},"coordinates":[3,4],"type":"Point"},"properties":{"a":49,"b":17},"type":"Feature"}`
	if strings.Compare(string(b), ref) != 0 {
		fmt.Println("recieved    ", string(b))
		fmt.Println("but expected", ref)
		t.Fail()
	}
}
