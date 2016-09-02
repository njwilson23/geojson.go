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
	ref := "{\"crs\":{\"type\":\"name\",\"properties\":{\"name\":\"urn:ogc:def:crs:OGC:1.3:CRS84\"}},\"type\":\"Point\",\"coordinates\":[3.000000,4.000000]}"
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
	ref := "{\"crs\":{\"type\":\"name\",\"properties\":{\"name\":\"urn:ogc:def:crs:OGC:1.3:CRS84\"}},\"type\":\"LineString\",\"coordinates\":[[2.000000,1.000000],[3.000000,-2.000000],[4.000000,-1.000000]]}"
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
	ref := "{\"crs\":{\"type\":\"name\",\"properties\":{\"name\":\"urn:ogc:def:crs:OGC:1.3:CRS84\"}},\"type\":\"MultiLineString\",\"coordinates\":[[[2.000000,1.000000],[3.000000,-2.000000],[4.000000,-1.000000]]]}"
	if strings.Compare(string(b), ref) != 0 {
		t.Fail()
	}
}
