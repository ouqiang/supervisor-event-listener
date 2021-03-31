package httpclient

// http-client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	client = &http.Client{Timeout: 6 * time.Second}
)

type ResponseWrapper struct {
	StatusCode int
	Body       string
	Header     http.Header
}

func (this *ResponseWrapper) IsOK() bool {
	return this != nil && this.StatusCode == 200
}

func (this *ResponseWrapper) String() string {
	if this == nil {
		return "null"
	}
	body := this.Body
	if runes := []rune(this.Body); len(runes) > 256 {
		body = string(runes[0:256])
	}
	return fmt.Sprintf("resp[%d] %s", this.StatusCode, body)
}

func (this *ResponseWrapper) Error() string {
	if this == nil {
		return "null"
	}
	if this.Body == "" {
		return "empty body"
	}
	return this.Body
}

func Get(url string, timeout int) *ResponseWrapper {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return createRequestError(err)
	}

	return request(req, timeout)
}

func PostParams(url string, params string, timeout int) *ResponseWrapper {
	buf := bytes.NewBufferString(params)
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return createRequestError(err)
	}
	req.Header.Set("Content-type", "application/x-www-form-urlencoded")

	return request(req, timeout)
}

func PostJson(url string, body string, timeout int) *ResponseWrapper {
	buf := bytes.NewBufferString(body)
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return createRequestError(err)
	}
	req.Header.Set("Content-type", "application/json")

	return request(req, timeout)
}

func request(req *http.Request, timeout int) *ResponseWrapper {
	wrapper := ResponseWrapper{StatusCode: 0, Body: "", Header: make(http.Header)}
	if timeout > 0 {
		client.Timeout = time.Duration(timeout) * time.Second
	}
	setRequestHeader(req)
	resp, err := client.Do(req)
	if err != nil {
		wrapper.Body = fmt.Sprintf("执行HTTP请求错误-%s", err.Error())
		return &wrapper
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		wrapper.Body = fmt.Sprintf("读取HTTP请求返回值失败-%s", err.Error())
		return &wrapper
	}
	wrapper.StatusCode = resp.StatusCode
	wrapper.Body = string(body)
	wrapper.Header = resp.Header
	return &wrapper
}

func setRequestHeader(req *http.Request) {
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.8,en;q=0.6")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2987.133 Safari/537.36 golang/gocron")
}

func createRequestError(err error) *ResponseWrapper {
	errorMessage := fmt.Sprintf("创建HTTP请求错误-%s", err.Error())
	return &ResponseWrapper{0, errorMessage, make(http.Header)}
}
