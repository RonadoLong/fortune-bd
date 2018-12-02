package helper

import "github.com/json-iterator/go"

func MarshalToByte(json interface{}) ([]byte, error){
	return jsoniter.Marshal(json)
}

func UnMarshal(source interface{}, data []byte) (error){
	return jsoniter.Unmarshal(data, &source)
}