package utils

import (
	"math/rand"
	"time"
)

func RandInt(min, max int) int {
	if min >= max || max == 0 {
		return max
	}
	rand.Seed(time.Now().Local().UnixNano())
	num := rand.Intn(max-min) + min
	return num
}
