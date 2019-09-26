package mathutil

// InEpsilon returns if two values are within the Epsilon of each other absolutely.
func InEpsilon(a, b float64) bool {
	return (a-b) < Epsilon && (b-a) < Epsilon
}
