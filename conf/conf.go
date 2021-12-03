package conf

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/zngw/cfg/util"
	"io/ioutil"
	"strings"
)

// 配置表中表明该配置是属性哪端
const (
	Null   = "Null"   // 不加入配置
	Common = "Common" // 共同配置
	Client = "Client" // 客户端配置
	Server = "Server" // 服务器配置
)

// 输入文件类型
const (
	BuildToJson = "json" // 输出Json格式
	BuildToJs   = "js"   // 输出js格式
	BuildToTs   = "ts"   // 输出ts格式
	BuildToMdb  = "mdb"  // 输出Mongodb的js格式
	BuildToCsv  = "csv"  // 输出csv
	BuildToPost = "post" // 以json格式post到指定uri上
)

type Conf struct {
	SrcFiles    []string `json:"files,omitempty"`  // Excel文件,文件存在时，所在目录配置失效
	SrcPath     string   `json:"path,omitempty"`   // Excel所在目录
	SheetPrefix string   `json:"pre,omitempty"`    // 转换表前缀
	BuildClient []string `json:"client,omitempty"` // 客户端输出文件类型
	BuildServer []string `json:"server,omitempty"` // 服务器输出文件类型
	ServerPath  string   `json:"sPath,omitempty"`  // 服务器生成目录
	ClientPath  string   `json:"cPath,omitempty"`  // 客户端生成目录
	PostUrl     string   `json:"url,omitempty"`    // Post Json数据地址
	PostKey     string   `json:"key,omitempty"`    // Post Json验签密钥
}

var Cfg Conf

func getDefaultConf() Conf {
	return Conf{
		SrcFiles:    []string{},
		SrcPath:     "./excel",
		SheetPrefix: "Table",
		BuildClient: []string{},
		BuildServer: []string{},
		ServerPath:  "./out/server",
		ClientPath:  "./out/client",
	}
}

// 读取配置
func readConf(file string) (cfg Conf, err error) {
	cfg = getDefaultConf()
	// 没有配置文件
	if len(file) == 0 {
		return
	}

	// 读取配置文件
	b, err := ioutil.ReadFile(util.GetAbsPath(file))
	if err != nil {
		err = fmt.Errorf("读取配置文件 %v 失败: %v", file, err)
		return
	}

	// 解析配置文件
	err = json.Unmarshal(b, &cfg)
	if err != nil {
		err = fmt.Errorf("解析配置文件 %v 失败: %v", file, err)
		return
	}

	return
}

// 初始化配置
func Init() (err error) {
	// 接收命令行参数，命令行参数会替换文件中的配置
	// -c ./conf.json -path ./excel -pre Table -type all -client json|ts -server mdb
	cfg := flag.String("c", "", "传入配置文件")
	files := flag.String("files", "", "传入转换文件，如果文件存在刚后面的目录无效")
	src := flag.String("path", "", "传入转换文件所在目录")
	pre := flag.String("pre", "", "转换表前缀")
	cli := flag.String("client", "", "客户端输出文件类型")
	ser := flag.String("server", "", "服务器输出文件类型")
	sp := flag.String("sp", "", "服务器生成目录")
	cp := flag.String("cp", "", "客户端生成目录")
	url := flag.String("url", "", "Post Json数据地址")
	key := flag.String("key", "", "Post Json验签密钥")
	flag.Parse() //解析输入的参数

	Cfg, err = readConf(*cfg)
	if err != nil {
		fmt.Println("配置文件不存在，使用默认配置:", err.Error())
		Cfg = getDefaultConf()
	}

	if len(*src) > 0 {
		Cfg.SrcPath = *src
	}

	if len(*pre) > 0 {
		Cfg.SheetPrefix = *pre
	}

	if len(*cli) > 0 {
		Cfg.BuildClient = strings.Split(*cli, "|")
	}

	if len(*ser) > 0 {
		Cfg.BuildServer = strings.Split(*ser, "|")
	}

	if len(*files) > 0 {
		Cfg.SrcFiles = strings.Split(*files, ",")
	}
	
	if len(*sp) > 0 {
		Cfg.ServerPath = *sp
		Cfg.ServerPath = util.GetAbsPath(Cfg.ServerPath)
	}

	if len(*cp) > 0 {
		Cfg.ClientPath = *cp
		Cfg.ClientPath = util.GetAbsPath(Cfg.ClientPath)
	}

	if len(*url) > 0 {
		Cfg.PostUrl = *url
	}

	if len(*key) > 0{
		Cfg.PostKey = *key
	}	

	return
}
