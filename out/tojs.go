package out

import (
	"encoding/json"
	"fmt"
	"github.com/zngw/cfg/conf"
	"io/ioutil"
	"os"
	"path/filepath"
)

type OutJs struct {
	Type string
	Path string
}

func (o *OutJs) Init(path string) (err error) {
	o.Type = conf.BuildToJs
	if conf.Cfg.CreateTypePath {
		o.Path = filepath.Join(path, o.Type)
	}else {
		o.Path = path
	}
	err = os.MkdirAll(o.Path, os.ModePerm) //创建目录
	return
}

func (o *OutJs) GetType() (t string) {
	return o.Type
}

// 转JS
func (o *OutJs) OutTo(file string, attr bool, keys *[]string, s *[]map[string]interface{}) (err error) {
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

	filename := filepath.Join(o.Path, file+".js")
	err = ioutil.WriteFile(filename, []byte("var "+file+" = "+string(data)+"\nmodules.export = "+file), os.ModePerm)
	if err != nil {
		err = fmt.Errorf("保存文件: %v 失败: %v ", filename, err)
		return
	}

	return
}
