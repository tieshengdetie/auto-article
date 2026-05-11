package middleware

import (
	"AutoArticle/utils"
	"bytes"
	"fmt"
	"io"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
)

func DecryptParam() gin.HandlerFunc {

	return func(ctx *gin.Context) {

		param := ctx.GetHeader("request-param")
		aes, err := utils.DecryptAES([]byte(param), []byte("Y86yIcbyISTQGzPr3+hml5QaNmKYNibaVLkuH3aTbTU="))
		if err != nil {
			fmt.Println(err)
			ctx.Abort()
			return
		}
		fmt.Println(aes)
		switch ctx.Request.Method {
		case "POST":
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer([]byte(aes)))
		case "GET":
			var queryMap = map[string]interface{}{}
			originalURL := ctx.Request.URL
			//err = json.Unmarshal([]byte(aes), &queryMap)
			// int类型在json.Unmarshal时会被解析为float64，用interface接收时会丢失精度，所以这里使用json.Number
			d := json.NewDecoder(bytes.NewReader([]byte(aes)))
			d.UseNumber()
			err = d.Decode(&queryMap)
			var originQuery = originalURL.Query()
			if err != nil {
				fmt.Println(err)
			}
			for key, value := range queryMap {

				originQuery.Set(key, fmt.Sprintf("%v", value))
			}
			newURL := &url.URL{
				Path:     originalURL.Path,
				RawQuery: originQuery.Encode(),
			}
			ctx.Request.URL = newURL
		}

		ctx.Next()
	}

}
