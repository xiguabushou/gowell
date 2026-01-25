package utils

import (
	"encoding/json"
	"strings"
)

func EncodeJson(input string)([]byte ,error){

    // 步骤1: 按中文和英文逗号分割
    // 这里使用 strings.FieldsFunc 可以更灵活地处理多种分隔符
    var separators = func(r rune) bool {
        return r == ',' || r == '，' // 英文逗号和中文逗号
    }
    parts := strings.FieldsFunc(input, separators)

    // 步骤2: 清理每个部分的空格（trim）
    var result []string
    for _, part := range parts {
        trimmed := strings.TrimSpace(part)
        if trimmed != "" { // 忽略空字符串
            result = append(result, trimmed)
        }
    }

    // 步骤3: 转换为 JSON
    jsonData, err := json.Marshal(result)
    if err != nil {
        return nil ,err
    }
	return  jsonData, nil
}

func UnencodeJson(jsonData []byte)([]string,error){
    var str []string
    err := json.Unmarshal(jsonData,&str)
    if err != nil {
        return nil,err  
    }
    return str,nil
}
