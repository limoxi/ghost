package utils

import "encoding/json"

func Decode(jsonStr string, container interface{}){
	err := json.Unmarshal([]byte(jsonStr), container)
	if err != nil{
		panic(err)
	}
}