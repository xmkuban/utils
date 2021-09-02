package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func KeyInSlice(key string, slice []string) (index int) {
	for i, v := range slice {
		if v == key {
			return i
		}
	}
	return -1
}

//MStrDdf 如果map里面的key不是目标key，则返回def，否则返回key对应的value
func MStrDdf(source map[string]interface{}, key, def string) string {
	_data, ok := source[key]
	if ok {
		res, ok1 := _data.(string)
		if ok1 {
			return res
		}
	}
	return def
}

func MInt64Def(source map[string]interface{}, key string, def int64) int64 {
	_data, ok := source[key]
	if ok {
		switch tmp := _data.(type) {
		case int64:
			return tmp
		case int:
			return int64(tmp)
		case float64:
			return int64(tmp)
		case string:
			res, _ := strconv.ParseInt(tmp, 10, 32)
			return res
		}
	}

	return def
}

func MIntDef(source map[string]interface{}, key string, def int) int {
	_data, ok := source[key]
	if ok {
		switch tmp := _data.(type) {
		case int:
			return tmp
		case int64:
			return int(tmp)
		case float64:
			return int(tmp)
		case string:
			res, _ := strconv.Atoi(tmp)
			return res
		}
	}

	return def
}

func MBoolDef(source map[string]interface{}, key string, def bool) bool {
	_data, ok := source[key]
	if ok {
		switch res := _data.(type) {
		case bool:
			return res
		case int:
			if res <= 0 {
				return false
			} else {
				return true
			}
		case int64:
			if res <= 0 {
				return false
			} else {
				return true
			}
		case float64:
			if res <= 0 {
				return false
			} else {
				return true
			}
		}
	}
	return def
}

func MMapStrInterfaceDef(source map[string]interface{}, key string, def map[string]interface{}) map[string]interface{} {
	_data, ok := source[key]
	if ok {
		switch res := _data.(type) {
		case map[string]interface{}:
			return res
		case map[string]map[string]interface{}:
			_res := make(map[string]interface{})
			for k, v := range res {
				_res[k] = v
			}
			return _res
		}
	}

	return def
}

func MMapStrStrDef(source map[string]interface{}, key string, def map[string]string) map[string]string {
	_data, ok := source[key]
	if !ok {
		return def
	}
	switch data := _data.(type) {
	case map[string]interface{}:
		res := make(map[string]string)
		for k, v := range data {
			switch vData := v.(type) {
			case string:
				res[k] = vData
			case []byte:
				res[k] = string(vData)
			default:
				continue
			}
		}
		return res
	case map[string]string:
		return data
	default:
		return def
	}
}

func MMapStrMapStrInterfaceDef(source map[string]interface{}, key string, def map[string]map[string]interface{}) map[string]map[string]interface{} {
	_data, ok := source[key]
	if !ok {
		return def
	}

	switch data := _data.(type) {
	case map[string]interface{}:
		res := make(map[string]map[string]interface{})
		for k, v := range data {
			switch vData := v.(type) {
			case map[string]interface{}:
				res[k] = vData
			default:
				continue
			}
		}
		return res
	case map[string]map[string]interface{}:
		return data
	default:
		return def
	}
}

func MSliceStrDef(source map[string]interface{}, key string, def []string) []string {
	_data, ok := source[key]
	if !ok {
		return def
	}

	switch data := _data.(type) {
	case []string:
		return data
	default:
		return def
	}
}

func MSliceInterfaceDef(source map[string]interface{}, key string, def []interface{}) []interface{} {
	_data, ok := source[key]
	if !ok {
		return def
	}

	switch data := _data.(type) {
	case []interface{}:
		return data
	case []string:
		var __data []interface{}
		for _, v := range data {
			__data = append(__data, v)
		}
		return __data
	case []int:
		var __data []interface{}
		for _, v := range data {
			__data = append(__data, v)
		}
		return __data
	case []map[string]interface{}:
		var __data []interface{}
		for _, v := range data {
			__data = append(__data, v)
		}
		return __data
	default:
		return def
	}
}

func CheckMKeys(source map[string]interface{}, keys []string) (err error) {
	var missParams []string
	for _, key := range keys {
		if _, ok := source[key]; ok {
			continue
		} else {
			missParams = append(missParams, key)
		}
	}
	if len(missParams) == 0 {
		return nil
	}
	missStr := strings.Join(missParams, ",")
	return fmt.Errorf("%s are/is missing", missStr)
}

//MMerge 2个map合并
func MMerge(base map[string]interface{}, args ...map[string]interface{}) map[string]interface{} {
	for _, arg := range args {
		for k, v := range arg {
			base[k] = v
		}
	}
	return base
}

func Slice2MapStrInterface(keys []string, arr []interface{}) map[string]interface{} {
	if len(keys) != len(arr) || len(keys) <= 0 || len(arr) <= 0 {
		return nil
	}
	res := make(map[string]interface{})
	for k, v := range arr {
		res[keys[k]] = v
	}
	return res
}

func Interface2DSlice2MapStrSlice(keys []string, arr [][]interface{}) []map[string]interface{} {
	res := make([]map[string]interface{}, 0)
	for _, v := range arr {
		res = append(res, Slice2MapStrInterface(keys, v))
	}
	return res
}

func Int2Bool(source int) bool {
	if source == 0 {
		return false
	}
	return true
}

func Bool2Int(source bool) int {
	if source {
		return 1
	}
	return 0
}
