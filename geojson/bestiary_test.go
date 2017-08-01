package geojson

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestBestiary100(t *testing.T) {
	b, err := ioutil.ReadFile("testdata/bestiary-100.json")
	geo, err := UnmarshalGeoJSON2(b)
	if err != nil {
		fmt.Println(err)
		t.Error()
	}
	g := geo.FeatureCollection.Features[0].Geometry
	fmt.Println(g)
	if g.Type != "GeometryCollection" {
		t.Fail()
	}
}

func TestBestiary200(t *testing.T) {
	b, err := ioutil.ReadFile("testdata/bestiary-200.json")
	geo, err := UnmarshalGeoJSON2(b)
	if err != nil {
		fmt.Println(err)
		t.Error()
	}
	g := geo.FeatureCollection.Features[0].Geometry
	if g.Type != "GeometryCollection" {
		t.Fail()
	}
}
