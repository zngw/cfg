package gen

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// 转JS
func toJs(path, file string, attr bool, s *[]map[string]interface{}) (err error) {
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

	filename := path + file + ".js"
	err = ioutil.WriteFile(filename, []byte("var "+file+" = "+string(data)+"\nmodules.export = "+file), os.ModePerm)
	if err != nil {
		err = fmt.Errorf("保存文件: %v 失败: %v ", filename, err)
		return
	}

	return
}
