package geojson

import (
	"encoding/json"
	"fmt"
)

func (p *Position2) MarshalJSON() ([]byte, error) {
	b := make([]byte, 24)
	b = []byte(fmt.Sprintf("[%.6f, %.6f]", p.x, p.y))
	return b, nil
}

func (p *Position3) MarshalJSON() ([]byte, error) {
	b := make([]byte, 36)
	b = []byte(fmt.Sprintf("[%.6f, %.6f, %.6f]", p.x, p.y, p.z))
	return b, nil
}

func AsGeoJSON(g Geometry) ([]byte, error) {
	b, err := json.Marshal(g)
	if err != nil {
		fmt.Println("error", err)
	}
	return b, err
}
