package pkg1

//
import (
	"math"
	//
	//
	"github.com/samples/sample1/pkg2"
)

//

func Add(x, y int) int {
	return x + y
}

func AddAndSubtract(x, y int) (int, int) {
	return Add(x, y), pkg2.Subtract(x, y)
}

func Floor(x float64) float64 {
	return math.Floor(x)
}
