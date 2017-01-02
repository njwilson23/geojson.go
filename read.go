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
type GeoJSONResult struct {
	Features   []Feature
	Geometries []Geometry
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
func UnmarshalGeoJSON(data []byte) (GeoJSONResult, error) {
	var uknType unknownGeoJSONType
	var result GeoJSONResult
	err := json.Unmarshal(data, &uknType)
	if err != nil {
		return result, err
	}
	switch uknType.Type {
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
		for i := 0; i != len(featureCollection.Features); i++ {
			result.Features = append(result.Features, featureCollection.Features[i])
		}
	default:
		g, err := UnmarshalGeometry(data)
		if err != nil {
			return result, err
		}
		result.Geometries = append(result.Geometries, g)
	}
	return result, nil
}

// UnmarshalGeometry recursively unpacks Geometries (including GeometryCollections) from a JSON byte array
func UnmarshalGeometry(data []byte) (Geometry, error) {
	var uknType unknownGeoJSONType
	err := json.Unmarshal(data, &uknType)
	if err != nil {
		return new(Point), err
	}
	switch uknType.Type {
	case "Point":
		var pt Point
		err = json.Unmarshal(data, &pt)
		if err != nil {
			return pt, errors.New("invalid GeoJSON: malformed Point")
		}
		return pt, nil
	case "LineString":
		var ls LineString
		err = json.Unmarshal(data, &ls)
		if err != nil {
			return ls, errors.New("invalid GeoJSON: malformed LineString")
		}
		return ls, nil
	case "Polygon":
		var poly Polygon
		err = json.Unmarshal(data, &poly)
		if err != nil {
			return poly, errors.New("invalid GeoJSON: malformed Polygon")
		}
		return poly, nil
	case "MultiPoint":
		var mp MultiPoint
		err = json.Unmarshal(data, &mp)
		if err != nil {
			return mp, errors.New("invalid GeoJSON: malformed MultiPoint")
		}
		return mp, nil
	case "MultiLineString":
		var mls MultiLineString
		err = json.Unmarshal(data, &mls)
		if err != nil {
			return mls, errors.New("invalid GeoJSON: malformed MultiLineString")
		}
		return mls, nil
	case "MultiPolygon":
		var mpoly MultiPolygon
		err = json.Unmarshal(data, &mpoly)
		if err != nil {
			return mpoly, errors.New("invalid GeoJSON: malformed MultiPolygon")
		}
		return mpoly, nil
	case "GeometryCollection":
		var collection GeometryCollection
		var partial partialGeometryCollection
		err = json.Unmarshal(data, &partial)
		if err != nil {
			fmt.Println(err)
			return collection, errors.New("invalid GeoJSON: malformed GeometryCollection")
		}
		var geom Geometry
		for i := 0; i != len(partial.Geometries); i++ {
			geom, err = UnmarshalGeometry(partial.Geometries[i])
			if err != nil {
				return collection, err
			}
			collection.Geometries = append(collection.Geometries, geom)
		}
		return collection, nil
	default:
		return new(Point), errors.New(fmt.Sprintf("unrecognized GeoJSON type '%s'", uknType.Type))
	}
}
