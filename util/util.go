package utils

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"regexp"
	"strconv"
	"time"
)

var (
	chars = []string{
		"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
		"a", "b", "c", "d", "e", "f", "g", "h", "i", "j",
		"k", "l", "m", "n", "o", "p", "q", "r", "s", "t",
		"u", "v", "w", "x", "y", "z", "A", "B", "C", "D",
		"E", "F", "G", "H", "I", "J", "K", "L", "M", "N",
		"O", "P", "Q", "R", "S", "T", "U", "V", "W", "X",
		"Y", "Z", "~", "!", "@", "#", "$", "%", "^", "&",
		"*", "(", ")", "-", "_", "=", "+", "[", "]", "{",
		"}", "|", "<", ">", "?", "/", ".", ",", ";", ":"}
)

// GetComplexRandomString 获取复杂的随机n位字符
func GetComplexRandomString(n int) string {
	if n < 1 {
		return ""
	}
	var ret string
	for i := 0; i < n; i++ {
		ret += chars[rand.Intn(len(chars))]
	}
	return ret
}

// GetSimpleRandomString 获取纯英文大写的随机n位字符串
func GetSimpleRandomString(n int) string {
	ret := []byte{}
	for i := 0; i < n; i++ {
		tmp := rand.Intn(26)
		ret = append(ret, byte(tmp+65))
	}
	return string(ret)
}

// PhoneRegexp 验证是否手机
func PhoneRegexp(phone string) bool {
	b := false
	if phone != "" {
		reg := regexp.MustCompile(`^(86)*0*1\d{10}$`)
		b = reg.FindString(phone) != ""
	}
	return b
}

// StructCopy 对结构体中相同字段进行拷贝
func StructCopy(dst, src interface{}) {
	if dst == nil || src == nil {
		return
	}

	srcVal := reflect.Indirect(reflect.ValueOf(src))
	dstVal := reflect.Indirect(reflect.ValueOf(dst))
	if !(srcVal.Kind() == reflect.Struct && dstVal.Kind() == reflect.Struct) {
		return
	}

	srcType := srcVal.Type()
	for i := 0; i < srcVal.NumField(); i++ {
		srcField := srcVal.Field(i)
		name := srcType.Field(i).Name
		dstField := dstVal.FieldByName(name)
		if !dstField.CanSet() {
			continue
		}

		if getKind(srcField) != getKind(dstField) {
			continue
		}
		switch getKind(srcField) {
		case reflect.Int64:
			dstField.SetInt(srcField.Int())
		case reflect.Uint64:
			dstField.SetUint(srcField.Uint())
		case reflect.Float64:
			dstField.SetFloat(srcField.Float())
		case reflect.Bool, reflect.String, reflect.Slice:
			dstField.Set(srcField)
		case reflect.Ptr:
			StructCopy(dstField.Interface(), srcField.Interface())
		default:
		}
	}
}

func getKind(v reflect.Value) reflect.Kind {
	kind := v.Kind()
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return reflect.Int64
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return reflect.Uint64
	case reflect.Float32, reflect.Float64:
		return reflect.Float64
	default:
		return kind
	}
}

// FormatFloat 取小数点后n位非零小数
func FormatFloat(num float64, decimal int) string {
	if math.Trunc(num) == num || decimal == 0 {
		return fmt.Sprintf("%.f", math.Trunc(num))
	}
	format := "%." + strconv.Itoa(decimal) + "f"
	return fmt.Sprintf(format, num)
}

// RoundFloat 取小数点后n位非零小数,四舍五入
func RoundFloat(num float64, decimal int) string {
	for i := 0; i < decimal; i++ {
		num = num * 10
	}
	return fmt.Sprintf("%.1f", math.Round(num)/10)
}

// GetFirstDateOfMonth 获取本月第一天零点
func GetFirstDateOfMonth(d time.Time) time.Time {
	d = d.AddDate(0, 0, -d.Day()+1)
	return GetZeroTime(d)
}

// GetLastDateOfMonth 获取本月最后一天零点
func GetLastDateOfMonth(d time.Time) time.Time {
	return GetFirstDateOfMonth(d).AddDate(0, 1, -1)
}

// GetWeekZeroTime 获取本周第一天的0点时间
func GetWeekZeroTime(d time.Time) time.Time {
	wd := int(d.Weekday())
	if wd == 0 {
		wd = 7
	}
	todayZero := GetZeroTime(d)
	return todayZero.AddDate(0, 0, 1-wd)
}

// GetZeroTime 获取某一天的0点时间
func GetZeroTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

//GetString 从interface中获取string
func GetString(data interface{}) string {
	switch data.(type) {
	case string:
		return data.(string)
	default:
		return ""
	}
}

//GetInt 从interface中获取int
func GetInt(data interface{}) int {
	switch data.(type) {
	case int:
		return data.(int)
	case int64:
		return int(data.(int64))
	case float32:
		return int(data.(float32))
	case float64:
		return int(data.(float64))
	default:
		return 0
	}
}

//GetInt64 从interface中获取int64
func GetInt64(data interface{}) int64 {
	switch data.(type) {
	case int:
		return data.(int64)
	case int64:
		return data.(int64)
	case float32:
		return int64(data.(float32))
	case float64:
		return int64(data.(float64))
	default:
		return 0
	}
}

// VerifyCode 从interface中确认code是否为0
func VerifyCode(data interface{}) bool {
	switch data.(type) {
	case int:
		return data.(int) == 0
	case int64:
		return data.(int64) == 0
	case float32:
		return data.(float32) == 0
	case float64:
		return data.(float64) == 0
	case string:
		return data.(string) == "0"
	default:
		return false
	}
}

// GetDivideUpInt 除法，向上取整
func GetDivideUpInt(target, up, down int64) int64 {
	tmp := target * up
	if tmp%down != 0 {
		return tmp/down + 1
	}
	return tmp / down
}
