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
	OutTo(subPath, file string, attr bool, keys *[]string, s *[]map[string]interface{}) (err error)
}

var ClientOut []Out
var ServerOut []Out

func Init() {
	// 加载客户端输出
	for _, t := range conf.Cfg.BuildClient {
		o, err := create(conf.Cfg.ClientPath, t)
		if err != nil {
			fmt.Println(err)
			continue
		}

		ClientOut = append(ClientOut, o)
	}

	// 加载服务器输出
	for _, t := range conf.Cfg.BuildServer {
		o, err := create(conf.Cfg.ServerPath, t)
		if err != nil {
			fmt.Println(err)
			continue
		}

		ServerOut = append(ServerOut, o)
	}

	// 加载post
	if len(conf.Cfg.PostUrl) > 0 {
		o, err := create("", conf.BuildToPost)
		if err != nil {
			fmt.Println(err)
		} else {
			ServerOut = append(ServerOut, o)
		}
	}
}

func Client(subPath, file string, attr bool, k *[]string, s *[]map[string]interface{}) (err error) {
	for _, o := range ClientOut {
		outerr := o.OutTo(subPath, file, attr, k, s)
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

func Server(subPath, file string, attr bool, k *[]string, s *[]map[string]interface{}) (err error) {
	for _, o := range ServerOut {
		outerr := o.OutTo(subPath, file, attr, k, s)
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
		o = new(Mdb)
	case conf.BuildToJson:
		o = new(Json)
	case conf.BuildToJs:
		o = new(Js)
	case conf.BuildToTs:
		o = new(Ts)
	case conf.BuildToCsv:
		o = new(Csv)
	case conf.BuildToPost:
		o = new(Post)
	default:
		err = fmt.Errorf("生成类型[%s]不存在", _type)
		return
	}

	err = o.Init(path)
	return
}
