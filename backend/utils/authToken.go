// Package utils
// @Description: 中台2.0生成token方法
package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"
)

// GetToken 生成认证字符串
// @param	appKey	      string	 "中台下发的AppKey"
// @param	appSecret     string	 "中台下发的AppSecret"
// @param	host	      string	 "请求host"
// @param	path	      string	 "请求的接口路由path"
// @param	method	      string	 "请求method，Get/Post等"
// @param       validTime     int        "签名有效期（秒）"
// @return                    string 	 "最终生成的认证字符串"
func GetToken(appKey, appSecret, host, path, method string, validTime int) string {
	header := make(map[string]string)
	header["Host"] = host
	canonicalRequest := GenCanonicalRequest(Request{
		Method: method,
		Path:   path,
		Header: header,
	})
	authStringPrefix := GenAuthStringPrefix(appKey, validTime)
	token := GenSignature(appSecret, authStringPrefix, canonicalRequest)

	return authStringPrefix + "/host/" + token
}

// Request 请求参数
type Request struct {
	Method string
	Path   string
	Query  map[string]string
	Header map[string]string
}

// URLEncode 对绝对路径进行编码
// @param	str	string	"绝对路径"
// @return      string 	"编码结果"
func URLEncode(str string) string {
	return strings.Replace(url.QueryEscape(str), "+", "%20", -1)
}

// URIEncodeExceptSlash 对绝对路径进行编码，但是不对"/"进行编码
// @param	uri string	"绝对路径"
// @return          string 	"编码结果"
func URIEncodeExceptSlash(uri string) string {
	var result string
	for _, char := range uri {
		str := fmt.Sprintf("%c", char)
		if str == "/" {
			result += str
		} else {
			result += URLEncode(str)
		}
	}

	return result
}

// ToCanonicalQueryString 对URL中的Query String进行编码
// @param	params	map	"Query String组成的map"
// @return              string 	"编码结果"
func ToCanonicalQueryString(params map[string]string) string {
	if params == nil {
		return ""
	}

	encodedQueryStrings := make([]string, 0, 10)
	var query string

	for key, value := range params {
		if key != "" {
			query = URLEncode(key) + "="
			if value != "" {
				query += URLEncode(value)
			}
			encodedQueryStrings = append(encodedQueryStrings, query)
		}
	}

	sort.Strings(encodedQueryStrings)

	return strings.Join(encodedQueryStrings, "&")
}

// ToCanonicalHeaderString 对请求中的Header进行选择性编码
// @param	headerMap map	"选择的Header组成的map"
// @return          	string 	"编码结果"
func ToCanonicalHeaderString(headerMap map[string]string) string {
	headers := make([]string, 0, len(headerMap))
	for key, value := range headerMap {
		headers = append(headers,
			fmt.Sprintf("%s:%s", URLEncode(strings.ToLower(key)),
				URLEncode(strings.TrimSpace(value))))
	}

	sort.Strings(headers)

	return strings.Join(headers, "\n")
}

// GenCanonicalRequest 对http请求进行规范化编码
// @param	dto	Request	"http请求信息"
// @return              string 	"编码结果"
func GenCanonicalRequest(req Request) string {
	md := strings.ToUpper(req.Method)
	curi := URIEncodeExceptSlash(req.Path)
	cqs := ToCanonicalQueryString(req.Query)
	chs := ToCanonicalHeaderString(req.Header)
	return md + "\n" + curi + "\n" + cqs + "\n" + chs
}

// GenAuthStringPrefix 生成前缀字符串
// @param	accessKeyId		string	"访问密钥ID"
// @param	expirationPeriodInSecs	int 	"签名有效期（秒）"
// @return          			string 	"前缀字符串"
func GenAuthStringPrefix(accessKeyId string, expirationPeriodInSecs int) string {
	ts := time.Now().UTC().Format(time.RFC3339)
	authStringPrefix := fmt.Sprintf("readboy-edu-middle-v2/%s/%s/%d", accessKeyId, ts, expirationPeriodInSecs)
	return authStringPrefix
}

// GenSignature 生成签名摘要
// @param	secretAccessKey		string	"秘密访问密钥"
// @param	authStringPrefix	string	"前缀字符串"
// @param	canonicalRequest	string	"规范化请求"
// @return            			string 	"签名摘要"
func GenSignature(secretAccessKey, authStringPrefix, canonicalRequest string) string {
	h1 := hmac.New(sha256.New, []byte(secretAccessKey))
	h1.Write([]byte(authStringPrefix))
	signingKey := hex.EncodeToString(h1.Sum(nil))

	h2 := hmac.New(sha256.New, []byte(signingKey))
	h2.Write([]byte(canonicalRequest))
	signature := hex.EncodeToString(h2.Sum(nil))
	return signature
}
