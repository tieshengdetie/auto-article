package utils

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
	"unicode"

	"github.com/google/uuid"
)

// 定义一个类型约束，允许使用常见基本类型

type BasicType interface {
	~int | ~int64 | ~string | ~bool | ~float32 | ~float64
}

// IfGenerics 泛型函数
func IfGenerics[T BasicType](condition bool, trueVal, falseVal T) T {
	if condition {
		return trueVal
	}
	return falseVal
}

// If 模拟三元运算
func If(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}
func GenerateUniqueId() string {
	return uuid.NewString()
}
func camelToSnake(s string) string {
	var result []rune
	for i, r := range s {
		if unicode.IsUpper(r) && i > 0 {
			result = append(result, '_')
		}
		result = append(result, unicode.ToLower(r))
	}
	return string(result)
}

func IsEmptyMap(data interface{}) bool {
	if data == nil {
		return false
	}
	// 将 data 断言为 map[string]interface{}
	m, ok := data.(map[string]interface{})
	if !ok {
		return false
	}
	// 检查 map 是否为空
	return len(m) == 0
}

func PrintStructToJson(params interface{}) {
	var byteData []byte
	byteData, _ = json.Marshal(&params)
	fmt.Printf("%+v\n", string(byteData))
}

// Int64SliceToStringSlice 将 []int64 转换为 []string
func Int64SliceToStringSlice(slice []int64) []string {
	strSlice := make([]string, len(slice))
	for i, num := range slice {
		strSlice[i] = strconv.FormatInt(num, 10)
	}
	return strSlice
}

// StringToInt64Slice 将逗号分隔的字符串转换为 []int64
func StringToInt64Slice(s string) ([]int64, error) {
	if s == "" {
		return []int64{}, nil
	}
	strSlice := strings.Split(s, ",")
	int64Slice := make([]int64, len(strSlice))
	for i, str := range strSlice {
		num, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return nil, err
		}
		int64Slice[i] = num
	}
	return int64Slice, nil
}

// GenerateUniqueNumericID 生成不重复的数字字符串
func GenerateUniqueNumericID(length int) (string, error) {
	const digits = "0123456789"
	id := make([]byte, length)
	for i := range id {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		id[i] = digits[num.Int64()]
	}
	return string(id), nil
}
func Round(num float64, places int) float64 {
	factor := math.Pow(10, float64(places))
	return math.Round(num*factor) / factor
}

// Contains 判断切片中是否包含某一元素
func Contains[T comparable](slice []T, element T) bool {
	for _, v := range slice {
		if v == element {
			return true
		}
	}
	return false
}

// MergeSlices 合并多个切片
func MergeSlices[T any](slices ...[]T) []T {
	var mergedSlice []T
	for _, slice := range slices {
		mergedSlice = append(mergedSlice, slice...)
	}
	return mergedSlice
}

// SliceDifference 函数接受两个切片并返回第二个切片中不包含第一个切片中的元素
func SliceDifference[T comparable](slice1, slice2 []T) []T {
	// 创建一个映射，用于存储第二个切片中的元素
	elementMap := make(map[T]struct{})

	// 将第二个切片的元素存入映射
	for _, v := range slice2 {
		elementMap[v] = struct{}{}
	}

	var difference []T

	// 遍历第一个切片，检查每个元素是否在第二个切片中
	for _, v := range slice1 {
		if _, found := elementMap[v]; !found {
			difference = append(difference, v)
		}
	}

	return difference
}
