package helper

import "github.com/json-iterator/go"

func StructToJsonStr(query interface{}) string {
	bytes, _ := jsoniter.Marshal(query)
	return string(bytes)
}
