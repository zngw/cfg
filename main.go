package main

import (
	"fmt"
	"github.com/zngw/cfg/conf"
	"github.com/zngw/cfg/gen"
	"github.com/zngw/cfg/out"
	"github.com/zngw/cfg/util"
	"os"
)

func main() {
	// 读取配置
	err := conf.Init()
	if err != nil {
		fmt.Println("加载配置文件失败")
		util.WaitExit(1)
		return
	}

	// 目录清理
	_ = os.RemoveAll(conf.Cfg.ServerPath)
	_ = os.RemoveAll(conf.Cfg.ClientPath)

	// 初始化输出规则
	out.Init()

	// 输出文件
	files := conf.Cfg.SrcFiles
	if len(files) == 0 {
		// 不存在Excel文件，从目录中读取
		fmt.Println("正在生成" + conf.Cfg.SrcPath + "目录下的配置文件...")
		// 递归遍历xlsx文件
		files, err = util.GetFilesInDirs(conf.Cfg.SrcPath)
		if err != nil {
			fmt.Println("获取配置文件失败", err.Error())
			util.WaitExit(1)
			return
		}
	}

	// 遍历成生
	success := 0
	for _, file := range files {
		s, err := gen.Generate(file)
		if err != nil {
			fmt.Println(err.Error())
			util.WaitExit(1)
		}

		success += s
	}

	fmt.Println("=================================================================")
	fmt.Println("生成成功:", success)

	util.WaitExit(0)
}
