package money

import "math"

func roundToTwo(value float64) float64 {
	return math.Round(value*100) / 100
}
