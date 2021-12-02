package out

import (
	"encoding/json"
	"fmt"
	"github.com/zngw/cfg/conf"
	"io/ioutil"
	"os"
	"path/filepath"
)

type OutJson struct {
	Type string
	Path string
}

func (o *OutJson) Init(path string) (err error) {
	o.Type = conf.BuildToJson
	o.Path = filepath.Join(path, o.Type)
	err = os.MkdirAll(o.Path, os.ModePerm) //创建目录
	return
}

func (o *OutJson) GetType() (t string) {
	return o.Type
}

//转json
func (o *OutJson) OutTo(file string, attr bool, keys *[]string, s *[]map[string]interface{}) (err error) {
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

	filename := filepath.Join(o.Path, file+".json")
	err = ioutil.WriteFile(filename, data, os.ModePerm)
	if err != nil {
		err = fmt.Errorf("保存文件: %v 失败: %v ", filename, err)
		return
	}

	return
}
