package geojson

import (
	"fmt"
	"testing"
)

func TestMarshallPoint(t *testing.T) {
	point := NewPoint(3.0, 4.0)
	b, err := AsGeoJSON(point)
	if err != nil {
		fmt.Println("error", err)
		t.Fail()
	}
	fmt.Println(string(b))
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
	fmt.Println(string(b))
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
	fmt.Println(string(b))
}
