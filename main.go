package main

import (
	"fmt"
	"github.com/zngw/cfg/conf"
	"github.com/zngw/cfg/gen"
	"github.com/zngw/cfg/util"
)

func main() {
	// 目录清理
	util.CleanPath()

	// 读取配置
	cft := conf.GetConf()

	fmt.Println("正在生成" + cft.SrcPath + "目录下的配置文件...")

	// 递归遍历xlsx文件
	files, err := util.GetFilesInDirs(cft.SrcPath)
	if err != nil {
		fmt.Println("获取配置文件失败", err.Error())
		util.WaitExit(1)
	}

	// 遍历成生
	success := 0
	for _, file := range files {
		s, err := gen.Generate(file, &cft)
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
