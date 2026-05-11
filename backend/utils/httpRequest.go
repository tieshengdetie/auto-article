package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	GET int = iota
	PostForm
	PostReq
	PostJson
)

type CommonRespStruct struct {
	Code int
	Data interface{}
	Msg  string
}

func HandleRequest(requestType int, targetUrl string, param map[string]interface{}, header map[string]string, files map[string]io.Reader) (r *http.Response, err error) {
	client := &http.Client{}

	var req *http.Request

	switch requestType {
	case GET:
		req, err = http.NewRequest("GET", targetUrl, nil)
		if err != nil {
			return nil, err
		}
		q := req.URL.Query()
		for k, v := range param {
			q.Add(k, v.(string))
		}
		req.URL.RawQuery = q.Encode()
	case PostForm:
		var b bytes.Buffer
		w := multipart.NewWriter(&b)
		for key, r := range files {
			var fw io.Writer
			if x, ok := r.(io.Closer); ok {
				defer x.Close()
			}
			if x, ok := r.(*os.File); ok {
				if fw, err = w.CreateFormFile(key, x.Name()); err != nil {
					return nil, err
				}
			} else {
				if fw, err = w.CreateFormField(key); err != nil {
					return nil, err
				}
			}
			if _, err = io.Copy(fw, r); err != nil {
				return nil, err
			}
		}
		w.Close()
		req, err = http.NewRequest("POST", targetUrl, &b)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", w.FormDataContentType())
	case PostReq:
		form := url.Values{}
		for k, v := range param {
			form.Add(k, v.(string))
		}
		req, err = http.NewRequest("POST", targetUrl, strings.NewReader(form.Encode()))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	case PostJson:
		bytesData, err := json.Marshal(param)
		if err != nil {
			return nil, err
		}
		req, err = http.NewRequest("POST", targetUrl, bytes.NewBuffer(bytesData))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
	}

	for k, v := range header {
		req.Header.Set(k, v)
	}
	log.Printf("发送http请求,url: %s\n", targetUrl)
	strParam := MapToJson(param)
	log.Printf("请求参数 : %s\n", strParam)
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("请求出错 %s", err))
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("请求出错, 状态码 %d", resp.StatusCode))
	}
	return resp, nil
}

// HTTPRequest 发送HTTP请求
func HTTPRequest(ctx context.Context, requestType int, targetUrl string, param map[string]interface{}, header map[string]string, files map[string]io.Reader) (content []byte, err error) {
	var resp *http.Response
	resp, err = HandleRequest(requestType, targetUrl, param, header, files)

	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	var body []byte
	body, err = io.ReadAll(resp.Body)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("解析内容出错 %s", err))
	}
	log.Printf("返回结果: %s\n", string(body))

	return body, nil
}

func HTTPRequestByReceiver(requestType int, targetUrl string, param map[string]interface{}, header map[string]string, files map[string]io.Reader, receiver interface{}) (err error) {

	var resp *http.Response
	resp, err = HandleRequest(requestType, targetUrl, param, header, files)
	if err != nil {
		return err
	}

	return json.NewDecoder(resp.Body).Decode(&receiver)
}

func HTTPRequestByStream(requestType int, targetUrl string, param map[string]interface{}, header map[string]string, files map[string]io.Reader) (resp *http.Response, err error) {

	resp, err = HandleRequest(requestType, targetUrl, param, header, files)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
