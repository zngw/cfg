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

// 编译输入类型
const (
	BuildTypeAll    = "all" // 编译服务和客户端
	BuildTypeClient = "cli" // 编译客户端
	BuildTypeServer = "ser" // 编译服务器
)

// 输入文件类型
const (
	BuildToJson = "json" // 输出Json格式
	BuildToJs   = "js"   // 输出js格式
	BuildToTs   = "ts"   // 输出ts格式
	BuildToMdb  = "mdb"  // 输出Mongodb的js格式
	BuildToCsv  = "csv"  // 输出csv
)

type Conf struct {
	SrcFiles    []string `json:"files,omitempty"`  // Excel文件,文件存在时，所在目录配置失效
	SrcPath     string   `json:"path,omitempty"`   // Excel所在目录
	SheetPrefix string   `json:"pre,omitempty"`    // 转换表前缀
	BuildType   string   `json:"type,omitempty"`   // 输出类型
	BuildClient []string `json:"client,omitempty"` // 客户端输出文件类型
	BuildServer []string `json:"server,omitempty"` // 服务器输出文件类型
	ServerPath  string   `json:"server,omitempty"` // 服务器生成目录
	ClientPath  string   `json:"server,omitempty"` // 客户端生成目录
}

func getDefaultConf() Conf {
	return Conf{
		SrcFiles:    []string{},
		SrcPath:     "./excel",
		SheetPrefix: "Table",
		BuildType:   BuildTypeAll,
		BuildClient: []string{BuildToJson},
		BuildServer: []string{BuildToMdb},
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

// 获取配置
func GetConf() (conf Conf, err error) {
	// 接收命令行参数，命令行参数会替换文件中的配置
	// -c ./conf.json -path ./excel -pre Table -type all -client json|ts -server mdb
	cfg := flag.String("c", "", "传入配置文件")
	files := flag.String("files", "", "传入转换文件，如果文件存在刚后面的目录无效")
	src := flag.String("path", "", "传入转换文件所在目录")
	pre := flag.String("pre", "", "转换表前缀")
	typ := flag.String("type", "", "传入转换类型")
	cli := flag.String("client", "", "客户端输出文件类型")
	ser := flag.String("server", "", "服务器输出文件类型")
	flag.Parse() //解析输入的参数

	conf, err = readConf(*cfg)
	if err != nil {
		fmt.Println("配置文件不存在，使用默认配置:", err.Error())
		conf = getDefaultConf()
	}

	if len(*src) > 0 {
		conf.SrcPath = *src
	}

	if len(*pre) > 0 {
		conf.SheetPrefix = *pre
	}

	if len(*typ) > 0 {
		conf.BuildType = *typ
	}

	if len(*cli) > 0 {
		conf.BuildClient = strings.Split(*cli, "|")
	}

	if len(*ser) > 0 {
		conf.BuildServer = strings.Split(*ser, "|")
	}

	if len(*files) > 0 {
		conf.SrcFiles = strings.Split(*files, ",")
	}

	conf.ServerPath = util.GetAbsPath(conf.ServerPath)
	conf.ClientPath = util.GetAbsPath(conf.ClientPath)

	return
}
