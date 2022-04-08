package function

import (
	"github.com/shoppehub/sjet/common"
	"math"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/CloudyKit/jet/v6"
	"github.com/shoppehub/sjet/engine"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var globalFunc = make(map[string]jet.Func)

// 初始化全局函数
func InitGlobalFunc(t *engine.TemplateEngine) {

	for k, v := range globalFunc {
		t.Views.AddGlobalFunc(k, v)
	}

	// 把数字转换为int数组
	t.Views.AddGlobalFunc("numArray", numArrayFunc)
	// 支持把数据转换为字符串，比如 objectId
	t.Views.AddGlobalFunc("oid", oidFunc)
	t.Views.AddGlobalFunc("newObjectId", newObjectIdFunc)

	//function/strings.go:31 已完成{{ formatTime(info.createdAt,"2006-01-02")}}
	t.Views.AddGlobalFunc("time", timeFunc)
	t.Views.AddGlobalFunc("timeNowFormat", timeNowFormatFunc)
	t.Views.AddGlobalFunc("timeNowAddDateFormat", timeNowAddDateFormatFunc)
	t.Views.AddGlobalFunc("timeBefore", timeBeforeFunc)
	t.Views.AddGlobalFunc("timeOffset", timeOffsetFunc)

	t.Views.AddGlobalFunc("formatUrlPath", formatUrlPathFunc)

	t.Views.AddGlobalFunc("map", mapFunc)
	t.Views.AddGlobalFunc("deleteMapProperty", deleteMapPropertyFunc)
	t.Views.AddGlobalFunc("put", putFunc)
	t.Views.AddGlobalFunc("delete", deleteFunc)
	t.Views.AddGlobalFunc("append", appendFunc)

	t.Views.AddGlobalFunc("array", arrayFunc)
	t.Views.AddGlobalFunc("arraySort", arraySortFunc)
	t.Views.AddGlobalFunc("arrayAppend", arrayAppendFunc)

	t.Views.AddGlobalFunc("aggregate", aggregateFunc)
	t.Views.AddGlobalFunc("pipeline", aggregateFunc)

	t.Views.AddGlobalFunc("m", mFunc)
	t.Views.AddGlobalFunc("d", dFunc)

	t.Views.AddGlobalFunc("parseInt", parseIntFunc)
	t.Views.AddGlobalFunc("parseFloat", parseFloatFunc)
	t.Views.AddGlobalFunc("ceil", ceilFunc)
	t.Views.AddGlobalFunc("floor", floorFunc)
	t.Views.AddGlobalFunc("randomInt", randomIntFunc)

	t.Views.AddGlobalFunc("log", logFunc)

	t.Views.AddGlobalFunc("exit", exitFunc)
}

func oidFunc(a jet.Arguments) reflect.Value {
	if !a.Get(0).IsValid() {
		return reflect.ValueOf("")
	}
	oid, _ := primitive.ObjectIDFromHex(a.Get(0).String())
	return reflect.ValueOf(oid)
}
func newObjectIdFunc(a jet.Arguments) reflect.Value {
	oid := primitive.NewObjectID()
	return reflect.ValueOf(oid)
}
func timeFunc(a jet.Arguments) reflect.Value {
	value := a.Get(0)
	if !value.IsValid() {
		return reflect.ValueOf("")
	}
	str := value.String()
	layout := "2006-01-02 15:04:05"
	if strings.IndexAny(str, ":") == -1 {
		layout = "2006-01-02"
	}

	if strings.Contains(str, "T") {
		layout = time.RFC3339
	}

	val, _ := time.Parse(layout, str)
	//if err != nil {
	//	ejson, _ := json.Marshal(field)
	//	return nil, errors.New("value :" + value.(string) + " " + err.Error() + " ; the field config is " + string(ejson))
	//}
	return reflect.ValueOf(val)
}

func timeNowFormatFunc(a jet.Arguments) reflect.Value {
	layout := "2006-01-02"

	if a.IsSet(0) {
		format := a.Get(0)
		if !format.IsValid() {
			//layout = "2006-01-02 15:04:05"
			layout = time.RFC3339
		} else {
			layout = format.String()
		}
	}

	val := time.Now().Format(layout)
	return reflect.ValueOf(val)
}

func timeNowAddDateFormatFunc(a jet.Arguments) reflect.Value {
	originalTime := time.Now()

	days := 0
	if a.IsSet(0) {
		daysValue := a.Get(0)
		if daysValue.IsValid() && daysValue.Kind() == reflect.Float64 {
			days = int(daysValue.Float())
		} else {
			days = int(daysValue.Int())
		}
	}
	months := 0
	if a.IsSet(1) {
		monthsValue := a.Get(1)
		if monthsValue.IsValid() && monthsValue.Kind() == reflect.Float64 {
			months = int(monthsValue.Float())
		} else {
			months = int(monthsValue.Int())
		}
	}
	years := 0
	if a.IsSet(2) {
		yearsValue := a.Get(2)
		if yearsValue.IsValid() && yearsValue.Kind() == reflect.Float64 {
			years = int(yearsValue.Float())
		} else {
			years = int(yearsValue.Int())
		}
	}
	//layout = "2006-01-02 15:04:05"
	layout := time.RFC3339

	if a.IsSet(3) {
		format := a.Get(3)
		if format.IsValid() {
			layout = format.String()
		}
	}

	val := originalTime.AddDate(years, months, days).Format(layout)

	return reflect.ValueOf(val)
}

func timeBeforeFunc(a jet.Arguments) reflect.Value {
	time1 := time.Now()
	layout1 := "2006-01-02 15:04:05"
	layout2 := "2006-01-02 15:04:05"

	if a.IsSet(3) {
		// 如果传入四个参数，则格式 为 t1,l1, t2, l2
		format2 := a.Get(3)
		if format2.IsValid() {
			layout2 = format2.String()
		}
		format1 := a.Get(1)
		if format1.IsValid() {
			layout1 = format1.String()
		}
		time1, time1Error := time.Parse(layout1, a.Get(0).String())
		if time1Error != nil {
			return reflect.ValueOf(time1Error.Error())
		}
		time2, time2Error := time.Parse(layout2, a.Get(2).String())
		if time2Error != nil {
			return reflect.ValueOf(time2Error.Error())
		}
		return reflect.ValueOf(time1.Before(time2))
	} else if a.IsSet(1) { // 第四个参数没值，且第二个参数有值
		// 如果传入2个参数，则格式 为 t1, t2
		time1, time1Error := time.Parse(layout1, a.Get(0).String())
		if time1Error != nil {
			return reflect.ValueOf(time1Error.Error())
		}
		time2, time2Error := time.Parse(layout2, a.Get(1).String())
		if time2Error != nil {
			return reflect.ValueOf(time2Error.Error())
		}
		return reflect.ValueOf(time1.Before(time2))
	} else if a.IsSet(0) {
		// 如果传入1个参数，则为 t2, t1 默认为当前时间
		time2, time2Error := time.Parse(layout2, a.Get(1).String())
		if time2Error != nil {
			return reflect.ValueOf(time2Error.Error())
		}
		return reflect.ValueOf(time1.Before(time2))
	}

	return reflect.ValueOf("参数异常！")
}
func timeOffsetFunc(a jet.Arguments) reflect.Value {
	precision := "year"
	time1 := time.Now()
	time2 := time.Now()
	layout1 := "2006-01-02 15:04:05"
	layout2 := "2006-01-02 15:04:05"
	if a.NumOfArguments() > 1 {
		precision = a.Get(0).String()
	}
	if a.NumOfArguments() == 5 {
		layout1 = a.Get(2).String()
		layout2 = a.Get(4).String()
		time1, _ = time.Parse(layout1, a.Get(1).String())
		time2, _ = time.Parse(layout2, a.Get(3).String())
	} else if a.NumOfArguments() == 3 {
		layout2 = a.Get(2).String()
		time2, _ = time.Parse(layout2, a.Get(1).String())
	} else if a.NumOfArguments() == 1 {
		layout2 = "2006"
		time2, _ = time.Parse(layout2, a.Get(0).String())
	}
	switch precision {
	case "year":
		return reflect.ValueOf(time1.Year() - time2.Year())
		break
	case "day":
		return reflect.ValueOf(time1.Sub(time2).Hours() / 24)
		break
	case "hours":
		return reflect.ValueOf(time1.Sub(time2).Hours())
		break
	default:
		return reflect.ValueOf("参数异常！")
		break
	}

	return reflect.ValueOf("参数异常！")
}

func formatUrlPathFunc(a jet.Arguments) reflect.Value {
	if !a.Get(0).IsValid() {
		return reflect.ValueOf("")
	}
	u, _ := url.Parse(a.Get(0).Interface().(string))
	return reflect.ValueOf(u.Path)
}

// 把数字转换为int数组
func numArrayFunc(a jet.Arguments) reflect.Value {
	var total int
	k := a.Get(0).Kind()
	switch k {
	case reflect.Float64:
		total = int(a.Get(0).Float())
	default:
		total = int(a.Get(0).Int())
	}

	nums := make([]int64, total)
	for i := 0; i < total; i++ {
		nums[i] = int64(i + 1)
	}
	return reflect.ValueOf(nums)
}

func mapFunc(a jet.Arguments) reflect.Value {
	if a.NumOfArguments()%2 > 0 {
		return reflect.ValueOf(make(map[string]interface{}))
	}
	m := reflect.ValueOf(make(map[string]interface{}, a.NumOfArguments()/2))
	for i := 0; i < a.NumOfArguments(); i += 2 {

		m.SetMapIndex(a.Get(i), a.Get(i+1))
	}
	return m
}

func deleteMapPropertyFunc(a jet.Arguments) reflect.Value {
	if a.NumOfArguments() != 2 {
		return reflect.ValueOf(a.Get(0))
	}
	m := a.Get(0).Interface().(map[string]interface{})

	delete(m, a.Get(1).String())

	return reflect.ValueOf(m)
}

func deleteFunc(a jet.Arguments) reflect.Value {
	name := a.Get(0).Type().Name()

	if name == "M" {
		m := a.Get(0).Interface().(bson.M)
		m[a.Get(1).String()] = a.Get(2).Interface()
		return reflect.ValueOf(m)
	} else {
		m := a.Get(0).Interface().(map[string]interface{})
		delete(m, a.Get(1).String())
		return reflect.ValueOf(m)
	}
}

func putFunc(a jet.Arguments) reflect.Value {
	name := a.Get(0).Type().Name()

	if name == "M" {
		m := a.Get(0).Interface().(bson.M)
		m[a.Get(1).String()] = a.Get(2).Interface()
		return reflect.ValueOf(m)
	} else {
		m := a.Get(0).Interface().(map[string]interface{})
		m[a.Get(1).String()] = a.Get(2).Interface()
		return reflect.ValueOf(m)
	}
}

func appendFunc(a jet.Arguments) reflect.Value {
	name := a.Get(0).Type().Name()
	kind := a.Get(0).Type().Kind()

	if name == "D" {
		m := a.Get(0).Interface().(bson.D)
		e := bson.E{}
		e.Key = a.Get(1).Interface().(string)
		e.Value = a.Get(2).Interface()
		m = append(m, e)
		return reflect.ValueOf(m)
	} else if name == "M" {
		m := a.Get(0).Interface().(bson.M)
		if m[a.Get(1).String()] != nil {
			val := append(m[a.Get(1).String()].([]bson.M), a.Get(2).Interface().(bson.M))
			m[a.Get(1).String()] = val
		} else {
			val := []bson.M{a.Get(2).Interface().(bson.M)}
			m[a.Get(1).String()] = val
		}
		return reflect.ValueOf(m)
	} else if kind == reflect.Map {
		m := a.Get(0).Interface().(map[string]interface{})
		if m[a.Get(1).String()] != nil {
			val := append(m[a.Get(1).String()].([]interface{}), a.Get(2).Interface())
			m[a.Get(1).String()] = val
		} else {
			val := []interface{}{a.Get(2).Interface()}
			m[a.Get(1).String()] = val
		}
		return reflect.ValueOf(m)
	} else if kind == reflect.Slice {
		m := a.Get(0).Interface().([]interface{})
		m = append(m, a.Get(1).Interface())
		return reflect.ValueOf(m)
	}
	return reflect.ValueOf("")
}

func parseIntFunc(a jet.Arguments) reflect.Value {
	value := a.Get(0).Interface()
	val, _ := strconv.ParseInt(value.(string), 10, 64)
	return reflect.ValueOf(val)
}

func parseFloatFunc(a jet.Arguments) reflect.Value {
	value := a.Get(0).Interface()
	val, _ := strconv.ParseFloat(value.(string), 64)
	return reflect.ValueOf(val)
}

func ceilFunc(a jet.Arguments) reflect.Value {
	value := a.Get(0).Interface()
	return reflect.ValueOf(int(math.Ceil(value.(float64))))
}

func floorFunc(a jet.Arguments) reflect.Value {
	value := a.Get(0).Interface()
	return reflect.ValueOf(int(math.Floor(value.(float64))))
}

func randomIntFunc(a jet.Arguments) reflect.Value {
	if a.NumOfArguments() != 2 {
		return reflect.ValueOf("this func only support 2 params")
	}
	min := 0
	k := a.Get(0).Kind()
	switch k {
	case reflect.Float64:
		min = int(a.Get(0).Float())
	default:
		min = int(a.Get(0).Int())
	}

	max := 100
	k2 := a.Get(1).Kind()
	switch k2 {
	case reflect.Float64:
		max = int(a.Get(1).Float())
	default:
		max = int(a.Get(1).Int())
	}

	return reflect.ValueOf(common.RandomInt(min, max))
}

func mFunc(a jet.Arguments) reflect.Value {
	d := bson.M{}
	for i := 0; i < a.NumOfArguments(); i += 2 {
		d[a.Get(i).String()] = a.Get(i + 1).Interface()
	}
	m := reflect.ValueOf(d)
	return m
}

func dFunc(a jet.Arguments) reflect.Value {
	d := bson.D{}
	for i := 0; i < a.NumOfArguments(); i += 2 {
		d = append(d, bson.E{
			Key:   a.Get(i).String(),
			Value: a.Get(i + 1).Interface(),
		})
	}
	m := reflect.ValueOf(d)
	return m
}

func aggregateFunc(a jet.Arguments) reflect.Value {
	var p []bson.D
	for i := 0; i < a.NumOfArguments(); i++ {
		p = append(p, a.Get(i).Interface().(bson.D))
	}
	m := reflect.ValueOf(p)
	return m
}

func arrayFunc(a jet.Arguments) reflect.Value {
	var p []interface{}
	for i := 0; i < a.NumOfArguments(); i++ {
		p = append(p, a.Get(i).Interface())
	}
	m := reflect.ValueOf(p)
	return m
}
func arraySortFunc(a jet.Arguments) reflect.Value {
	paramSlice := a.Get(0).Slice(0, a.Get(0).Len())

	var result []map[string]interface{}
	for i := 0; i < paramSlice.Len(); i++ {
		item := paramSlice.Index(i)
		temp := item.Interface().(map[string]interface{})
		result = append(result, temp)
	}

	sort.Slice(result, func(i, j int) bool {
		var iSort, jSort int
		if _, ok := result[i]["sort"]; !ok {
			result[i]["sort"] = 99999
		}
		if _, ok := result[j]["sort"]; !ok {
			result[j]["sort"] = 99999
		}
		t := reflect.TypeOf(result[i]["sort"]).Kind()
		switch t {
		case reflect.Float64:
			iSort = int(result[i]["sort"].(float64))
			jSort = int(result[j]["sort"].(float64))
		default:
			iSort = result[i]["sort"].(int)
			jSort = result[j]["sort"].(int)
		}

		if iSort < jSort {
			return true
		}
		return false
	})
	return reflect.ValueOf(result)
}
func arrayAppendFunc(a jet.Arguments) reflect.Value {
	var p []interface{}
	for i := 0; i < a.NumOfArguments(); i++ {
		p = append(p, a.Get(i).Interface().([]interface{})...)
	}
	m := reflect.ValueOf(p)
	return m
}
func logFunc(a jet.Arguments) reflect.Value {

	level := a.Get(0).Interface().(string)
	logVal := a.Get(1).Interface()
	switch level {
	case "err":
		logrus.Error(logVal)
	case "info":
		logrus.Info(logVal)
	case "warn":
		logrus.Warn(logVal)
	}

	return reflect.ValueOf("")
}

func exitFunc(a jet.Arguments) reflect.Value {
	panic("exit::::")
}
