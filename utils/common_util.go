package utils

import (
	"encoding/json"
	"runtime"
)

//Condition 模拟三元运算
func Condition(cond bool, a, b interface{}) (result interface{}) {
	if cond {
		result = a
	} else {
		result = b
	}
	return
}

func PrintObjToJsonStr(obj interface{}) string {
	j, _ := json.Marshal(obj)
	return string(j)
}

func MemStats() uint64 {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	return ms.Alloc
}
