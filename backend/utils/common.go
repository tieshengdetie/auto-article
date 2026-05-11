package utils

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// 宏定义
const (
	//Status
	FINISH   = 1 //完成
	UNFINISH = 2 //未完成
	FROZEN   = 3 //冻结
	ABANDON  = 4 //弃用
)

// 把Float32,Float64和int64类型统一为int
func IntId(val interface{}) (Id int) {
	if val != nil {
		switch reflect.TypeOf(val).Kind() {
		case reflect.Int64:
			Id = int(val.(int64))
		case reflect.Float32:
			Id = int(val.(float32))
		case reflect.Float64:
			Id = int(val.(float64))
		case reflect.Int32:
			Id = int(val.(int32))
		case reflect.Int:
			Id = val.(int)
		}
	}
	return Id
}

// 把获取状态名称
func StatusName(status int) (name string) {
	switch status {
	case FINISH:
		name = "正常"
	case UNFINISH:
		name = "未完成"
	case FROZEN:
		name = "冻结"
	case ABANDON:
		name = "弃用"
	default:
		name = "正常"
	}
	return name
}

// tiku状态
func TikuStatusName(status int) (name string) {
	switch status {
	case 0:
		name = "正常"
	case 1:
		name = "旧版本"
	case 2:
		name = "未完成"
	case 3:
		name = "冻结"
	case 4:
		name = "弃用"
	}
	return name
}

// 把获取状态名称
func StatusPptName(status int) (name string) {
	switch status {
	case FINISH:
		name = "初稿"
	case UNFINISH:
		name = "未完成"
	case FROZEN:
		name = "终稿"
	case ABANDON:
		name = "弃用"
	default:
		name = "正常"
	}
	return name
}

// Strval 获取变量的字符串值
// 浮点型 3.0将会转换成字符串3, "3"
// 非数值或字符类型的变量将会被转换成JSON格式字符串
// interface类型转string类型
func Strval(value interface{}) string {
	var key string
	if value == nil {
		return key
	}

	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}

	return key
}

func DeepCopy(value interface{}) interface{} {
	if valueMap, ok := value.(map[string]interface{}); ok {
		newMap := make(map[string]interface{})
		for k, v := range valueMap {
			newMap[k] = DeepCopy(v)
		}

		return newMap
	} else if valueSlice, ok := value.([]interface{}); ok {
		newSlice := make([]interface{}, len(valueSlice))
		for k, v := range valueSlice {
			newSlice[k] = DeepCopy(v)
		}

		return newSlice
	}

	return value
}

/**
 * @Author Dong
 * @Description 获得当前月的初始和结束日期
 * @Date 16:29 2020/8/6
 * @Param  * @param null
 * @return
 **/
func GetMonthDay() (string, string) {
	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
	f := firstOfMonth.Unix()
	l := lastOfMonth.Unix()
	return time.Unix(f, 0).Format("2006-01-02") + " 00:00:00", time.Unix(l, 0).Format("2006-01-02") + " 23:59:59"
}

/**
 * @Author Dong
 * @Description 获得当前周的初始和结束日期
 * @Date 16:32 2020/8/6
 * @Param  * @param null
 * @return
 **/
func GetWeekDay() (string, string) {
	now := time.Now()
	offset := int(time.Monday - now.Weekday())
	//周日做特殊判断 因为time.Monday = 0
	if offset > 0 {
		offset = -6
	}

	lastoffset := int(time.Saturday - now.Weekday())
	//周日做特殊判断 因为time.Monday = 0
	if lastoffset == 6 {
		lastoffset = -1
	}

	firstOfWeek := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset)
	lastOfWeeK := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, lastoffset+1)
	f := firstOfWeek.Unix()
	l := lastOfWeeK.Unix()
	return time.Unix(f, 0).Format("2006-01-02") + " 00:00:00", time.Unix(l, 0).Format("2006-01-02") + " 23:59:59"
}

/**
 * @Author Dong
 * @Description //获得当前季度的初始和结束日期
 * @Date 16:33 2020/8/6
 * @Param  * @param null
 * @return
 **/
func GetQuarterDay() (string, string) {
	year := time.Now().Format("2006")
	month := int(time.Now().Month())
	var firstOfQuarter string
	var lastOfQuarter string
	if month >= 1 && month <= 3 {
		//1月1号
		firstOfQuarter = year + "-01-01 00:00:00"
		lastOfQuarter = year + "-03-31 23:59:59"
	} else if month >= 4 && month <= 6 {
		firstOfQuarter = year + "-04-01 00:00:00"
		lastOfQuarter = year + "-06-30 23:59:59"
	} else if month >= 7 && month <= 9 {
		firstOfQuarter = year + "-07-01 00:00:00"
		lastOfQuarter = year + "-09-30 23:59:59"
	} else {
		firstOfQuarter = year + "-10-01 00:00:00"
		lastOfQuarter = year + "-12-31 23:59:59"
	}
	return firstOfQuarter, lastOfQuarter
}

// DistToOption 无序
func DistToOption(dict map[string]interface{}) (option []map[string]interface{}) {
	for k, v := range dict {
		optionItem := make(map[string]interface{})
		optionItem["value"] = k
		optionItem["label"] = v
		option = append(option, optionItem)
	}
	return
}

func Empty(params interface{}) bool {
	if params == nil {
		return true
	}
	//初始化变量
	var (
		flag          bool = true
		default_value reflect.Value
	)

	r := reflect.ValueOf(params)

	//获取对应类型默认值
	default_value = reflect.Zero(r.Type())
	//由于params 接口类型 所以default_value也要获取对应接口类型的值 如果获取不为接口类型 一直为返回false
	if !reflect.DeepEqual(r.Interface(), default_value.Interface()) {
		flag = false
	}
	return flag
}

func BytesCombine(pBytes ...[]byte) []byte {
	return bytes.Join(pBytes, []byte(""))
}

// bson二维结构排序,obj:要排序的数组,key：要排序的字段,sort:(desc:降序; asc:升序)
func BsonSort(obj []bson.M, key string, sort string) []bson.M {
	for i := len(obj) - 1; i >= 0; i-- {
		for i2 := len(obj) - 1; i2 >= 0; i2-- {
			if sort == "asc" {
				if IntId(obj[i][key]) > IntId(obj[i2][key]) {
					obj[i], obj[i2] = obj[i2], obj[i]
				}
			} else if sort == "desc" {
				if IntId(obj[i][key]) < IntId(obj[i2][key]) {
					obj[i], obj[i2] = obj[i2], obj[i]
				}
			}

		}
	}
	var newObj []bson.M
	for i := 0; i < len(obj); i++ {
		newObj = append(newObj, obj[i])
	}
	return newObj
}

// 把Float32,Float64和int64,string类型统一为int64
func Int64Id(val interface{}) (Id int64) {
	if val != nil {
		switch reflect.TypeOf(val).Kind() {
		case reflect.Int64:
			Id = val.(int64)
		case reflect.Float32:
			Id = int64(val.(float32))
		case reflect.Float64:
			Id = int64(val.(float64))
		case reflect.Int:
			Id = int64(val.(int))
		case reflect.Int32:
			Id = int64(val.(int32))
		case reflect.String:
			jsonId, err := val.(json.Number)
			if !err {
				Id, _ = strconv.ParseInt(val.(string), 10, 64)
			} else {
				Id, _ = jsonId.Int64()
			}
		}
	}
	return Id
}
