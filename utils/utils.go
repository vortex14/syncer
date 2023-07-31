package utils

import (
	"math/rand"
	"time"
)

func GetRandomIntRange(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min

}
