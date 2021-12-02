// @Title
// @Description $
// @Author  55
// @Date  2021/12/3
package out

import (
	"fmt"
	"github.com/zngw/cfg/conf"
)

type Out interface {
	Init(path string) (err error)
	GetType() (t string)
	OutTo(file string, attr bool, keys *[]string, s *[]map[string]interface{}) (err error)
}

var ClientOut []Out
var ServerOut []Out

func Init(cfg *conf.Conf) {
	// 加载客户端输出
	for _, t := range cfg.BuildClient {
		o, err := create(cfg.ClientPath, t)
		if err != nil {
			fmt.Println(err)
			continue
		}

		ClientOut = append(ClientOut, o)
	}

	// 加载服务器输出
	for _, t := range cfg.BuildServer {
		o, err := create(cfg.ServerPath, t)
		if err != nil {
			fmt.Println(err)
			continue
		}

		ServerOut = append(ServerOut, o)
	}
}

func OutClient(file string, attr bool, k *[]string, s *[]map[string]interface{}) (err error) {
	for _, o := range ClientOut {
		outerr := o.OutTo(file, attr, k, s)
		if outerr != nil {
			if err == nil {
				err = fmt.Errorf("%s 生成 %s 失败", file, o.GetType())
			} else {
				err = fmt.Errorf("%s\n%s 生成 %s 失败", err.Error(), file, o.GetType())
			}
		}
	}
	return
}

func OutServer(file string, attr bool, k *[]string, s *[]map[string]interface{}) (err error) {
	for _, o := range ServerOut {
		outerr := o.OutTo(file, attr, k, s)
		if outerr != nil {
			if err == nil {
				err = fmt.Errorf("%s 生成 %s 失败", file, o.GetType())
			} else {
				err = fmt.Errorf("%s\n%s 生成 %s 失败", err.Error(), file, o.GetType())
			}
		}
	}
	return
}

func create(path, _type string) (o Out, err error) {
	switch _type {
	case conf.BuildToMdb:
		o = new(OutMdb)
	case conf.BuildToJson:
		o = new(OutJson)
	case conf.BuildToJs:
		o = new(OutJs)
	case conf.BuildToTs:
		o = new(OutTs)
	case conf.BuildToCsv:
		o = new(OutCsv)
	default:
		err = fmt.Errorf("生成类型[%s]不存在", _type)
		return
	}

	err = o.Init(path)
	return
}
