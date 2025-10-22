package utils

import (
	"encoding/json"
	"fmt"
)

// MapPrint 打印 map 类型的数据
func MapPrint(mapData interface{}) {
	dataJson, err := json.MarshalIndent(mapData, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(dataJson))
}
