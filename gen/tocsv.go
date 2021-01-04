package gen

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

// 生成csv
func toCsv(path, file string, attr bool, keys *[]string, s *[]map[string]interface{}) (err error) {
	filename := path + file + ".csv"
	nfs, err := os.Create(filename)
	if err != nil {
		err = fmt.Errorf("保存文件: %v 失败: %v", filename, err)
		return
	}
	defer nfs.Close()
	_, _ = nfs.Seek(0, io.SeekEnd)

	w := csv.NewWriter(nfs)

	var newContent [][]string

	// 写入key
	newContent = append(newContent, *keys)

	// 写入值
	for i, _ := range *s {
		var line []string
		for _, key := range *keys {
			if v, ok := (*s)[i][key]; ok {
				line = append(line, getCsvValue(v))
			}
		}

		newContent = append(newContent, line)
	}

	err = w.WriteAll(newContent)
	return
}

func getCsvValue(value interface{}) (str string) {
	switch value.(type) {
	case string:
		str = value.(string)
	case int:
		str = strconv.Itoa(value.(int))
	case int64:
		str = strconv.FormatInt(value.(int64), 10)
	case float64:
		str = strconv.FormatFloat(value.(float64), 'f', -1, 64)
		break
	case bool:
		if value.(bool) {
			str = "TRUE"
		} else {
			str = "FALSE"
		}
	default:
		str = toCsvObj(value)
	}

	return
}

// csv中不存在嵌套情况，如果配置表中存在了嵌套时。
// 如果是map，按key排序，只取值，并以'='号分割, 程序读取时注意拆分和赋值。如： {"min":10,"max":100} => 100=10
// 如果是数组,用'|'分割。如： 【1，2，3，4】=> 1|2|3|4
// 如果是一层数组map。如 【{"min":10,"max":100，"weight":20},{"min":50,"max":200，"weight":55}】 => 100=10=20|200=50=55
// 如果是多层嵌套。劝诫你还是不要转csv了，转出来的数据会出错！！！
func toCsvObj(value interface{}) (str string) {
	str = ""
	if valueMap, ok := value.(map[string]interface{}); ok {
		if len(valueMap) == 0 {
			str += ""
		} else {
			var keys []string
			for k, _ := range valueMap {
				keys = append(keys, k)
			}
			sort.Strings(keys)

			for _, key := range keys {
				if v, ok := valueMap[key]; ok {
					str += getCsvValue(v)
					str += "="
				}
			}
			str = strings.TrimRight(str, "=")
		}
	} else if valueSlice, ok := value.([]interface{}); ok {
		if len(valueSlice) == 0 {
			str += ""
		} else {
			for _, v := range valueSlice {
				str += getCsvValue(v)
				str += "|"
			}
			str = strings.TrimRight(str, "|")
		}
	}

	return
}
