package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

// EncodeURLParams 获取sign 加密  order = 1 正序 order =2 倒叙
func EncodeURLParams(bm map[string]interface{}, order int) string {
	var (
		buf  strings.Builder
		keys []string
	)
	for k := range bm {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for i, _ := range keys {
		num := i
		if order == 2 {
			num = len(keys) - i - 1
		}
		if v, ok := bm[keys[num]]; ok {
			buf.WriteString(url.QueryEscape(keys[num]))
			buf.WriteByte('=')
			buf.WriteString(url.QueryEscape(fmt.Sprintf("%v", v)))
			buf.WriteByte('&')
		}
	}
	if buf.Len() <= 0 {
		return ""
	}
	return MD5V([]byte(buf.String()[:buf.Len()-1]))
}

// SingCheng 获取sign加密对比  order = 1 正序 order =2 倒叙
func SingCheng(bm map[string]interface{}, sign string, order int) bool {
	var (
		buf  strings.Builder
		keys []string
	)
	for k := range bm {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for i, _ := range keys {
		num := i
		if order == 2 {
			num = len(keys) - i - 1
		}
		if v, ok := bm[keys[num]]; ok {
			buf.WriteString(url.QueryEscape(keys[num]))
			buf.WriteByte('=')
			buf.WriteString(url.QueryEscape(fmt.Sprintf("%v", v)))
			buf.WriteByte('&')
		}
	}
	if buf.Len() <= 0 {
		return false
	}
	if MD5V([]byte(buf.String()[:buf.Len()-1])) == sign {
		return true
	}
	return false
}

func MD5V(str []byte, b ...byte) string {
	h := md5.New()
	h.Write(str)
	return hex.EncodeToString(h.Sum(b))
}
