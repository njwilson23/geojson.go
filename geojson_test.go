package geojson

import (
	"fmt"
	"strings"
	"testing"
)

func TestMarshallPoint(t *testing.T) {
	point := NewPoint(3.0, 4.0)
	b, err := AsGeoJSON(point)
	if err != nil {
		fmt.Println("error", err)
		t.Fail()
	}
	ref := "{\"crs\":{\"type\":\"name\",\"properties\":{\"name\":\"urn:ogc:def:crs:OGC:1.3:CRS84\"}},\"type\":\"Point\",\"coordinates\":[3,4]}"
	if strings.Compare(string(b), ref) != 0 {
		t.Fail()
	}
}

func TestMarshallLineString(t *testing.T) {
	X := []float64{2.0, 3.0, 4.0}
	Y := []float64{1.0, -2.0, -1.0}
	point := NewLineString(X, Y)
	b, err := AsGeoJSON(point)
	if err != nil {
		fmt.Println("error", err)
		t.Fail()
	}
	ref := "{\"crs\":{\"type\":\"name\",\"properties\":{\"name\":\"urn:ogc:def:crs:OGC:1.3:CRS84\"}},\"type\":\"LineString\",\"coordinates\":[[2,1],[3,-2],[4,-1]]}"
	if strings.Compare(string(b), ref) != 0 {
		t.Fail()
	}
}

func TestMarshallPolygon(t *testing.T) {
	X := []float64{2.0, 3.0, 4.0}
	Y := []float64{1.0, -2.0, -1.0}
	point := NewPolygon2(X, Y)
	b, err := AsGeoJSON(point)
	if err != nil {
		fmt.Println("error", err)
		t.Fail()
	}
	ref := "{\"crs\":{\"type\":\"name\",\"properties\":{\"name\":\"urn:ogc:def:crs:OGC:1.3:CRS84\"}},\"type\":\"MultiLineString\",\"coordinates\":[[[2,1],[3,-2],[4,-1]]]}"
	if strings.Compare(string(b), ref) != 0 {
		t.Fail()
	}
}

func TestUnmarshallInvalid(t *testing.T) {
	_, err := UnmarshalGeoJSON([]byte(`{ "type": "FauxPoint", "coordinates": [100.0, 0.0] }`))
	if err == nil {
		t.Fail()
	}
}

func TestUnmarshallPoint(t *testing.T) {
	_, err := UnmarshalGeoJSON([]byte(`{ "type": "Point", "coordinates": [100.0, 0.0] }`))
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
}

func TestUnmarshallLineString(t *testing.T) {
	_, err := UnmarshalGeoJSON([]byte(`{ "type": "LineString", "coordinates": [ [100.0, 0.0], [101.0, 1.0] ] }`))
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
}

func TestUnmarshallPolygonNoHoles(t *testing.T) {
	_, err := UnmarshalGeoJSON([]byte(`{ "type": "Polygon", "coordinates": [ [ [100.0, 0.0], [101.0, 0.0], [101.0, 1.0], [100.0, 1.0], [100.0, 0.0] ] ] }`))
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
}

func TestUnmarshallPolygonHoles(t *testing.T) {
	_, err := UnmarshalGeoJSON([]byte(`{ "type": "Polygon", "coordinates": [ [ [100.0, 0.0], [101.0, 0.0], [101.0, 1.0], [100.0, 1.0], [100.0, 0.0] ],
      [ [100.2, 0.2], [100.8, 0.2], [100.8, 0.8], [100.2, 0.8], [100.2, 0.2] ] ] }`))
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
}

func TestUnmarshallMultiPoint(t *testing.T) {
	_, err := UnmarshalGeoJSON([]byte(`{"type": "MultiPoint", "coordinates": [ [100.0, 0.0], [101.0, 1.0] ] }`))
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
}

func TestUnmarshallMultiLineString(t *testing.T) {
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

func TestUnmarshallMultiPolygon(t *testing.T) {
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

func TestUnmarshallGeometryCollection(t *testing.T) {
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
