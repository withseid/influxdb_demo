package utils

import (
	"math/rand"
	"time"
)

func GenFloatNum(min, max int) float64 {
	rand.Seed(time.Now().UnixNano())
	randNum := rand.Intn(max-min) + min
	return float64(randNum)
}
