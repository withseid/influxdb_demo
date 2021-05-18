package utils

import (
	"math/rand"
	"time"
)

func GenFloatNum(min, max int) float32 {

	rand.Seed(time.Now().UnixNano())
	randNum := rand.Intn(max-min) + min
	return float32(randNum)
}
