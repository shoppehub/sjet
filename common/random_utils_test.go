package common

import (
	"fmt"
	"testing"
)

func TestRandomInt(t *testing.T) {
	for i := 0; i < 100; i++ {
		result := RandomInt(0, 5)
		fmt.Printf("%v\n", result)
	}
}
