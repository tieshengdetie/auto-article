package sendHttp

type urlMap map[string]string

// 接口地址映射
var urlMapping = map[string]urlMap{}

// RegisterUrl 注册url映射
func registerUrl(key string, value urlMap) {
	if _, ok := urlMapping[key]; !ok {
		urlMapping[key] = value
	}

}

type CommonRespStruct struct {
	Code int
	Data interface{}
	Msg  string
}
