package geojson

import (
	"fmt"
	"strings"
	"testing"
)

/* CONVENIENCE FUNCTIONS FOR TESTING */

// NameCRS returns a named CRS object
// Note that for RFC 7946 compliance, WGS84 may be used
func NameCRS(name string) *CRS {
	prop := make(map[string]string)
	prop["name"] = name
	return &CRS{"name", prop}
}

var WGS84 *CRS = NameCRS("urn:ogc:def:crs:OGC::CRS84")

// NewPoint creates a point with the provided coordinates
func NewPoint(x ...float64) *Point {
	g := new(Point)
	g.Coordinates = x
	g.Crs = WGS84
	return g
}

func NewLineString(x ...[]float64) *LineString {
	var ivert int
	var nVertices int
	var pos []float64
	var coordinates [][]float64

	if len(x) == 2 {

		nVertices = len(x[0])
		coordinates = make([][]float64, nVertices)
		for ivert = 0; ivert != nVertices; ivert++ {
			pos = []float64{x[0][ivert], x[1][ivert]}
			coordinates[ivert] = pos
		}

	} else if len(x) == 3 {

		nVertices = len(x[0])
		coordinates = make([][]float64, nVertices)
		for ivert = 0; ivert != nVertices; ivert++ {
			pos = []float64{x[0][ivert], x[1][ivert], x[2][ivert]}
			coordinates[ivert] = pos
		}

	} else {
		panic("NewLineString takes either 2 or 3 arguments of type []float64")
	}

	g := new(LineString)
	g.Coordinates = coordinates
	g.Crs = WGS84
	return g
}

// NewPolygon2 is a convenience constructor for a 2D Polygon. It is called as
// NewPolygon2(x, y, [x_sub1, y_sub1, [x_sub2, y_sub2]]...) where areguments
// are slices of floats.
func NewPolygon2(x ...[]float64) *Polygon {
	var ip, ivert int
	var nParts, nVertices int
	var pos []float64
	var coordinates [][][]float64

	if (len(x) % 2) == 0 {

		nParts = len(x) / 2
		coordinates = make([][][]float64, nParts)
		for ip = 0; ip != nParts; ip++ {
			nVertices = len(x[ip*2])
			coordinates[ip] = make([][]float64, nVertices)
			for ivert = 0; ivert != nVertices; ivert++ {
				pos = []float64{x[ip*2][ivert], x[ip*2+1][ivert]}
				coordinates[ip][ivert] = pos
			}
		}

	} else {
		panic("NewPolygon2 called with odd number of arguments")
	}

	g := new(Polygon)
	g.Coordinates = coordinates
	g.Crs = WGS84
	return g
}

/* TEST FUNCTIONS */

func TestMarshalPointNoCrs(t *testing.T) {
	point := new(Point)
	point.Coordinates = []float64{3, 4}
	b, err := MarshalGeometry(point)
	if err != nil {
		fmt.Println("error", err)
		t.Fail()
	}
	ref := `{"type":"Point","coordinates":[3,4]}`
	if strings.Compare(string(b), ref) != 0 {
		fmt.Println("recieved    ", string(b))
		fmt.Println("but expected", ref)
		t.Fail()
	}
}

func TestMarshalPoint(t *testing.T) {
	point := NewPoint(3.0, 4.0)
	b, err := MarshalGeometry(point)
	if err != nil {
		fmt.Println("error", err)
		t.Fail()
	}
	ref := `{"crs":{"type":"name","properties":{"name":"urn:ogc:def:crs:OGC::CRS84"}},"type":"Point","coordinates":[3,4]}`
	if strings.Compare(string(b), ref) != 0 {
		fmt.Println("recieved    ", string(b))
		fmt.Println("but expected", ref)
		t.Fail()
	}
}

func TestMarshalLineString(t *testing.T) {
	X := []float64{2.0, 3.0, 4.0}
	Y := []float64{1.0, -2.0, -1.0}
	point := NewLineString(X, Y)
	b, err := MarshalGeometry(point)
	if err != nil {
		fmt.Println("error", err)
		t.Error()
	}
	ref := `{"crs":{"type":"name","properties":{"name":"urn:ogc:def:crs:OGC::CRS84"}},"type":"LineString","coordinates":[[2,1],[3,-2],[4,-1]]}`
	if strings.Compare(string(b), ref) != 0 {
		t.Fail()
	}
}

func TestMarshalPolygon(t *testing.T) {
	X := []float64{2.0, 3.0, 4.0}
	Y := []float64{1.0, -2.0, -1.0}
	point := NewPolygon2(X, Y)
	b, err := MarshalGeometry(point)
	if err != nil {
		fmt.Println("error", err)
		t.Error()
	}
	ref := `{"crs":{"type":"name","properties":{"name":"urn:ogc:def:crs:OGC::CRS84"}},"type":"Polygon","coordinates":[[[4,-1],[3,-2],[2,1]]]}`
	if strings.Compare(string(b), ref) != 0 {
		t.Fail()
	}
}

func TestMarshallMultiPolygon(t *testing.T) {
	// creates a two-part multipolygon, with a hole in the second part
	mpoly := new(MultiPolygon)
	mpoly.Crs = WGS84
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
	b, err := MarshalGeometry(mpoly)
	ref := `{"crs":{"type":"name","properties":{"name":"urn:ogc:def:crs:OGC::CRS84"}},"type":"MultiPolygon","coordinates":[[[[102,2],[103,2],[103,3],[102,3],[102,2]]],[[[102,2],[103,2],[103,3],[102,3],[102,2]],[[100,0],[101,0],[101,1],[100,1],[100,0]],[[100.2,0.2],[100.2,0.8],[100.8,0.8],[100.8,0.2],[100.2,0.2]]]]}`
	if err != nil {
		fmt.Println("error", err)
		t.Error()
	}
	if strings.Compare(string(b), ref) != 0 {
		t.Fail()
	}
}

func TestMarshalFeature(t *testing.T) {
	f := new(Feature)
	geom := NewPoint(3.0, 4.0)
	prop := make(map[string]interface{})
	prop["a"] = 49
	prop["b"] = 17
	f.Crs = WGS84
	f.Properties = prop
	f.Geometry = geom
	b, err := MarshalFeature(*f)
	if err != nil {
		fmt.Println("error", err)
		t.Error()
	}
	ref := `{"crs":{"type":"name","properties":{"name":"urn:ogc:def:crs:OGC::CRS84"}},"type":"Feature","geometry":{"crs":{"type":"name","properties":{"name":"urn:ogc:def:crs:OGC::CRS84"}},"type":"Point","coordinates":[3,4]},"properties":{"a":49,"b":17}}`
	if strings.Compare(string(b), ref) != 0 {
		t.Fail()
	}
}

// TestUnmarshalInvalid ensures that an error is emitted when the type is not a valid GeoJSON type
func TestUnmarshalInvalid(t *testing.T) {
	_, err := UnmarshalGeoJSON([]byte(`{ "type": "FauxPoint", "coordinates": [100.0, 0.0] }`))
	if err == nil {
		t.Fail()
	}
}

func TestUnmarshalPoint(t *testing.T) {
	_, err := UnmarshalGeoJSON([]byte(`{ "type": "Point", "coordinates": [100.0, 0.0] }`))
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
}

func TestUnmarshalLineString(t *testing.T) {
	_, err := UnmarshalGeoJSON([]byte(`{ "type": "LineString", "coordinates": [ [100.0, 0.0], [101.0, 1.0] ] }`))
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
}

func TestUnmarshalPolygonNoHoles(t *testing.T) {
	_, err := UnmarshalGeoJSON([]byte(`{ "type": "Polygon", "coordinates": [ [ [100.0, 0.0], [101.0, 0.0], [101.0, 1.0], [100.0, 1.0], [100.0, 0.0] ] ] }`))
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
}

func TestUnmarshalPolygonHoles(t *testing.T) {
	_, err := UnmarshalGeoJSON([]byte(`{ "type": "Polygon", "coordinates": [ [ [100.0, 0.0], [101.0, 0.0], [101.0, 1.0], [100.0, 1.0], [100.0, 0.0] ],
      [ [100.2, 0.2], [100.8, 0.2], [100.8, 0.8], [100.2, 0.8], [100.2, 0.2] ] ] }`))
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
}

func TestUnmarshalMultiPoint(t *testing.T) {
	_, err := UnmarshalGeoJSON([]byte(`{"type": "MultiPoint", "coordinates": [ [100.0, 0.0], [101.0, 1.0] ] }`))
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
}

func TestUnmarshalMultiLineString(t *testing.T) {
	_, err := UnmarshalGeoJSON([]byte(`{"type": "MultiLineString",
    "coordinates": [
        [ [100.0, 0.0], [101.0, 1.0] ],
        [ [102.0, 2.0], [103.0, 3.0] ]
      ]
    }`))
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
}

func TestUnmarshalMultiPolygon(t *testing.T) {
	_, err := UnmarshalGeoJSON([]byte(`{"type": "MultiPolygon",
    "coordinates": [
      [[[102.0, 2.0], [103.0, 2.0], [103.0, 3.0], [102.0, 3.0], [102.0, 2.0]]],
      [[[100.0, 0.0], [101.0, 0.0], [101.0, 1.0], [100.0, 1.0], [100.0, 0.0]],
       [[100.2, 0.2], [100.8, 0.2], [100.8, 0.8], [100.2, 0.8], [100.2, 0.2]]]
      ]
    }`))
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
}

func TestUnmarshalGeometryCollection(t *testing.T) {
	_, err := UnmarshalGeoJSON([]byte(`{ "type": "GeometryCollection",
    "geometries": [
      { "type": "Point",
        "coordinates": [100.0, 0.0]
        },
      { "type": "LineString",
        "coordinates": [ [101.0, 0.0], [102.0, 1.0] ]
        }
    ]
  }`))
	if err != nil {
		fmt.Println(err)
		t.Error()
	}
}

func TestUnmarshalFeature(t *testing.T) {
	contents, err := UnmarshalGeoJSON([]byte(`{ "type": "Feature",
        "geometry": {"type": "Point", "coordinates": [102.0, 0.5]},
        "properties": {"prop0": "value0"}
        }`))
	if err != nil {
		fmt.Println("error:", err)
		t.Error()
	}
	if len(contents.Features) != 1 {
		t.Fail()
	}
	_, ok := contents.Features[0].Geometry.(*Point)
	if !ok {
		t.Fail()
	}
	_, ok = contents.Features[0].Geometry.(*Polygon)
	if ok {
		t.Fail()
	}
	if contents.Features[0].Properties["prop0"] != "value0" {
		t.Fail()
	}
}

func TestUnmarshalFeatureCollection(t *testing.T) {
	contents, err := UnmarshalGeoJSON([]byte(`{ "type": "FeatureCollection",
    "features": [
      { "type": "Feature",
        "geometry": {"type": "Point", "coordinates": [102.0, 0.5]},
        "properties": {"prop0": "value0"}
        },
      { "type": "Feature",
        "geometry": {
          "type": "LineString",
          "coordinates": [
            [102.0, 0.0], [103.0, 1.0], [104.0, 0.0], [105.0, 1.0]
            ]
          },
        "properties": {
          "prop0": "value0",
          "prop1": 0.0
          }
        },
      { "type": "Feature",
         "geometry": {
           "type": "Polygon",
           "coordinates": [
             [ [100.0, 0.0], [101.0, 0.0], [101.0, 1.0],
               [100.0, 1.0], [100.0, 0.0] ]
             ]
         },
         "properties": {
           "prop0": "value0",
           "prop1": {"this": "that"}
           }
         }
       ]
     }`))
	if err != nil {
		fmt.Println("error:", err)
		t.Error()
	}
	if len(contents.Features) != 3 {
		t.Fail()
	}
	if contents.Features[0].Properties["prop0"] != "value0" {
		t.Fail()
	}
	if contents.Features[1].Properties["prop1"] != 0.0 {
		t.Fail()
	}
	_, ok := contents.Features[2].Properties["prop1"].(map[string]interface{})
	if !ok {
		t.Fail()
	}
	_, ok = contents.Features[0].Geometry.(*Point)
	if !ok {
		t.Fail()
	}
	_, ok = contents.Features[1].Geometry.(*LineString)
	if !ok {
		t.Fail()
	}
	_, ok = contents.Features[2].Geometry.(*Polygon)
	if !ok {
		t.Fail()
	}
}
