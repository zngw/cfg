package out

import (
	"encoding/json"
	"fmt"
	"github.com/zngw/cfg/conf"
	"io/ioutil"
	"os"
	"path/filepath"
)

type OutTs struct {
	Type string
	Path string
}

func (o *OutTs) Init(path string) (err error) {
	o.Type = conf.BuildToTs
	o.Path = filepath.Join(path, o.Type)
	err = os.MkdirAll(o.Path, os.ModePerm) //创建目录
	return
}

func (o *OutTs) GetType() (t string) {
	return o.Type
}

// 转ts
func (o *OutTs) OutTo(file string, attr bool, keys *[]string, s *[]map[string]interface{}) (err error) {
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

	filename := filepath.Join(o.Path, file+".ts")
	err = ioutil.WriteFile(filename, []byte("export let "+file+" = "+string(data)), os.ModePerm)
	if err != nil {
		err = fmt.Errorf("保存文件: %v 失败: %v ", filename, err)
		return
	}

	return
}
