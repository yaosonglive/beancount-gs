package service

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func FileMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	hash := md5.New()
	_, _ = io.Copy(hash, file)
	return hex.EncodeToString(hash.Sum(nil)), nil
}

// 获取所有图标
func GetIcons(c *gin.Context) {
	var resList = make([]string, 0)
	var fileMd5 = make(map[string]bool, 0)
	files, _ := filepath.Glob("./public/icons/*")
	for _, v := range files {
		tmp_md5, _ := FileMD5(v)
		if _, ok := fileMd5[tmp_md5]; ok {
			continue
		}
		fileMd5[tmp_md5] = true
		_, name := filepath.Split(v)
		accIcon := strings.Trim(name, filepath.Ext(name))
		if accIcon != "" {
			resList = append(resList, accIcon)
		}
	}
	OK(c, resList)
}
