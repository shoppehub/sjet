package network

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

// 发送Post请求
func HttpPostWithHeader(url string, body string, headerMap map[string]string) (*http.Response, error) {
	//生成client 参数为默认
	client := &http.Client{}

	//提交请求
	request, err := http.NewRequest(http.MethodPost, url, strings.NewReader(body))
	if err != nil {
		panic(err)
	}
	request.Header.Set("Content-Type", "application/json")
	//增加header选项
	for key, value := range headerMap {
		request.Header.Set(key, value)
	}

	//处理返回结果
	return client.Do(request)
}

// 发送POST请求(JSON)
func HttpPostJson(url string, body interface{}) ([]byte, error) {
	postBodyStr := ""
	if reflect.TypeOf(body).Name() != "string" {
		bodyStr, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		postBodyStr = string(bodyStr)
	} else {
		postBodyStr = reflect.ValueOf(body).String()
	}

	return httpPost(url, "application/json", postBodyStr)
}

func HttpPostMap(urlStr string, contentType string, body map[string]string) ([]byte, error) {
	if contentType == "" {
		contentType = "application/x-www-form-urlencoded"
	}
	var data = make(url.Values)
	for key, value := range body {
		data.Set(key, value)
	}
	return httpPost(urlStr, contentType, data.Encode())
}

// 发送POST请求(XML)
func HttpPostXml(url string, xmlBody string) ([]byte, error) {
	return httpPost(url, "application/xml", xmlBody)
}

// 发送通用的POST请求
func httpPost(url string, contentType string, body string) ([]byte, error) {
	rsp, err := client.Post(url, contentType, strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	return ioutil.ReadAll(rsp.Body)
}
