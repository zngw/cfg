package out

import (
	"encoding/json"
	"fmt"
	"github.com/zngw/cfg/conf"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Ts struct {
	Type string
	Path string
}

func (o *Ts) Init(path string) (err error) {
	o.Type = conf.BuildToTs
	if conf.Cfg.CreateTypePath {
		o.Path = filepath.Join(path, o.Type)
	}else {
		o.Path = path
	}
	err = os.MkdirAll(o.Path, os.ModePerm) //创建目录
	return
}

func (o *Ts) GetType() (t string) {
	return o.Type
}

// 转ts
func (o *Ts) OutTo(subPath, file string, attr bool, keys *[]string, s *[]map[string]interface{}) (err error) {
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

	path := o.Path
	if conf.Cfg.ReserveSubPath && len(subPath) > 0 {
		path = filepath.Join(o.Path, subPath)
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			fmt.Println("创建子目录错误", err)
			path = o.Path
		}
	}

	filename := filepath.Join(path, file+".ts")
	err = ioutil.WriteFile(filename, []byte("export let "+file+" = "+string(data)), os.ModePerm)
	if err != nil {
		err = fmt.Errorf("保存文件: %v 失败: %v ", filename, err)
		return
	}

	return
}
