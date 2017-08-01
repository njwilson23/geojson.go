/*
 * Implements functions for unmarshaling GeoJSON bytes to structs
 */
package geojson

import (
	"encoding/json"
	"fmt"
)

func (g *Geo) UnmarshalJSON(data []byte) (err error) {

	g.Point = new(Point)
	g.LineString = new(LineString)
	g.Polygon = new(Polygon)
	g.MultiPoint = new(MultiPoint)
	g.MultiLineString = new(MultiLineString)
	g.MultiPolygon = new(MultiPolygon)
	g.Feature = new(Feature)
	g.GeometryCollection = new(GeometryCollection)
	g.FeatureCollection = new(FeatureCollection)

	partial := &struct {
		Type string `json:"type"`
	}{}
	json.Unmarshal(data, partial)
	g.Type = partial.Type
	switch partial.Type {
	case "Point":
		err = json.Unmarshal(data, g.Point)
	case "LineString":
		err = json.Unmarshal(data, g.LineString)
	case "Polygon":
		err = json.Unmarshal(data, g.Polygon)
	case "MultiPoint":
		err = json.Unmarshal(data, g.MultiPoint)
	case "MultiLineString":
		err = json.Unmarshal(data, g.MultiLineString)
	case "MultiPolygon":
		err = json.Unmarshal(data, g.MultiPolygon)
	case "Feature":
		err = json.Unmarshal(data, g.Feature)
	case "GeometryCollection":
		err = json.Unmarshal(data, g.GeometryCollection)
	case "FeatureCollection":
		err = json.Unmarshal(data, g.FeatureCollection)
	default:
		err = fmt.Errorf("unhandled object type: '%s'", partial.Type)
	}
	return
}

func UnmarshalGeoJSON2(data []byte) (*Geo, error) {
	g := new(Geo)
	err := g.UnmarshalJSON(data)
	return g, err
}
