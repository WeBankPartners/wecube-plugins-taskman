package rpc

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// HttpGet  Get请求
func HttpGet(url, userToken, language string) (byteArr []byte, err error) {
	req, newReqErr := http.NewRequest(http.MethodGet, url, strings.NewReader(""))
	if newReqErr != nil {
		err = fmt.Errorf("try to new http request fail,%s ", newReqErr.Error())
		return
	}
	req.Header.Set("Authorization", userToken)
	req.Header.Set("Accept-Language", language)
	resp, respErr := http.DefaultClient.Do(req)
	if respErr != nil {
		err = fmt.Errorf("try to do http request fail,%s ", respErr.Error())
		return
	}
	byteArr, _ = io.ReadAll(resp.Body)
	defer resp.Body.Close()
	return
}

// HttpPost Post请求
func HttpPost(url, userToken, language string, postBytes []byte) (byteArr []byte, err error) {
	req, reqErr := http.NewRequest(http.MethodPost, url, bytes.NewReader(postBytes))
	if reqErr != nil {
		err = fmt.Errorf("new http reqeust fail,%s ", reqErr.Error())
		return
	}
	req.Header.Set("Authorization", userToken)
	req.Header.Set("Accept-Language", language)
	req.Header.Set("Content-Type", "application/json")
	resp, respErr := http.DefaultClient.Do(req)
	if respErr != nil {
		err = fmt.Errorf("do http reqeust fail,%s ", respErr.Error())
		return
	}
	byteArr, _ = io.ReadAll(resp.Body)
	defer resp.Body.Close()
	return
}

// HttpPut Put请求
func HttpPut(url, userToken, language string, postBytes []byte) (byteArr []byte, err error) {
	req, reqErr := http.NewRequest(http.MethodPut, url, bytes.NewReader(postBytes))
	if reqErr != nil {
		err = fmt.Errorf("new http reqeust fail,%s ", reqErr.Error())
		return
	}
	req.Header.Set("Authorization", userToken)
	req.Header.Set("Accept-Language", language)
	req.Header.Set("Content-Type", "application/json")
	resp, respErr := http.DefaultClient.Do(req)
	if respErr != nil {
		err = fmt.Errorf("do http reqeust fail,%s ", respErr.Error())
		return
	}
	byteArr, _ = io.ReadAll(resp.Body)
	defer resp.Body.Close()
	return
}

// HttpDelete Delete
func HttpDelete(url, userToken, language string) (byteArr []byte, err error) {
	req, reqErr := http.NewRequest(http.MethodDelete, url, strings.NewReader(""))
	if reqErr != nil {
		err = fmt.Errorf("new http reqeust fail,%s ", reqErr.Error())
		return
	}
	req.Header.Set("Authorization", userToken)
	req.Header.Set("Accept-Language", language)
	resp, respErr := http.DefaultClient.Do(req)
	if respErr != nil {
		err = fmt.Errorf("do http reqeust fail,%s ", respErr.Error())
		return
	}
	byteArr, _ = io.ReadAll(resp.Body)
	defer resp.Body.Close()
	return
}
