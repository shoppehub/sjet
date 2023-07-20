package network

import (
	"encoding/json"
	"fmt"
	"github.com/CloudyKit/jet/v6"
	"github.com/shoppehub/sjet/util/common"
	"io/ioutil"
	"net"
	"net/http"
	"reflect"
	"time"
)

var client *http.Client

func init() {
	client = &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			IdleConnTimeout:     3 * time.Minute,
			TLSHandshakeTimeout: 5 * time.Second,
			DialContext: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 10 * time.Minute,
				DualStack: true,
			}).DialContext,
		},
	}
}

func PostJson() jet.Func {
	return func(a jet.Arguments) reflect.Value {
		urlStr := a.Get(0).String()
		paramMap := make(map[string]interface{})
		iter := a.Get(1).MapRange()
		for iter.Next() {
			k := iter.Key().String()
			v := fmt.Sprintf("%s", iter.Value())
			paramMap[k] = v
			//fmt.Printf("%s: %s\r\n", k, v)
		}
		resp, error := HttpPostJson(urlStr, paramMap)
		if error != nil {
			return reflect.ValueOf("出错了!")
		} else {
			response := string(resp)
			return reflect.ValueOf(response)
		}
	}
}

func PostWithHeader() jet.Func {
	return func(a jet.Arguments) reflect.Value {
		urlStr := a.Get(0).String()
		param := ""
		if a.IsSet(1) {
			body := a.Get(1).Interface()
			bodyStr, e := json.Marshal(body)
			if e != nil {
				panic(e)
			}
			param = string(bodyStr)
		}

		paramHeaderMap := make(map[string]string)
		if a.IsSet(2) {
			iterHeader := a.Get(2).MapRange()
			for iterHeader.Next() {
				k := iterHeader.Key().String()
				v := fmt.Sprintf("%s", iterHeader.Value())
				paramHeaderMap[k] = v
				//fmt.Printf("%s: %s\r\n", k, v)
			}
		}
		resp, error := HttpPostWithHeader(urlStr, param, paramHeaderMap)
		defer resp.Body.Close()
		if error != nil {
			return reflect.ValueOf("出错了!")
		} else {
			respBody, error := ioutil.ReadAll(resp.Body)
			if error != nil {

			}
			responseBody := make(map[string]interface{})
			common.ToObject(respBody, &responseBody)
			result := make(map[string]interface{})
			result["body"] = responseBody
			result["header"] = resp.Header

			return reflect.ValueOf(result)
		}
	}
}

func GetWithHeader() jet.Func {
	return func(a jet.Arguments) reflect.Value {
		urlStr := a.Get(0).String()
		paramMap := make(map[string]string)
		iter := a.Get(1).MapRange()
		for iter.Next() {
			k := iter.Key().String()
			v := fmt.Sprintf("%s", iter.Value())
			paramMap[k] = v
			//fmt.Printf("%s: %s\r\n", k, v)
		}
		resp, error := HttpGetWidthHeader(urlStr, paramMap)
		if error != nil {
			return reflect.ValueOf("出错了!")
		} else {
			response := string(resp)
			result := make(map[string]interface{})
			common.ToObject([]byte(response), &result)
			return reflect.ValueOf(result)
		}
	}
}
