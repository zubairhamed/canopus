package goap

import (
	"math/rand"
	"time"
)

var MESSAGEID_CURR = 0
func init() {
	rand.Seed(time.Now().UTC().UnixNano())

	MESSAGEID_CURR = rand.Intn(65535)
}
