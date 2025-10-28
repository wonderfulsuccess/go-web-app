package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// MapPrint 打印 map 类型的数据
func MapPrint(mapData interface{}) {
	if mapData == nil {
		fmt.Println("null")
		return
	}
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(mapData); err != nil {
		fmt.Println("Error:", err)
		fmt.Printf("%#v\n", mapData)
		return
	}
	// json.Encoder adds a trailing newline; Print keeps it as-is
	fmt.Print(buf.String())
}
