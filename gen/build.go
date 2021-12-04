package gen

import (
	"encoding/json"
	"fmt"
	"github.com/tealeg/xlsx/v3"
	"github.com/zngw/cfg/conf"
	"github.com/zngw/cfg/out"
	"github.com/zngw/cfg/util"
	"path/filepath"
	"strings"
)

// 单文件生成，传入文件路径和配置
func Generate(file string) (success int, err error) {
	xlf, err := xlsx.OpenFile(file)
	if err != nil {
		err = fmt.Errorf("读取文件: %v 失败: %v", file, err)
		return
	}

	subPath := ""
	fp, _ := filepath.Split(file)
	fp = util.GetAbsPath(fp)
	if len(conf.Cfg.SrcPath) > 0 {
		index := strings.Index(fp, conf.Cfg.SrcPath)
		if index == 0{
			subPath = fp[len(conf.Cfg.SrcPath)+1:]
		}
	}


	fmt.Println("-----------------------------------------------------------------")
	fmt.Println("读取配置文件:", file)
	pre := len(conf.Cfg.SheetPrefix)
	for _, sheet := range xlf.Sheets {
		// 判断表格生成前缀，如果配置为空则生成所有的表格
		if sheet.Name[0:pre] == conf.Cfg.SheetPrefix {
			if sheet.MaxRow == 0 || sheet.MaxCol == 0 {
				err = fmt.Errorf("生成失败: %v 数据为空！", sheet.Name)
				return
			}

			// 客户端数据
			var cs []map[string]interface{}
			// 客户端Key
			var ck []string

			// 服务器数据
			var ss []map[string]interface{}
			// 服务器Key
			var sk []string

			row, _ := sheet.Row(0)
			attr := row.GetCell(0).Value == ""
			if attr {
				// 解析属性表
				err = parseAttributes(sheet, &ck, &sk, &cs, &ss)
				if err != nil {
					err = fmt.Errorf("生成失败: %v, %v", sheet.Name, err)
					return
				}
			} else {
				// 解析数组表
				err = parseArray(sheet, &ck, &sk, &cs, &ss)
				if err != nil {
					err = fmt.Errorf("生成失败: %v, %v", sheet.Name, err)
					return
				}
			}

			// 生成客户端配置
			if ck!=nil && len(ck) > 0 {
				err = out.Client(subPath,sheet.Name, attr, &ck, &cs)
				if err != nil {
					err = fmt.Errorf("生成失败: %v, %v", sheet.Name, err)
					return
				}
			}


			// 生成服务器配置
			if sk!=nil && len(sk) > 0 {
				err = out.Server(subPath,sheet.Name, attr, &sk, &ss)
				if err != nil {
					err = fmt.Errorf("生成失败: %v, %v", sheet.Name, err)
					return
				}
			}

			fmt.Println("生成配置成功:", sheet.Name)
			success++
		}
	}

	return
}

// 解析数组表
func parseArray(sheet *xlsx.Sheet, ck, sk *[]string, cs, ss *[]map[string]interface{}) (err error) {

	rowNum := sheet.MaxRow //获取行数
	colNum := sheet.MaxCol //获取列数

	if rowNum < 4 || colNum < 1 {
		err = fmt.Errorf("文件格式错误:")
		return
	}

	// 变量所属服务器还是客户端
	adsorptions := make([]string, colNum)

	// 变量名
	value := make([]string, colNum)

	// 变量类型
	tp := make([]string, colNum)

	rowIndex := 0
	//遍历每一行
	err = sheet.ForEachRow(func(row *xlsx.Row) error {
		c := make(map[string]interface{})
		s := make(map[string]interface{})

		cellIndex := 0
		//遍历每一个单元
		err = row.ForEachCell(func(cell *xlsx.Cell) error {
			if cellIndex >= colNum {
				return nil
			}

			if rowIndex == 0 {
				// 第一行标明是服务器还是客户端
				if cell.String() == "" {
					colNum = cellIndex
					return nil
				}
				// 变量所属服务器还是客户端
				adsorptions[cellIndex] = cell.String()
			} else if rowIndex == 1 {
				// 第二行变量名
				value[cellIndex] = cell.String()
				ads := adsorptions[cellIndex]
				if ads == conf.Common || ads == conf.Client {
					*ck = append(*ck, cell.String())
				}

				if ads == conf.Common || ads == conf.Server {
					*sk = append(*sk, cell.String())
				}
			} else if rowIndex == 2 {
				// 第三行变量类型
				tp[cellIndex] = cell.String()
			} else if rowIndex == 3 {
				// 第四行说明，不处理
			} else {
				if cellIndex == 0 && cell.String() == "" {
					// 首行为空时结束
					return nil
				}

				key := value[cellIndex]
				ads := adsorptions[cellIndex]
				typ := tp[cellIndex]
				if ads != conf.Null {
					val, err := getValue(typ, cell)
					if err != nil {
						return fmt.Errorf("数据类型错误: key=%v, row=%v, cell=%v, %v", key, rowIndex, cell.String(), err)
					}

					if ads == conf.Common || ads == conf.Client {
						c[key] = val
					}

					if ads == conf.Common || ads == conf.Server {
						s[key] = val
					}
				}
			}

			cellIndex++
			return nil
		})

		if err != nil {
			return err
		}

		if len(c) > 0 {
			*cs = append(*cs, c)
		}

		if len(s) > 0 {
			*ss = append(*ss, s)
		}

		rowIndex++
		return nil
	})

	return
}

// 解析属性表
func parseAttributes(sheet *xlsx.Sheet, ck, sk *[]string, cs, ss *[]map[string]interface{}) (err error) {
	rowNum := sheet.MaxRow //获取行数
	colNum := sheet.MaxCol //获取列数

	if rowNum < 1 || colNum < 4 {
		err = fmt.Errorf("文件格式错误:")
		return
	}

	c := make(map[string]interface{})
	s := make(map[string]interface{})

	rowIndex := 0
	//遍历每一行
	err = sheet.ForEachRow(func(row *xlsx.Row) error {

		if rowIndex == 0 {
			rowIndex++
			return nil
		}
		rowIndex++

		ads := row.GetCell(0).String()
		key := row.GetCell(1).String()
		typ := row.GetCell(2).String()
		value := row.GetCell(3)

		if ads != conf.Null {
			val, err := getValue(typ, value)
			if err != nil {
				return fmt.Errorf("数据类型错误: key=%v, row=%v, cell=%v, %v", key, rowIndex, value.String(), err)
			}

			if ads == conf.Common || ads == conf.Client {
				*ck = append(*ck, key)
				c[key] = val
			}

			if ads == conf.Common || ads == conf.Server {
				*sk = append(*sk, key)
				s[key] = val
			}
		}

		return err
	})

	if len(c) > 0 {
		*cs = append(*cs, c)
	}

	if len(s) > 0 {
		*ss = append(*ss, s)
	}

	return
}

// 变量类型解析
func getValue(tp string, cell *xlsx.Cell) (val interface{}, err error) {
	switch tp {
	case "BOOL":
		val = cell.Bool()
	case "INT":
		val, err = cell.Int()
	case "LONG":
		val, err = cell.Int64()
	case "FLOAT":
		val, err = cell.Float()
	case "STRING":
		val = cell.String()
	case "OBJ":
		str := cell.String()
		if len(str) == 0 {
			val = make(map[string]interface{})
		} else {
			err = json.Unmarshal([]byte(str), &val)
		}
	case "ARRAY":
		str := cell.String()
		if len(str) == 0 {
			val = []string{}
		} else {
			err = json.Unmarshal([]byte(str), &val)
		}
	default:
		err = fmt.Errorf("%v  类型无法识别", tp)
		return
	}

	return
}
