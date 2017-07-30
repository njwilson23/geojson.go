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

type partialFeature struct {
	CRSReferencable
	Type       string                 `json:"type"`
	ID         string                 `json:"id,omitempty"`
	Properties map[string]interface{} `json:"properties"`
	Geometry   json.RawMessage        `json:"geometry"`
}

type partialFeatureCollection struct {
	CRSReferencable
	Type     string            `json:"type"`
	Features []json.RawMessage `json:"features"`
}

func (g *GeoJSONContents) String() string {
	return fmt.Sprintf("Points: %d\nLineStrings: %d\nPolygons: %d\nMultiPoints: %d\nMultiLineStrings: %d\nMultiPolygons: %d\nFeatures: %d\n",
		len(g.Points), len(g.LineStrings), len(g.Polygons), len(g.MultiPoints),
		len(g.MultiLineStrings), len(g.MultiPolygons), len(g.Features))
}

// CoalesceGeometry returns a Geometry containing all Geometry-types
func (g *GeoJSONContents) CoalesceGeometry() Geometry {
	cnt := len(g.Points) + len(g.LineStrings) + len(g.Polygons) + len(g.MultiPoints) + len(g.MultiLineStrings) + len(g.MultiPolygons)
	var retval Geometry
	if cnt == 1 {
		if len(g.Points) == 1 {
			retval = &g.Points[0]
		} else if len(g.LineStrings) == 1 {
			retval = &g.LineStrings[0]
		} else if len(g.Polygons) == 1 {
			retval = &g.Polygons[0]
		} else if len(g.MultiPoints) == 1 {
			retval = &g.MultiPoints[0]
		} else if len(g.MultiLineStrings) == 1 {
			retval = &g.MultiLineStrings[0]
		} else if len(g.MultiPolygons) == 1 {
			retval = &g.MultiPolygons[0]
		}
	} else {
		gc := new(GeometryCollection)
		for i := 0; i != len(g.Points); i++ {
			gc.Geometries = append(gc.Geometries, &g.Points[i])
		}
		for i := 0; i != len(g.LineStrings); i++ {
			gc.Geometries = append(gc.Geometries, &g.LineStrings[i])
		}
		for i := 0; i != len(g.Polygons); i++ {
			gc.Geometries = append(gc.Geometries, &g.Polygons[i])
		}
		for i := 0; i != len(g.MultiPoints); i++ {
			gc.Geometries = append(gc.Geometries, &g.MultiPoints[i])
		}
		for i := 0; i != len(g.MultiLineStrings); i++ {
			gc.Geometries = append(gc.Geometries, &g.MultiLineStrings[i])
		}
		for i := 0; i != len(g.MultiPolygons); i++ {
			gc.Geometries = append(gc.Geometries, &g.MultiPolygons[i])
		}
		retval = gc
	}
	return retval
}

// UnmarshalGeoJSON unpacks Features and Geometries from a JSON byte array The
// ouput, GeoJSONContents, is a flat listing of all Geometry and Feature objects
// in the input bytes. Nested structure of the original GeoJSON is not
// preserved.
func UnmarshalGeoJSON(data []byte) (result GeoJSONContents, err error) {

	defer func() {
		if r := recover(); r != nil {
			_, ok := r.(error)
			if !ok {
				err = errors.New("unknown panic occurred")
			} else {
				err = r.(error)
			}
		}
	}()

	var uknType unknownGeoJSONType
	err = json.Unmarshal(data, &uknType)
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

		ch := make(chan GeoJSONContents)

		for i := 0; i != len(partial.Geometries); i++ {
			go func(g []byte, ch chan GeoJSONContents) {
				subresult, err := UnmarshalGeoJSON(g)
				if err != nil {
					panic(errors.New("failure unmarshalling child geometry"))
				}
				ch <- subresult
			}(partial.Geometries[i], ch)
		}

		var subresult GeoJSONContents
		for i := 0; i != len(partial.Geometries); i++ {
			subresult = <-ch
			result.Points = append(result.Points, subresult.Points...)
			result.LineStrings = append(result.LineStrings, subresult.LineStrings...)
			result.Polygons = append(result.Polygons, subresult.Polygons...)
			result.MultiPoints = append(result.MultiPoints, subresult.MultiPoints...)
			result.MultiLineStrings = append(result.MultiLineStrings, subresult.MultiLineStrings...)
			result.MultiPolygons = append(result.MultiPolygons, subresult.MultiPolygons...)
		}

	case "Feature":
		var partial partialFeature
		err = json.Unmarshal(data, &partial)
		if err != nil {
			return result, errors.New("invalid GeoJSON: malformed Feature")
		}

		var subresult GeoJSONContents
		subresult, err = UnmarshalGeoJSON(partial.Geometry)
		if err != nil {
			return result, err
		}

		feature := new(Feature)
		feature.CRS = partial.CRS
		feature.ID = partial.ID
		feature.Properties = partial.Properties
		feature.Geometry = subresult.CoalesceGeometry()
		result.Features = append(result.Features, *feature)

	case "FeatureCollection":
		var partial partialFeatureCollection
		err = json.Unmarshal(data, &partial)
		if err != nil {
			fmt.Println(err)
			return result, errors.New("invalid GeoJSON: malformed FeatureCollection")
		}

		ch := make(chan GeoJSONContents)

		for i := 0; i != len(partial.Features); i++ {
			go func(b []byte, ch chan GeoJSONContents) {
				subresult, err := UnmarshalGeoJSON(b)
				if err != nil {
					panic(errors.New("failure unmarshalling child feature"))
				}
				ch <- subresult
			}(partial.Features[i], ch)
		}

		var subresult GeoJSONContents
		for i := 0; i != len(partial.Features); i++ {
			subresult = <-ch
			result.Features = append(result.Features, subresult.Features[0])
		}

	default:
		return result, errors.New(fmt.Sprintf("unrecognized type: %s", uknType.Type))
	}
	return result, err
}
