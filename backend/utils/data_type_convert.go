package utils

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

// StructToMap struct => map
func StructToMap(obj interface{}) map[string]interface{} {
	obj1 := reflect.TypeOf(obj)
	obj2 := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < obj1.NumField(); i++ {
		data[obj1.Field(i).Name] = obj2.Field(i).Interface()
	}
	return data
}

// StructToLittleHumpMap 利用反射将结构体转化为小驼峰map
func StructToLittleHumpMap(obj interface{}) map[string]interface{} {
	obj1 := reflect.TypeOf(obj)
	obj2 := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < obj1.NumField(); i++ {
		value := obj2.Field(i).Interface()
		if reflect.TypeOf(value).Kind() == reflect.Struct {
			data[LowerFirstLetter(obj1.Field(i).Name)] = StructToLittleHumpMap(value)
		} else if reflect.TypeOf(value).Kind() == reflect.Slice {
			var newValue []interface{}
			oldValue := reflect.ValueOf(value)
			for j := 0; j < oldValue.Len(); j++ {
				if reflect.TypeOf(oldValue.Index(j).Interface()).Kind() == reflect.Struct {
					newValue = append(newValue, StructToLittleHumpMap(oldValue.Index(j).Interface()))
				} else {
					newValue = append(newValue, oldValue.Index(j).Interface())
				}
			}
			data[LowerFirstLetter(obj1.Field(i).Name)] = newValue
		} else {
			data[LowerFirstLetter(obj1.Field(i).Name)] = value
		}
	}
	return data
}

// LowerFirstLetter 字符首字母小写
func LowerFirstLetter(str string) string {
	var upperStr string
	vv := []rune(str) // 后文有介绍
	for i := 0; i < len(vv); i++ {
		if i == 0 {
			if vv[i] >= 65 && vv[i] <= 90 { // 后文有介绍
				vv[i] += 32 // string的码表相差32位
				upperStr += string(vv[i])
			} else {
				fmt.Println("Not begins with lowercase letter,")
				return str
			}
		} else {
			upperStr += string(vv[i])
		}
	}
	return upperStr
}

// InterfaceToInt interface => int
func InterfaceToInt(source interface{}) (target int, err error) {
	if source != nil {
		switch reflect.TypeOf(source).Kind() {
		case reflect.Int:
			target = source.(int)
		case reflect.Int8:
			target = int(source.(int8))
		case reflect.Int16:
			target = int(source.(int16))
		case reflect.Int32:
			target = int(source.(int32))
		case reflect.Int64:
			target = int(source.(int64))
		case reflect.Float32:
			target = int(source.(float32))
		case reflect.Float64:
			target = int(source.(float64))
		case reflect.String:
			target, err = strconv.Atoi(source.(string))
		default:
			target = int(0)
			err = errors.New("输入类型异常！")
		}
	}
	return
}

// InterfaceToInt8 interface => int8
func InterfaceToInt8(source interface{}) (target int8, err error) {
	if source != nil {
		switch reflect.TypeOf(source).Kind() {
		case reflect.Int:
			target = int8(source.(int))
		case reflect.Int8:
			target = source.(int8)
		case reflect.Int16:
			target = int8(source.(int16))
		case reflect.Int32:
			target = int8(source.(int32))
		case reflect.Int64:
			target = int8(source.(int64))
		case reflect.Float32:
			target = int8(source.(float32))
		case reflect.Float64:
			target = int8(source.(float64))
		case reflect.String:
			var temp int
			temp, err = strconv.Atoi(source.(string))
			target = int8(temp)
		default:
			target = int8(0)
			err = errors.New("输入类型异常！")
		}
	}
	return
}

// InterfaceToInt16 interface => int16
func InterfaceToInt16(source interface{}) (target int16, err error) {
	if source != nil {
		switch reflect.TypeOf(source).Kind() {
		case reflect.Int:
			target = int16(source.(int))
		case reflect.Int8:
			target = int16(source.(int8))
		case reflect.Int16:
			target = source.(int16)
		case reflect.Int32:
			target = int16(source.(int32))
		case reflect.Int64:
			target = int16(source.(int64))
		case reflect.Float32:
			target = int16(source.(float32))
		case reflect.Float64:
			target = int16(source.(float64))
		case reflect.String:
			var temp int
			temp, err = strconv.Atoi(source.(string))
			target = int16(temp)
		default:
			target = int16(0)
			err = errors.New("输入类型异常！")
		}
	}
	return
}

// InterfaceToInt32 interface => int32
func InterfaceToInt32(source interface{}) (target int32, err error) {
	if source != nil {
		switch reflect.TypeOf(source).Kind() {
		case reflect.Int:
			target = int32(source.(int))
		case reflect.Int8:
			target = int32(source.(int8))
		case reflect.Int16:
			target = int32(source.(int16))
		case reflect.Int32:
			target = source.(int32)
		case reflect.Int64:
			target = int32(source.(int64))
		case reflect.Float32:
			target = int32(source.(float32))
		case reflect.Float64:
			target = int32(source.(float64))
		case reflect.String:
			var temp int64
			temp, err = strconv.ParseInt(source.(string), 10, 64)
			target = int32(temp)
		default:
			target = int32(0)
			err = errors.New("输入类型异常！")
		}
	}
	return
}

// InterfaceToInt64 interface => int64
func InterfaceToInt64(source interface{}) (target int64, err error) {
	if source != nil {
		switch reflect.TypeOf(source).Kind() {
		case reflect.Int:
			target = int64(source.(int))
		case reflect.Int8:
			target = int64(source.(int8))
		case reflect.Int16:
			target = int64(source.(int16))
		case reflect.Int32:
			target = int64(source.(int32))
		case reflect.Int64:
			target = source.(int64)
		case reflect.Float32:
			target = int64(source.(float32))
		case reflect.Float64:
			target = int64(source.(float64))
		case reflect.String:
			target, err = strconv.ParseInt(source.(string), 10, 64)
		default:
			target = int64(0)
			err = errors.New("输入类型异常！")
		}
	}
	return
}

// InterfaceToFloat32 interface => float32
func InterfaceToFloat32(source interface{}) (target float32, err error) {
	if source != nil {
		switch reflect.TypeOf(source).Kind() {
		case reflect.Int:
			target = float32(source.(int))
		case reflect.Int8:
			target = float32(source.(int8))
		case reflect.Int16:
			target = float32(source.(int16))
		case reflect.Int32:
			target = float32(source.(int32))
		case reflect.Int64:
			target = float32(source.(int64))
		case reflect.Float32:
			target = source.(float32)
		case reflect.Float64:
			target = float32(source.(float64))
		case reflect.String:
			var temp float64
			temp, err = strconv.ParseFloat(source.(string), 64)
			target = float32(temp)
		default:
			target = float32(0)
			err = errors.New("输入类型异常！")
		}
	}
	return
}

// InterfaceToFloat64 interface => float64
func InterfaceToFloat64(source interface{}) (target float64, err error) {
	if source != nil {
		switch reflect.TypeOf(source).Kind() {
		case reflect.Int:
			target = float64(source.(int))
		case reflect.Int8:
			target = float64(source.(int8))
		case reflect.Int16:
			target = float64(source.(int16))
		case reflect.Int32:
			target = float64(source.(int32))
		case reflect.Int64:
			target = float64(source.(int64))
		case reflect.Float32:
			target = float64(source.(float32))
		case reflect.Float64:
			target = source.(float64)
		case reflect.String:
			target, err = strconv.ParseFloat(source.(string), 64)
		default:
			target = float64(0)
			err = errors.New("输入类型异常！")
		}
	} else {
		target = float64(0)
	}
	return
}

// InterfaceToString interface => string
func InterfaceToString(source interface{}) (target string, err error) {
	if source != nil {
		switch reflect.TypeOf(source).Kind() {
		case reflect.Int:
			target = strconv.Itoa(source.(int))
		case reflect.Int8:
			target = strconv.Itoa(int(source.(int8)))
		case reflect.Int16:
			target = strconv.Itoa(int(source.(int16)))
		case reflect.Int32:
			target = strconv.FormatInt(int64(source.(int32)), 10)
		case reflect.Int64:
			target = strconv.FormatInt(source.(int64), 10)
		case reflect.Float32:
			target = strconv.FormatFloat(float64(source.(float32)), 'f', 6, 64)
		case reflect.Float64:
			target = strconv.FormatFloat(source.(float64), 'f', 6, 64)
		case reflect.String:
			target = source.(string)
		default:
			target = ""
			err = errors.New("输入类型异常！")
		}
	}
	return
}
