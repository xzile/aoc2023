package ez

func Sum[S ~[]E, E ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64](s S) E {
	var e E
	for _, v := range s {
		e += v
	}
	return e
}

// GCD & LCM copied from https://siongui.github.io/2017/06/03/go-find-lcm-by-gcd/

// GCD greatest common divisor via Euclidean algorithm
func GCD(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// LCM find Least Common Multiple via GCD
func LCM(a, b int, integers ...int) int {
	result := a * b / GCD(a, b)

	for i := 0; i < len(integers); i++ {
		result = LCM(result, integers[i])
	}

	return result
}

// Shoelace and Point copied from https://rosettacode.org/wiki/Shoelace_formula_for_polygonal_area#Go

// Point represents an x and y coordinate for a point along a polynomial
type Point struct{ X, Y float64 }

// Shoelace calculates the area of a polynomial given a set of points
func Shoelace(pts []Point) float64 {
	sum := 0.
	p0 := pts[len(pts)-1]
	for _, p1 := range pts {
		sum += p0.Y*p1.X - p0.X*p1.Y
		p0 = p1
	}
	return sum / 2
}

// Picks is a modified version of the Picks theorem formula to calculate the inner points
// https://en.wikipedia.org/wiki/Pick%27s_theorem
func Picks(area float64, pointCount int) float64 {
	return area + float64(1) - float64(pointCount/2)
}
