package gen

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

//转json
func toJson(path, file string, attr bool, s *[]map[string]interface{}) (err error) {
	var data []byte
	if attr {
		data, err = json.Marshal((*s)[0])

		if err != nil {
			return
		}
	} else {
		data, err = json.Marshal(*s)
		if err != nil {
			return
		}
	}

	filename := path + file + ".json"
	err = ioutil.WriteFile(filename, data, os.ModePerm)
	if err != nil {
		err = fmt.Errorf("保存文件: %v 失败: %v ", filename, err)
		return
	}

	return
}
