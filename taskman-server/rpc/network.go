package rpc

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// HttpGet  Get请求
func HttpGet(url, userToken, language string) (byteArr []byte, err error) {
	req, newReqErr := http.NewRequest(http.MethodGet, url, strings.NewReader(""))
	if newReqErr != nil {
		err = fmt.Errorf("Try to new http request fail,%s ", newReqErr.Error())
		return
	}
	req.Header.Set("Authorization", userToken)
	req.Header.Set("Accept-Language", language)
	resp, respErr := http.DefaultClient.Do(req)
	if respErr != nil {
		err = fmt.Errorf("Try to do http request fail,%s ", respErr.Error())
		return
	}
	byteArr, _ = ioutil.ReadAll(resp.Body)
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
	byteArr, _ = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	return
}
