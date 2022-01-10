package function

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestString(t *testing.T) {
	timestamp := time.Now().Unix()
	fmt.Println(timestamp)

	ti := primitive.DateTime(1626921098228)

	// tt := time.Unix(ti, 0)

	// primitive.DateTime.Time()

	fmt.Println(ti.Time().Format("2006-01-02"))

	//Golang 实现 float64 转 string
	var km = 9900.101
	str := fmt.Sprintf("%f", km)

	fmt.Println(str)

	strKm := strconv.FormatFloat(km, 'f', -1, 64)
	fmt.Println("StrKm = ", strKm)
}

func TestCasIdEncode(t *testing.T) {
	//casId := 10092462
	//casId := 100000001
	casId := int64(10092462)
	if casId < int64(10000000) {
		casId = casId + 100000000
	}
	encodedStr := strings.ToUpper(strconv.FormatInt(casId, 16))
	fmt.Println(strings.ToUpper(encodedStr))
	encodedStr = encodedStr[4:] + encodedStr[0:4]
	fmt.Printf("转16进制偏移后的casId：%v \n", encodedStr)

	encodedStr = encodedStr[len(encodedStr)-4:] + encodedStr[0:len(encodedStr)-4]
	fmt.Println(encodedStr)
	n, err := strconv.ParseUint(encodedStr, 16, 32)
	if err != nil {
		panic("Parse Error")
	}
	n2 := uint32(n)
	if n2 > uint32(100000000) {
		n2 = n2 - 100000000
	}
	fmt.Println(n2)
}
