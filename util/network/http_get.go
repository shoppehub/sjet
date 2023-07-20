package network

import (
	"io/ioutil"
	"net/http"
)

// 发送GET请求
func HttpGet(url string) ([]byte, error) {
	rsp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	return ioutil.ReadAll(rsp.Body)
}

// 发送GET请求
func HttpGetWidthHeader(url string, headerMap map[string]string) ([]byte, error) {
	//生成client 参数为默认
	client := &http.Client{}
	//提交请求
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	//增加header选项
	for key, value := range headerMap {
		request.Header.Add(key, value)
	}

	//处理返回结果
	response, error := client.Do(request)
	if error != nil {
		return nil, error
	}
	defer response.Body.Close()
	return ioutil.ReadAll(response.Body)
}
