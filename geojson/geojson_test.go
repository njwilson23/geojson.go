package geojson

import (
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

var WGS84 *CRS = NameCRS("urn:ogc:def:crs:OGC::CRS84")

/* TEST FUNCTIONS */

func TestMarshalPointNoCRS(t *testing.T) {
	point := new(Point)
	point.Coordinates = []float64{3, 4}
	b, err := MarshalGeometry(point)
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
	b, err := MarshalGeometry(point)
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
	b, err := MarshalGeometry(lineString)
	if err != nil {
		fmt.Println("error", err)
		t.Error()
	}
	ref := `{"crs":{"type":"name","properties":{"name":"urn:ogc:def:crs:OGC::CRS84"}},"coordinates":[[2,1],[3,-2],[4,-1]],"type":"LineString"}`
	if strings.Compare(string(b), ref) != 0 {
		t.Fail()
	}
}

func TestMarshalPolygon(t *testing.T) {
	poly := new(Polygon)
	poly.Coordinates = [][][]float64{
		[][]float64{
			[]float64{2.0, 1.0}, []float64{3.0, -2.0}, []float64{4.0, -1.0},
		},
	}
	poly.CRS = WGS84
	b, err := MarshalGeometry(poly)
	if err != nil {
		fmt.Println("error", err)
		t.Error()
	}
	ref := `{"crs":{"type":"name","properties":{"name":"urn:ogc:def:crs:OGC::CRS84"}},"coordinates":[[[4,-1],[3,-2],[2,1],[4,-1]]],"type":"Polygon"}`
	if strings.Compare(string(b), ref) != 0 {
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
	b, err := MarshalGeometry(mpoly)
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

func TestMarshalFeature(t *testing.T) {
	f := new(Feature)
	geom := &Point{CRSReferencable{WGS84}, []float64{3.0, 4.0}}
	prop := make(map[string]interface{})
	prop["a"] = 49
	prop["b"] = 17
	f.CRS = WGS84
	f.Properties = prop
	f.Geometry = geom
	b, err := MarshalFeature(f)
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

// TestUnmarshalInvalid ensures that an error is emitted when the type is not a valid GeoJSON type
func TestUnmarshalInvalid(t *testing.T) {
	_, err := UnmarshalGeoJSON([]byte(`{ "type": "FauxPoint", "coordinates": [100.0, 0.0] }`))
	if err == nil {
		t.Fail()
	}
}

func TestUnmarshalPoint(t *testing.T) {
	contents, err := UnmarshalGeoJSON([]byte(`{ "type": "Point", "coordinates": [100.0, 0.0] }`))
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	if len(contents.Points) != 1 {
		t.Fail()
	}
}

func TestUnmarshalLineString(t *testing.T) {
	contents, err := UnmarshalGeoJSON([]byte(`{ "type": "LineString", "coordinates": [ [100.0, 0.0], [101.0, 1.0] ] }`))
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	if len(contents.LineStrings) != 1 {
		t.Fail()
	}
}

func TestUnmarshalPolygonNoHoles(t *testing.T) {
	contents, err := UnmarshalGeoJSON([]byte(`{ "type": "Polygon", "coordinates": [ [ [100.0, 0.0], [101.0, 0.0], [101.0, 1.0], [100.0, 1.0], [100.0, 0.0] ] ] }`))
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	if len(contents.Polygons) != 1 {
		t.Fail()
	}
	if len(contents.Polygons[0].Coordinates) != 1 {
		t.Fail()
	}
}

func TestUnmarshalPolygonHoles(t *testing.T) {
	contents, err := UnmarshalGeoJSON([]byte(`{ "type": "Polygon", "coordinates": [ [ [100.0, 0.0], [101.0, 0.0], [101.0, 1.0], [100.0, 1.0], [100.0, 0.0] ],
      [ [100.2, 0.2], [100.8, 0.2], [100.8, 0.8], [100.2, 0.8], [100.2, 0.2] ] ] }`))
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	if len(contents.Polygons) != 1 {
		t.Fail()
	}
	if len(contents.Polygons[0].Coordinates) != 2 {
		t.Fail()
	}
}

func TestUnmarshalMultiPoint(t *testing.T) {
	contents, err := UnmarshalGeoJSON([]byte(`{"type": "MultiPoint", "coordinates": [ [100.0, 0.0], [101.0, 1.0] ] }`))
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	if len(contents.MultiPoints) != 1 {
		t.Fail()
	}
}

func TestUnmarshalMultiLineString(t *testing.T) {
	contents, err := UnmarshalGeoJSON([]byte(`{"type": "MultiLineString",
    "coordinates": [
        [ [100.0, 0.0], [101.0, 1.0] ],
        [ [102.0, 2.0], [103.0, 3.0] ]
      ]
    }`))
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	if len(contents.MultiLineStrings) != 1 {
		t.Fail()
	}
}

func TestUnmarshalMultiPolygon(t *testing.T) {
	contents, err := UnmarshalGeoJSON([]byte(`{"type": "MultiPolygon",
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
	if len(contents.MultiPolygons) != 1 {
		t.Fail()
	}
}

func TestUnmarshalGeometryCollection(t *testing.T) {
	content, err := UnmarshalGeoJSON([]byte(`{ "type": "GeometryCollection",
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
	if len(content.Points) != 1 {
		t.Fail()
	}
	if len(content.LineStrings) != 1 {
		t.Fail()
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
}

// func TestValidateClosedRing(t *testing.T) {
// 	ring := [][]float64{[]float64{0, 0}, []float64{1, 0}, []float64{1, 1}, []float64{0, 1}}
// 	corrections := make(chan bool)
// 	validateRing(0, ring, corrections)
// 	result := <-corrections
// 	fmt.Println(result)
// 	if result == false {
// 		t.Fail()
// 	}
// 	if len(ring) != 5 {
// 		t.Fail()
// 	}
// }
