package geojson

import (
	"encoding/json"
	"fmt"
)

func AsGeoJSON(g Geometry) ([]byte, error) {
	b, err := json.Marshal(g)
	if err != nil {
		fmt.Println("error", err)
	}
	return b, err
}
