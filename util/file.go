package util

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

// 将路径转成绝对路径
func GetAbsPath(path string) (out string) {
	dir := GetCurrentDirectory()
	if !filepath.IsAbs(path) {
		out = filepath.Join(dir, path)
	} else {
		out = path
	}
	return
}

// 获取当前程序文件所在目录
func GetCurrentDirectory() string {
	//返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	return strings.Replace(dir, "\\", "/", -1) //将\替换成/
}

// 递归遍历目录下所有的xlsx文件，返回带路径的文件地址
func GetFilesInDirs(dirPth string) (files []string, err error) {
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return
	}

	PthSep := string(os.PathSeparator)

	for _, fi := range dir {
		if fi.IsDir() { // 目录, 递归遍历
			f, e := GetFilesInDirs(dirPth + PthSep + fi.Name())
			if e == nil {
				files = append(files, f...)
			}
		} else {
			// 过滤指定格式
			if checkIsHidden(fi) {
				continue
			}

			ok := strings.HasSuffix(fi.Name(), ".xlsx")
			if ok {
				files = append(files, dirPth+PthSep+fi.Name())
			}
		}
	}

	return
}

// 过滤隐藏属性文件
func checkIsHidden(file os.FileInfo) bool {
	if runtime.GOOS == "windows" {
		//"通过反射来获取Win32FileAttributeData的FileAttributes
		fa := reflect.ValueOf(file.Sys()).Elem().FieldByName("FileAttributes").Uint()
		bytefa := []byte(strconv.FormatUint(fa, 2))
		if bytefa[len(bytefa)-2] == '1' {
			return true
		}
		return false
	} else {
		if len(file.Name()) > 0 {
			if file.Name()[0] == '.' {
				return true
			}
		}
	}

	return false
}
