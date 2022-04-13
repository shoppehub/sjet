package common

import (
	"strconv"
	"strings"
	"unicode"
)

func HasUnicodeHan(str string) bool {
	var count int
	for _, v := range str {
		if unicode.Is(unicode.Han, v) {
			count++
			break
		}
	}
	return count > 0
}

// 如果存在特殊字符，直接在特殊字符前添加 \\ 转义
/**
判断是否为字母： unicode.IsLetter(v)
判断是否为十进制数字： unicode.IsDigit(v)
判断是否为数字： unicode.IsNumber(v)
判断是否为空白符号： unicode.IsSpace(v)
判断是否为Unicode标点字符 :unicode.IsPunct(v)
判断是否为中文：unicode.Han(v)
*/
func SpecialLetters(letter rune) (bool, []rune) {
	if unicode.IsSymbol(letter) || unicode.IsPunct(letter) {
		var chars []rune
		chars = append(chars, '\\', letter)
		return true, chars
	}
	return false, nil
}

// 对自增序列进行编码
func EncodeId(id int64) string {
	// 1.转为16进制
	idStr := strings.ToUpper(strconv.FormatInt(id, 16))
	// 2.前4位移动到结尾
	idStr = idStr[4:] + idStr[0:4]
	return idStr
}

//解码id
func DecodeId(idStr string) uint32 {
	idStr = idStr[len(idStr)-4:] + idStr[0:len(idStr)-4]
	n, err := strconv.ParseUint(idStr, 16, 32)
	if err != nil {
		panic("Parse Error")
	}
	id := uint32(n)
	return id
}

func EncodeCasId(casId int64) string {
	if casId < int64(10000000) {
		casId = casId + 100000000
	}
	encodedStr := strings.ToUpper(strconv.FormatInt(casId, 16))
	encodedStr = encodedStr[4:] + encodedStr[0:4]
	return encodedStr
}

func DecodeCasId(encodedStr string) uint32 {
	encodedStr = encodedStr[len(encodedStr)-4:] + encodedStr[0:len(encodedStr)-4]
	n, err := strconv.ParseUint(encodedStr, 16, 32)
	if err != nil {
		panic("Parse Error")
	}
	casId := uint32(n)
	if casId > uint32(100000000) {
		casId = casId - 100000000
	}
	return casId
}
