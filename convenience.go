package geojson

// NameCRS returns a named CRS object
// Note that for RFC 7946 compliance, WGS84 may be used
func NameCRS(name string) *CRS {
	prop := make(map[string]string)
	prop["name"] = name
	return &CRS{"name", prop}
}

// LinkCRS returns a linked CRS object
// Note that for RFC 7946 compliance, WGS84 may be used
func LinkCRS(link string) *CRS {
	prop := make(map[string]string)
	prop["link"] = link
	return &CRS{"link", prop}
}

var WGS84 *CRS = NameCRS("urn:ogc:def:crs:OGC::CRS84")

// NewPoint creates a point with the provided coordinates
func NewPoint(x ...float64) *Point {
	g := new(Point)
	g.Coordinates = x
	g.Crs = *WGS84
	return g
}

func NewMultiPoint(x ...[]float64) *MultiPoint {
	var ivert int
	var nVertices int
	var pos []float64
	var coordinates [][]float64

	if len(x) == 2 {

		nVertices = len(x[0])
		coordinates = make([][]float64, nVertices)
		for ivert = 0; ivert != nVertices; ivert++ {
			pos = []float64{x[0][ivert], x[1][ivert]}
			coordinates[ivert] = pos
		}

	} else if len(x) == 3 {

		nVertices = len(x[0])
		coordinates = make([][]float64, nVertices)
		for ivert = 0; ivert != nVertices; ivert++ {
			pos = []float64{x[0][ivert], x[1][ivert], x[2][ivert]}
			coordinates[ivert] = pos
		}

	} else {
		panic("NewMultiPoint takes either 2 or 3 arguments of type []float64")
	}

	g := new(MultiPoint)
	g.Coordinates = coordinates
	g.Crs = *WGS84
	return g
}

func NewLineString(x ...[]float64) *LineString {
	var ivert int
	var nVertices int
	var pos []float64
	var coordinates [][]float64

	if len(x) == 2 {

		nVertices = len(x[0])
		coordinates = make([][]float64, nVertices)
		for ivert = 0; ivert != nVertices; ivert++ {
			pos = []float64{x[0][ivert], x[1][ivert]}
			coordinates[ivert] = pos
		}

	} else if len(x) == 3 {

		nVertices = len(x[0])
		coordinates = make([][]float64, nVertices)
		for ivert = 0; ivert != nVertices; ivert++ {
			pos = []float64{x[0][ivert], x[1][ivert], x[2][ivert]}
			coordinates[ivert] = pos
		}

	} else {
		panic("NewLineString takes either 2 or 3 arguments of type []float64")
	}

	g := new(LineString)
	g.Coordinates = coordinates
	g.Crs = *WGS84
	return g
}

// FIXME: bug - if supplied with 2m*3n arguments (e.g. 6), ambiguous as to
// whether intended geometry is 2d or 3d
// proposed fix: specialize 2d and 3d constructor functions
func NewMultiLineString2(x ...[]float64) *MultiLineString {
	var ip, ivert int
	var nParts, nVertices int
	var pos []float64
	var coordinates [][][]float64

	if (len(x) % 2) == 0 {

		nParts = len(x) / 2
		coordinates = make([][][]float64, nParts)
		for ip = 0; ip != nParts; ip++ {
			nVertices = len(x[ip*2])
			coordinates[ip] = make([][]float64, nVertices)
			for ivert = 0; ivert != nVertices; ivert++ {
				pos = []float64{x[ip*2][ivert], x[ip*2+1][ivert]}
				coordinates[ip][ivert] = pos
			}
		}

	} else {
		panic("NewMultiLineString2 called with odd number of arguments")
	}

	g := new(MultiLineString)
	g.Coordinates = coordinates
	g.Crs = *WGS84
	return g
}

// NewPolygon2 is a convenience constructor for a 2D Polygon. It is called as
// NewPolygon2(x, y, [x_sub1, y_sub1, [x_sub2, y_sub2]]...) where areguments
// are slices of floats.
func NewPolygon2(x ...[]float64) *Polygon {
	var ip, ivert int
	var nParts, nVertices int
	var pos []float64
	var coordinates [][][]float64

	if (len(x) % 2) == 0 {

		nParts = len(x) / 2
		coordinates = make([][][]float64, nParts)
		for ip = 0; ip != nParts; ip++ {
			nVertices = len(x[ip*2])
			coordinates[ip] = make([][]float64, nVertices)
			for ivert = 0; ivert != nVertices; ivert++ {
				pos = []float64{x[ip*2][ivert], x[ip*2+1][ivert]}
				coordinates[ip][ivert] = pos
			}
		}

	} else {
		panic("NewPolygon2 called with odd number of arguments")
	}

	g := new(Polygon)
	g.Coordinates = coordinates
	g.Crs = *WGS84
	return g
}
