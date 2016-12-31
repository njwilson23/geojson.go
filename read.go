package geojson

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Option type for reading GeoJSON data
type GeoJSONResult struct {
	Features   []Feature
	Geometries []Geometry
}

type unknownGeoJSONType struct {
	Type string `json:"type"`
}

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
		return NewPoint(0, 0), err
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
		return NewPoint(0, 0), errors.New(fmt.Sprintf("unrecognized GeoJSON type '%s'", uknType.Type))
	}
}

// func ParseGeoJSON(stream []byte) (GeoJSONResult, error) {
// 	var result GeoJSONResult
// 	streamMap := make(map[string]interface{})
// 	err := json.Unmarshal(stream, streamMap)
// 	if err != nil {
// 		return result, err
// 	}
//
// 	var geo Geometry
// 	var geoCollection *GeometryCollection
// 	var feature *Feature
// 	var featureCollection *FeatureCollection
//
// 	// FIXME
// 	data, ok := streamMap["type"]
// 	if !ok {
// 		return result, errors.New("invalid GeoJSON: JSON object missing 'type' key")
// 	}
//
// 	_type, ok := data.(string)
// 	if !ok {
// 		return result, errors.New("invalid GeoJSON: 'type' value not a string")
// 	}
// 	fmt.Println(_type)
// 	/*switch type _type {
// 	case string:
// 		break
// 	default:
// 		return result, new.Errors("type attribute not string")
// 	}*/
//
// 	switch _type {
// 	case "Point":
// 		geo, err = parsePoint(streamMap)
// 		if err != nil {
// 			return result, errors.New("invalid GeoJSON: malformed Point object")
// 		}
// 		result.Geometries = append(result.Geometries, geo)
// 	case "MultiPoint":
// 		geo, err = parseMultiPoint(streamMap)
// 		if err != nil {
// 			return result, errors.New("invalid GeoJSON: malformed MultiPoint object")
// 		}
// 		result.Geometries = append(result.Geometries, geo)
// 	case "LineString":
// 		geo, err = parseLineString(streamMap)
// 		if err != nil {
// 			return result, errors.New("invalid GeoJSON: malformed LineString object")
// 		}
// 		result.Geometries = append(result.Geometries, geo)
// 	case "MultiLineString":
// 		geo, err = parseMultiLineString(streamMap)
// 		if err != nil {
// 			return result, errors.New("invalid GeoJSON: malformed MultiLineString object")
// 		}
// 		result.Geometries = append(result.Geometries, geo)
// 	case "Polygon":
// 		geo, err = parsePolygon(streamMap)
// 		if err != nil {
// 			return result, errors.New("invalid GeoJSON: malformed Polygon object")
// 		}
// 		result.Geometries = append(result.Geometries, geo)
// 	//case "MultiPolygon":
// 	//	geo, err = parseMultiPolygon(streamMap)
// 	//	if err != nil {
// 	//		return result, errors.New("invalid GeoJSON: malformed MultiPolygon object")
// 	//	}
// 	//	result.Geometries = append(result.Geometries, geo)
// 	case "GeometryCollection":
// 		geoCollection, err = parseGeometryCollection(streamMap)
// 		if err != nil {
// 			err = errors.New(fmt.Sprintf("%s\n  %s", err,
// 				"invalid GeoJSON: malformed GeometryCollection object"))
// 			return result, err
// 		}
// 		result.GeometryCollections = append(result.GeometryCollections, *geoCollection)
// 	case "Feature":
// 		feature, err = parseFeature(streamMap)
// 		if err != nil {
// 			err = errors.New(fmt.Sprintf("%s\n  %s", err,
// 				"invalid GeoJSON: malformed Feature object"))
// 			return result, err
// 		}
// 		result.Features = append(result.Features, *feature)
// 	case "FeatureCollection":
// 		featureCollection, err = parseFeatureCollection(streamMap)
// 		if err != nil {
// 			err = errors.New(fmt.Sprintf("%s\n  %s", err,
// 				"invalid GeoJSON: malformed FeatureCollection object"))
// 			return result, err
// 		}
// 		result.FeatureCollections = append(result.FeatureCollections, *featureCollection)
// 	}
//
// 	return result, nil
// }
//
// func parsePoint(streamMap map[string]interface{}) (*Point, error) {
// 	data, ok := streamMap["coordinates"]
// 	if !ok {
// 		return nil, errors.New("Point object missing 'coordinates'")
// 	}
// 	coords, ok := data.([]float64)
// 	if !ok {
// 		return nil, errors.New("Point object coordinates of unexpected type")
// 	}
// 	point := NewPoint(coords...)
// 	return point, nil
// }
//
// func parseMultiPoint(streamMap map[string]interface{}) (*MultiPoint, error) {
// 	data, ok := streamMap["coordinates"]
// 	if !ok {
// 		return nil, errors.New("Point object missing 'coordinates'")
// 	}
// 	coords, ok := data.([][]float64)
// 	if !ok {
// 		return nil, errors.New("MultiPoint object coordinates of unexpected type")
// 	}
// 	mp := NewMultiPoint(coords...)
// 	return mp, nil
// }
//
// func parseLineString(streamMap map[string]interface{}) (*LineString, error) {
// 	data, ok := streamMap["coordinates"]
// 	if !ok {
// 		return nil, errors.New("LineString object missing 'coordinates'")
// 	}
// 	coords, ok := data.([][]float64)
// 	if !ok {
// 		return nil, errors.New("LineString object coordinates of unexpected type")
// 	}
// 	linestring := NewLineString(coords...)
// 	return linestring, nil
// }
//
// func parseMultiLineString(streamMap map[string]interface{}) (*MultiLineString, error) {
// 	data, ok := streamMap["coordinates"]
// 	if !ok {
// 		return nil, errors.New("MultiLineString object missing 'coordinates'")
// 	}
// 	coords, ok := data.([][]float64)
// 	if !ok {
// 		return nil, errors.New("MultiLineString object coordinates of unexpected type")
// 	}
// 	multilinestring := NewMultiLineString2(coords...)
// 	return multilinestring, nil
// }
//
// func parsePolygon(streamMap map[string]interface{}) (*Polygon, error) {
// 	data, ok := streamMap["coordinates"]
// 	if !ok {
// 		return nil, errors.New("Polygon object missing 'coordinates'")
// 	}
// 	coords, ok := data.([][]float64)
// 	if !ok {
// 		return nil, errors.New("Polygon object coordinates of unexpected type")
// 	}
// 	polygon := NewPolygon2(coords...)
// 	return polygon, nil
// }
//
// // func parseMultiPolygon(streamMap map[string]interface{}) (*MultiPolygon, error) {
// // 	data, ok := streamMap["coordinates"]
// // 	if !ok {
// // 		return nil, errors.New("MultiPolygon object missing 'coordinates'")
// // 	}
// // 	coords, ok := data.([][][]float64)
// // 	if !ok {
// // 		return nil, errors.New("MultiPolygon object coordinates of unexpected type")
// // 	}
// // 	multipolygon := NewMultiPolygon2(coords...)
// // 	return multipolygon, nil
// // }
//
// func parseGeometryCollection(streamMap map[string]interface{}) (*GeometryCollection, error) {
// 	return nil, errors.New("not implemented")
// }
//
// func parseFeature(streamMap map[string]interface{}) (*Feature, error) {
// 	return nil, errors.New("not implemented")
// }
//
// func parseFeatureCollection(streamMap map[string]interface{}) (*FeatureCollection, error) {
// 	return nil, errors.New("not implemented")
// }
