package geojson

func NewPoint(x ...float64) *Point {
	var pos Position
	if len(x) == 2 {
		pos = &Position2{x[0], x[1]}
	} else if len(x) == 3 {
		pos = &Position3{x[0], x[1], x[2]}
	} else {
		panic("too many arguments")
	}
	crs := NameCRS("urn:ogc:def:crs:OGC:1.3:CRS84")
	g := new(Point)
	g.Type = "Point"
	g.Coordinates = pos
	g.Crs = *crs
	return g
}

func NewMultiPoint(x ...[]float64) *MultiPoint {
	var ivert int
	var nVertices int
	var pos Position
	var coordinates []Position

	if len(x) == 2 {

		nVertices = len(x[0])
		coordinates = make([]Position, nVertices)
		for ivert = 0; ivert != nVertices; ivert++ {
			pos = &Position2{x[0][ivert], x[1][ivert]}
			coordinates[ivert] = pos
		}

	} else if len(x) == 3 {

		nVertices = len(x[0])
		coordinates = make([]Position, nVertices)
		for ivert = 0; ivert != nVertices; ivert++ {
			pos = &Position3{x[0][ivert], x[1][ivert], x[2][ivert]}
			coordinates[ivert] = pos
		}

	} else {
		panic("NewLineString called with odd number of arguments")
	}

	crs := NameCRS("urn:ogc:def:crs:OGC:1.3:CRS84")
	g := new(MultiPoint)
	g.Type = "MultiPoint"
	g.Coordinates = coordinates
	g.Crs = *crs
	return g
}

func NewLineString(x ...[]float64) *LineString {
	var ivert int
	var nVertices int
	var pos Position
	var coordinates []Position

	if len(x) == 2 {

		nVertices = len(x[0])
		coordinates = make([]Position, nVertices)
		for ivert = 0; ivert != nVertices; ivert++ {
			pos = &Position2{x[0][ivert], x[1][ivert]}
			coordinates[ivert] = pos
		}

	} else if len(x) == 3 {

		nVertices = len(x[0])
		coordinates = make([]Position, nVertices)
		for ivert = 0; ivert != nVertices; ivert++ {
			pos = &Position3{x[0][ivert], x[1][ivert], x[2][ivert]}
			coordinates[ivert] = pos
		}

	} else {
		panic("NewLineString called with odd number of arguments")
	}

	crs := NameCRS("urn:ogc:def:crs:OGC:1.3:CRS84")
	g := new(LineString)
	g.Type = "LineString"
	g.Coordinates = coordinates
	g.Crs = *crs
	return g
}

// FIXME: bug - if supplied with 2m*3n arguments (e.g. 6), ambiguous as to
// whether intended geometry is 2d or 3d
// proposed fix: specialize 2d and 3d constructor functions
func NewMultiLineString2(x ...[]float64) *MultiLineString {
	var ip, ivert int
	var nParts, nVertices int
	var pos Position
	var coordinates [][]Position

	if (len(x) % 2) == 0 {

		nParts = len(x) / 2
		coordinates = make([][]Position, nParts)
		for ip = 0; ip != nParts; ip++ {
			nVertices = len(x[ip*2])
			coordinates[ip] = make([]Position, nVertices)
			for ivert = 0; ivert != nVertices; ivert++ {
				pos = &Position2{x[ip*2][ivert], x[ip*2+1][ivert]}
				coordinates[ip][ivert] = pos
			}
		}

	} else {
		panic("NewMultiLineString2 called with odd number of arguments")
	}

	crs := NameCRS("urn:ogc:def:crs:OGC:1.3:CRS84")
	g := new(MultiLineString)
	g.Type = "MultiLineString"
	g.Coordinates = coordinates
	g.Crs = *crs
	return g
}

// NewPolygon2 is a convenience constructor for a 2D Polygon. It is called as
// NewPolygon2(x, y, [x_sub1, y_sub1, [x_sub2, y_sub2]]...) where areguments
// are slices of floats.
func NewPolygon2(x ...[]float64) *Polygon {
	var ip, ivert int
	var nParts, nVertices int
	var pos Position
	var coordinates [][]Position

	if (len(x) % 2) == 0 {

		nParts = len(x) / 2
		coordinates = make([][]Position, nParts)
		for ip = 0; ip != nParts; ip++ {
			nVertices = len(x[ip*2])
			coordinates[ip] = make([]Position, nVertices)
			for ivert = 0; ivert != nVertices; ivert++ {
				pos = &Position2{x[ip*2][ivert], x[ip*2+1][ivert]}
				coordinates[ip][ivert] = pos
			}
		}

	} else {
		panic("NewPolygon called with odd number of arguments")
	}

	crs := NameCRS("urn:ogc:def:crs:OGC:1.3:CRS84")
	g := new(Polygon)
	g.Type = "MultiLineString"
	g.Coordinates = coordinates
	g.Crs = *crs
	return g
}
