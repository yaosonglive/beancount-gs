package service

import (
	"encoding/json"
	"github.com/beancount-gs/script"
	"github.com/gin-gonic/gin"
	"github.com/mozillazg/go-pinyin"
	"strconv"
	"strings"
	"unicode"
)

// 判断tag是否包含中文，是的话转换并保存
func getTag(c *gin.Context, str string) string {

	hasHan := false
	tag := str
	for _, r := range []rune(str) {
		if unicode.Is(unicode.Han, r) {
			hasHan = true
		}
	}
	if hasHan {
		a := pinyin.NewArgs()
		a.Fallback = func(r rune, a pinyin.Args) []string {
			return []string{string(r)}
		}
		pyList := pinyin.LazyConvert(str, &a)
		var slugList []string
		var tmpStr string = ""
		for i, r := range []rune(str) {
			if unicode.Is(unicode.Han, r) {
				if len(tmpStr) > 0 {
					slugList = append(slugList, tmpStr)
					tmpStr = ""
				}
				slugList = append(slugList, pyList[i])
			} else {
				tmpStr = tmpStr + string(r)
			}
		}
		if len(tmpStr) > 0 {
			slugList = append(slugList, tmpStr)
			tmpStr = ""
		}

		tag = "zh_" + strings.Join(slugList, "_")

		ledgerConfig := script.GetLedgerConfigFromContext(c)
		tagsMap := script.GetLedgerTags(ledgerConfig.Id)

		tmpTag := tag
		index := 1
		for {
			if name, ok := tagsMap[tmpTag]; ok {
				if name == str {
					return tmpTag
				} else {
					tmpTag = tag + "_" + strconv.Itoa(index)
					index = index + 1
				}
			} else {
				// 写入
				tagsMap[tmpTag] = str
				updateTagsMap(ledgerConfig, tagsMap)
				return tmpTag
			}
		}
	}

	return tag
}

func updateTagsMap(ledgerConfig *script.Config, tagsMap map[string]string) {
	// 更新文件
	pathFile := script.GetLedgerTagsMapFilePath(ledgerConfig.DataPath)
	bytes, err := json.MarshalIndent(tagsMap, "", "\t")
	if err != nil {
		return
	}
	err = script.WriteFile(pathFile, string(bytes))
	if err != nil {
		return
	}
	// 更新缓存
	script.UpdateLedgerTags(ledgerConfig.Id, tagsMap)
}
