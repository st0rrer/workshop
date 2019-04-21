package report

import (
	"math"
	"math/rand"
)

func RandomInt(min, max, base int) int {

	if max < 0 || max < base {
		panic("invalid argument to RandomInt")
	}

	randInt := min + rand.Intn(max-min)
	if randInt == 0 {
		return 0
	} else if base == 1 {
		return randInt
	}

	power := math.Trunc(math.Log(float64(randInt)) / math.Log(float64(base)))
	return int(math.Pow(float64(base), power))
}

func RandomLatitude() float64 {
	return rand.Float64()*180 - 90
}

func RandomLongitude() float64 {
	return rand.Float64()*360 - 180
}
