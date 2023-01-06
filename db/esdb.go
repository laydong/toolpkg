package db

import (
	"github.com/laydong/toolpkg"
	"github.com/laydong/toolpkg/logx"
	"github.com/olivere/elastic/v6"
)

// InitEsClient ES初始化
func InitEsClient(addr, username, password string) (db *elastic.Client, err error) {
	// 创建client连接ES
	if username == "" || password == "" {
		db, err = elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(addr))
	} else {
		db, err = elastic.NewClient(
			elastic.SetSniff(false),
			// elasticsearch 服务地址，多个服务地址使用逗号分隔
			elastic.SetURL(addr),
			// 基于http base auth验证机制的账号和密码
			elastic.SetBasicAuth(username, password),
		)
	}
	if err != nil {
		logx.ErrorF(toolpkg.GetNewGinContext(), "es数据库链接错误", err.Error())
	}
	return
}
