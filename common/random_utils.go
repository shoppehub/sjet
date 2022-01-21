package common

import (
	"bytes"
	"fmt"
	"math/rand"
)

func RandomInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func RandomString(l int) string {
	var result bytes.Buffer
	var temp string
	for i := 0; i < l; {
		temp = fmt.Sprint(RandomInt(65, 90))
		result.WriteString(temp)
		i++

	}
	return result.String()
}
