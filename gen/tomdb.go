package gen

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// 转MongoDB使用的js脚本
func toMdb(path, file string, attr bool, keys *[]string, s *[]map[string]interface{}) (err error) {
	id := ""
	for _, K := range *keys {
		k := strings.ToLower(K)
		// 将id、Id、ID、已经identity视为mdb的_id值
		if k == "id" || k == "identity" {
			id = K
		}
	}

	if id != "" {
		for i, _ := range *s {
			id, ok := (*s)[i][id]
			if id != nil && ok {
				(*s)[i]["_id"] = strconv.Itoa(id.(int))
			}
		}
	}

	// 存本地js
	str := "db.getCollection(\"" + file + "\").drop();db.createCollection(\"" + file + "\");"
	for i, _ := range *s {
		str += "db.getCollection(\"" + file + "\").insert("
		str += toMdbJs((*s)[i])
		str += ");"
	}

	filename := path + file + ".js"
	err = ioutil.WriteFile(filename, []byte(str), os.ModePerm)
	if err != nil {
		fmt.Println("保存文件:", filename, "失败: ", err)
		return
	}

	return
}

func toMdbJs(value interface{}) (str string) {
	str = ""
	if valueMap, ok := value.(map[string]interface{}); ok {
		if len(valueMap) == 0 {
			str += "{}"
		} else {
			str += "{"
			for k, v := range valueMap {
				str += k
				str += ":"
				str += getMdbValue(v)
				str += ","
			}

			str = strings.TrimRight(str, ",")
			str += "}"
		}
	} else if valueSlice, ok := value.([]interface{}); ok {
		if len(valueSlice) == 0 {
			str += "[]"
		} else {
			str += "["
			for _, v := range valueSlice {
				str += getMdbValue(v)
				str += ","
			}
			str = strings.TrimRight(str, ",")
			str += "]"
		}
	}

	return
}

func getMdbValue(value interface{}) (str string) {
	switch value.(type) {
	case string:
		str = value.(string)
	case int:
		str = "NumberInt(\"" + strconv.Itoa(value.(int)) + "\")"
	case int64:
		str = "NumberLong(\"" + strconv.FormatInt(value.(int64), 10) + "\")"
	case float64:
		str = strconv.FormatFloat(value.(float64), 'f', -1, 64)
		break
	case bool:
		if value.(bool) {
			str = "true"
		} else {
			str = "false"
		}

	default:
		str = toMdbJs(value)
	}

	return
}
