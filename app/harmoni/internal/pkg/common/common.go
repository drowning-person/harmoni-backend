package common

import (
	"math/rand"
	"path/filepath"
	"strings"
	"time"
)

// RandStringRunes 返回随机字符串
func RandStringRunes(n int) string {
	var letterRunes = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// IsInExtensionList 返回文件的扩展名是否在给定的列表范围内
func IsInExtensionList(extList []string, fileName string) bool {
	ext := strings.ToLower(filepath.Ext(fileName))
	// 无扩展名时
	if len(ext) == 0 {
		return false
	}

	if ContainsString(extList, ext[1:]) {
		return true
	}

	return false
}

// ContainsString 返回list中是否包含
func ContainsString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
