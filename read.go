/*
 * Implements functions for unmarshaling GeoJSON bytes to structs
 */
package geojson

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Option type for returning the components of a GeoJSON object
type GeoJSONContents struct {
	Points           []Point
	LineStrings      []LineString
	Polygons         []Polygon
	MultiPoints      []MultiPoint
	MultiLineStrings []MultiLineString
	MultiPolygons    []MultiPolygon
	Features         []Feature
}

// unknownGeoJSONType is used to read in only the type member of a GeoJSON
// object
type unknownGeoJSONType struct {
	Type string `json:"type"`
}

// partialGeometryCollection is used to obtain a list of byte arrays
// representing the Geometries member of a GeometryCollection
type partialGeometryCollection struct {
	Type       string            `json:"type"`
	Geometries []json.RawMessage `json:"geometries"`
}

// UnmarshalGeoJSON unpacks Features and Geometries from a JSON byte array
func UnmarshalGeoJSON(data []byte) (GeoJSONContents, error) {
	var uknType unknownGeoJSONType
	var result GeoJSONContents
	err := json.Unmarshal(data, &uknType)
	if err != nil {
		return result, err
	}
	switch uknType.Type {
	case "Point":
		var pt Point
		err = json.Unmarshal(data, &pt)
		if err != nil {
			return result, errors.New("invalid GeoJSON: malformed Point")
		}
		result.Points = append(result.Points, pt)

	case "LineString":
		var ls LineString
		err = json.Unmarshal(data, &ls)
		if err != nil {
			return result, errors.New("invalid GeoJSON: malformed LineString")
		}
		result.LineStrings = append(result.LineStrings, ls)

	case "Polygon":
		var poly Polygon
		err = json.Unmarshal(data, &poly)
		if err != nil {
			return result, errors.New("invalid GeoJSON: malformed Polygon")
		}
		result.Polygons = append(result.Polygons, poly)

	case "MultiPoint":
		var mp MultiPoint
		err = json.Unmarshal(data, &mp)
		if err != nil {
			return result, errors.New("invalid GeoJSON: malformed MultiPoint")
		}
		result.MultiPoints = append(result.MultiPoints, mp)

	case "MultiLineString":
		var mls MultiLineString
		err = json.Unmarshal(data, &mls)
		if err != nil {
			return result, errors.New("invalid GeoJSON: malformed MultiLineString")
		}
		result.MultiLineStrings = append(result.MultiLineStrings, mls)

	case "MultiPolygon":
		var mpoly MultiPolygon
		err = json.Unmarshal(data, &mpoly)
		if err != nil {
			return result, errors.New("invalid GeoJSON: malformed MultiPolygon")
		}
		result.MultiPolygons = append(result.MultiPolygons, mpoly)

	case "GeometryCollection":

		var partial partialGeometryCollection
		err = json.Unmarshal(data, &partial)
		if err != nil {
			return result, errors.New("invalid GeoJSON: malformed GeometryCollection")
		}

		var subresult GeoJSONContents
		for i := 0; i != len(partial.Geometries); i++ {
			subresult, err = UnmarshalGeoJSON(partial.Geometries[i])
			if err != nil {
				return result, err
			}

			result.Points = append(result.Points, subresult.Points...)
			result.LineStrings = append(result.LineStrings, subresult.LineStrings...)
			result.Polygons = append(result.Polygons, subresult.Polygons...)
			result.MultiPoints = append(result.MultiPoints, subresult.MultiPoints...)
			result.MultiLineStrings = append(result.MultiLineStrings, subresult.MultiLineStrings...)
			result.MultiPolygons = append(result.MultiPolygons, subresult.MultiPolygons...)
		}

	case "Feature":
		var feature Feature
		err = json.Unmarshal(data, &feature)
		if err != nil {
			return result, errors.New("invalid GeoJSON: malformed Feature")
		}
		result.Features = append(result.Features, feature)

	case "FeatureCollection":
		var featureCollection FeatureCollection
		err = json.Unmarshal(data, &featureCollection)
		if err != nil {
			return result, errors.New("invalid GeoJSON: malformed FeatureCollection")
		}
		result.Features = append(result.Features, featureCollection.Features...)
	default:
		return result, errors.New(fmt.Sprintf("unrecognized type: %s", uknType.Type))
	}
	return result, nil
}
