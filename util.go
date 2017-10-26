package tweethog

import (
	"math/rand"
	"time"
)

func GetRandomInt(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}
