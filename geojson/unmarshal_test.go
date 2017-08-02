package geojson

import (
	"fmt"
	"testing"
)

func TestUnmarshalInvalid(t *testing.T) {
	_, err := UnmarshalGeoJSON2([]byte(`{ "type": "FauxPoint", "coordinates": [100.0, 0.0] }`))
	if err == nil {
		t.Fail()
	}
}

func TestUnmarshalPoint(t *testing.T) {
	geo, err := UnmarshalGeoJSON2([]byte(`{ "type": "Point", "coordinates": [100.0, 0.0] }`))
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	if len(geo.Point.Coordinates) != 2 {
		t.Fail()
	}
	expected := []float64{100, 0}
	for i, v := range expected {
		if geo.Point.Coordinates[i] != v {
			t.Fail()
		}
	}
}

func TestUnmarshalLineString(t *testing.T) {
	geo, err := UnmarshalGeoJSON2([]byte(`{ "type": "LineString", "coordinates": [ [100.0, 0.0], [101.0, 1.0] ] }`))
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	if len(geo.LineString.Coordinates) != 2 {
		t.Fail()
	}
}

func TestUnmarshalPolygonNoHoles(t *testing.T) {
	geo, err := UnmarshalGeoJSON2([]byte(`{ "type": "Polygon", "coordinates": [ [ [100.0, 0.0], [101.0, 0.0], [101.0, 1.0], [100.0, 1.0], [100.0, 0.0] ] ] }`))
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	if len(geo.Polygon.Coordinates) != 1 {
		t.Fail()
	}
	if len(geo.Polygon.Coordinates[0]) != 5 {
		t.Fail()
	}
}

func TestUnmarshalPolygonHoles(t *testing.T) {
	geo, err := UnmarshalGeoJSON2([]byte(`{ "type": "Polygon", "coordinates": [ [ [100.0, 0.0], [101.0, 0.0], [101.0, 1.0], [100.0, 1.0], [100.0, 0.0] ],
      [ [100.2, 0.2], [100.8, 0.2], [100.8, 0.8], [100.2, 0.8], [100.2, 0.2] ] ] }`))
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	if len(geo.Polygon.Coordinates) != 2 {
		t.Fail()
	}
	if len(geo.Polygon.Coordinates[0]) != 5 {
		t.Fail()
	}
	if len(geo.Polygon.Coordinates[1]) != 5 {
		t.Fail()
	}
}

func TestUnmarshalMultiPoint(t *testing.T) {
	geo, err := UnmarshalGeoJSON2([]byte(`{"type": "MultiPoint", "coordinates": [ [100.0, 0.0], [101.0, 1.0] ] }`))
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	if len(geo.MultiPoint.Coordinates) != 2 {
		t.Fail()
	}
}

func TestUnmarshalMultiLineString(t *testing.T) {
	geo, err := UnmarshalGeoJSON2([]byte(`{"type": "MultiLineString",
    "coordinates": [
        [ [100.0, 0.0], [101.0, 1.0] ],
        [ [102.0, 2.0], [103.0, 3.0] ]
      ]
    }`))
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	if len(geo.MultiLineString.Coordinates) != 2 {
		t.Fail()
	}
}

func TestUnmarshalMultiPolygon(t *testing.T) {
	geo, err := UnmarshalGeoJSON2([]byte(`{"type": "MultiPolygon",
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
	if len(geo.MultiPolygon.Coordinates) != 2 {
		t.Fail()
	}
}

func TestUnmarshalGeometryCollection(t *testing.T) {
	geo, err := UnmarshalGeoJSON2([]byte(`{ "type": "GeometryCollection",
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
	if len(geo.GeometryCollection.Geometries) != 2 {
		t.Fail()
	}
}

func TestUnmarshalFeature(t *testing.T) {
	geo, err := UnmarshalGeoJSON2([]byte(`{ "type": "Feature",
         "geometry": {"type": "Point", "coordinates": [102.0, 0.5]},
         "properties": {"prop0": "value0"}
         }`))
	if err != nil {
		fmt.Println("error:", err)
		t.Error()
	}
	if geo.Feature.Geometry.Type != "Point" {
		t.Fail()
	}
	if len(geo.Feature.Geometry.Point.Coordinates) != 2 {
		t.Fail()
	}
	if geo.Feature.Properties["prop0"] != "value0" {
		t.Fail()
	}
}

func TestUnmarshalFeatureCollection(t *testing.T) {
	geo, err := UnmarshalGeoJSON2([]byte(`{ "type": "FeatureCollection",
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
	if len(geo.FeatureCollection.Features) != 3 {
		t.Fail()
	}
}
