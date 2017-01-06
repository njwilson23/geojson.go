package geojson

import (
	"fmt"
	"io/ioutil"
	"testing"
)

/*func beastiaryDepth(content GeoJSONContents) int {
	var c *GeometryCollection
	c = GeoJSONContents.GeometryCollection[0]
	cnt := 0
	itGoesDeeper := true
	var ok bool
	for itGoesDeeper {
		c = c.Geometries[0]
		_, itGoesDeeper = c.Geometries[0].(GeometryCollection)
		cnt++
	}
	return cnt
}*/

func TestBestiary100(t *testing.T) {
	b, err := ioutil.ReadFile("testdata/bestiary-100.json")
	content, err := UnmarshalGeoJSON(b)
	if err != nil {
		fmt.Println(err)
		t.Error()
	}
	g := content.Features[0].Geometry
	gc, ok := g.(*GeometryCollection)
	if !ok {
		t.Fail()
	}
	if len(gc.Geometries) != 2 {
		t.Fail()
	}
}

func TestBestiary200(t *testing.T) {
	b, err := ioutil.ReadFile("testdata/bestiary-200.json")
	content, err := UnmarshalGeoJSON(b)
	if err != nil {
		fmt.Println(err)
		t.Error()
	}
	g := content.Features[0].Geometry
	gc, ok := g.(*GeometryCollection)
	if !ok {
		t.Fail()
	}
	if len(gc.Geometries) != 2 {
		t.Fail()
	}
}
