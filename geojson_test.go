package geojson

import (
	"fmt"
	"strings"
	"testing"
)

func TestMarshalPoint(t *testing.T) {
	point := NewPoint(3.0, 4.0)
	b, err := MarshalGeometry(point)
	if err != nil {
		fmt.Println("error", err)
		t.Fail()
	}
	ref := "{\"crs\":{\"type\":\"name\",\"properties\":{\"name\":\"urn:ogc:def:crs:OGC::CRS84\"}},\"type\":\"Point\",\"coordinates\":[3,4]}"
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
	ref := "{\"crs\":{\"type\":\"name\",\"properties\":{\"name\":\"urn:ogc:def:crs:OGC::CRS84\"}},\"type\":\"LineString\",\"coordinates\":[[2,1],[3,-2],[4,-1]]}"
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
	ref := "{\"crs\":{\"type\":\"name\",\"properties\":{\"name\":\"urn:ogc:def:crs:OGC::CRS84\"}},\"type\":\"Polygon\",\"coordinates\":[[[2,1],[3,-2],[4,-1]]]}"
	if strings.Compare(string(b), ref) != 0 {
		t.Fail()
	}
}

func TestMarshalFeature(t *testing.T) {
	f := new(Feature)
	geom := NewPoint(3.0, 4.0)
	prop := make(map[string]int64)
	prop["a"] = 49
	prop["b"] = 17
	f.Crs = *WGS84
	f.Properties = prop
	f.Geometry = geom
	b, err := MarshalFeature(*f)
	if err != nil {
		fmt.Println("error", err)
		t.Error()
	}
	ref := "{\"crs\":{\"type\":\"name\",\"properties\":{\"name\":\"urn:ogc:def:crs:OGC::CRS84\"}},\"type\":\"Feature\",\"geometry\":{\"crs\":{\"type\":\"name\",\"properties\":{\"name\":\"urn:ogc:def:crs:OGC::CRS84\"}},\"type\":\"Point\",\"coordinates\":[3,4]},\"properties\":{\"a\":49,\"b\":17}}"
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
		t.Fail()
	}
}
