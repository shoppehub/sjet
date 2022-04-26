package function

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/CloudyKit/jet/v6"
	"github.com/shoppehub/sjet/common"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func init() {
	globalFunc["string"] = stringFunc
	globalFunc["md5"] = md5Func
	globalFunc["base64"] = base64Func
	globalFunc["base64Decode"] = base64DecodeFunc

	globalFunc["substring"] = substringFunc

	globalFunc["urlEncode"] = urlEncodeFunc

	globalFunc["indexOf"] = indexOfFunc

	globalFunc["lenStr"] = lenStrFunc

	globalFunc["formatTime"] = formatTimeFunc

	globalFunc["writeJson"] = writeJsonFunc

	globalFunc["regexStringFormat"] = regexStringFormatFunc

	globalFunc["hasUnicodeHan"] = hasUnicodeHanFunc

	globalFunc["replaceAllRegex"] = replaceAllRegexFunc

	globalFunc["title"] = titleFunc

	globalFunc["toUpper"] = toUpperFunc

	globalFunc["toLower"] = toLowerFunc

	globalFunc["contains"] = containsFunc

	globalFunc["idEncode"] = idEncodeFunc
	globalFunc["idDecode"] = idDecodeFunc

	globalFunc["casIdEncode"] = casIdEncodeFunc
	globalFunc["casIdDecode"] = casIdDecodeFunc

}

func stringFunc(a jet.Arguments) reflect.Value {

	if !a.Get(0).IsValid() {
		return reflect.ValueOf("")
	}

	value := a.Get(0).Interface()
	kind := a.Get(0).Type().Kind()

	if kind == reflect.Float64 {
		num := value.(float64)
		return reflect.ValueOf(strconv.FormatFloat(num, 'f', -1, 64))
	}
	if kind == reflect.Bool {
		num := value.(bool)
		return reflect.ValueOf(strconv.FormatBool(num))
	}
	if kind == reflect.Int64 {
		num := value.(int64)
		return reflect.ValueOf(strconv.FormatInt(num, 10))
	}
	if kind == reflect.Int {
		num := value.(int)
		return reflect.ValueOf(fmt.Sprint(num))
	}

	name := a.Get(0).Type().Name()

	switch name {
	case "ObjectID":
		oid := a.Get(0).Interface().(primitive.ObjectID)
		return reflect.ValueOf(oid.Hex())
	}

	return reflect.ValueOf(a.Get(0).Interface())
}

func substringFunc(a jet.Arguments) reflect.Value {
	value := a.Get(0).Interface()

	if a.Get(0).Type().Kind() == reflect.Float64 {
		num := value.(float64)
		prefix := int32(a.Get(1).Interface().(float64))
		val, _ := decimal.NewFromFloat(num).Round(prefix).Float64()
		return reflect.ValueOf(val)
	} else {
		str := value.(string)

		strs := []rune(str)
		start := int32(0)
		end := int32(a.Get(1).Interface().(float64))
		if a.NumOfArguments() == 3 {
			start = int32(a.Get(1).Interface().(float64))
			end = int32(a.Get(2).Interface().(float64))
		}
		return reflect.ValueOf(string(strs[start:end]))
	}
}

func urlEncodeFunc(a jet.Arguments) reflect.Value {
	value := a.Get(0).String()
	value = url.PathEscape(value)
	return reflect.ValueOf(value)
}

func lenStrFunc(a jet.Arguments) reflect.Value {
	value := a.Get(0).Interface()
	strs := []rune(value.(string))
	return reflect.ValueOf(float64(len(strs)))
}

func indexOfFunc(a jet.Arguments) reflect.Value {
	value := a.Get(0).Interface().(string)
	key := a.Get(1).Interface().(string)
	return reflect.ValueOf(strings.Index(value, key))
}

func md5Func(a jet.Arguments) reflect.Value {
	str := a.Get(0).Interface().(string)
	h := md5.New()
	h.Write([]byte(str))
	val := hex.EncodeToString(h.Sum(nil))
	return reflect.ValueOf(val)
}

func base64Func(a jet.Arguments) reflect.Value {
	str := a.Get(0).Interface().(string)
	b := []byte(str)
	sEnc := base64.StdEncoding.EncodeToString(b)
	return reflect.ValueOf(sEnc)
}

func base64DecodeFunc(a jet.Arguments) reflect.Value {
	str := a.Get(0).Interface().(string)
	sDec, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		logrus.Error(err, str)
	}
	return reflect.ValueOf(sDec)
}

func formatTimeFunc(a jet.Arguments) reflect.Value {

	f := a.Get(1).Interface().(string)
	name := a.Get(0).Type().Name()

	if name == "DateTime" {
		str := a.Get(0).Interface().(primitive.DateTime)
		return reflect.ValueOf(str.Time().Format(f))
	}
	str := a.Get(0).Interface().(time.Time)
	return reflect.ValueOf(str.Format(f))
}

func writeJsonFunc(a jet.Arguments) reflect.Value {
	st, _ := json.Marshal(a.Get(0).Interface())
	return reflect.ValueOf(string(st))
}

func regexStringFormatFunc(a jet.Arguments) reflect.Value {
	str := a.Get(0).String()
	chars := AddEscapedChar(str)
	return reflect.ValueOf(chars)
}

func hasUnicodeHanFunc(a jet.Arguments) reflect.Value {
	str := a.Get(0).String()
	return reflect.ValueOf(common.HasUnicodeHan(str))
}

func AddEscapedChar(str string) string {
	var chars []rune
	for _, letter := range str {
		ok, letters := common.SpecialLetters(letter)
		if ok {
			chars = append(chars, letters...)
		} else {
			chars = append(chars, letter)
		}
	}
	return string(chars)
}

func replaceAllRegexFunc(a jet.Arguments) reflect.Value {
	src := a.Get(0).String()
	oldReg := a.Get(1).String()
	repl := a.Get(2).String()
	regExp := regexp.MustCompile(oldReg)
	result := regExp.ReplaceAllString(src, repl)
	return reflect.ValueOf(result)
}

func titleFunc(a jet.Arguments) reflect.Value {
	str := a.Get(0).String()
	return reflect.ValueOf(strings.Title(str))
}

func toUpperFunc(a jet.Arguments) reflect.Value {
	str := a.Get(0).String()
	return reflect.ValueOf(strings.ToUpper(str))
}

func toLowerFunc(a jet.Arguments) reflect.Value {
	str := a.Get(0).String()
	return reflect.ValueOf(strings.ToLower(str))
}

func containsFunc(a jet.Arguments) reflect.Value {
	str := a.Get(0).String()
	subStr := a.Get(1).String()
	return reflect.ValueOf(strings.Contains(str, subStr))
}

func idEncodeFunc(a jet.Arguments) reflect.Value {
	var casId int64
	k := a.Get(0).Kind()
	switch k {
	case reflect.Float64:
		casId = int64(a.Get(0).Float())
	default:
		casId = a.Get(0).Int()
	}
	encodedStr := common.EncodeId(casId)
	return reflect.ValueOf(encodedStr)
}

func idDecodeFunc(a jet.Arguments) reflect.Value {
	encodedStr := a.Get(0).String()
	id := common.DecodeId(encodedStr)
	return reflect.ValueOf(id)
}

func casIdEncodeFunc(a jet.Arguments) reflect.Value {
	var casId int64
	k := a.Get(0).Kind()
	switch k {
	case reflect.Float64:
		casId = int64(a.Get(0).Float())
	default:
		casId = a.Get(0).Int()
	}
	encodedStr := common.EncodeCasId(casId)
	return reflect.ValueOf(encodedStr)
}

func casIdDecodeFunc(a jet.Arguments) reflect.Value {
	encodedStr := a.Get(0).String()
	casId := common.DecodeCasId(encodedStr)
	return reflect.ValueOf(casId)
}
