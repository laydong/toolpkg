package cloud_utlis

import (
	"github.com/olivere/elastic/v6"
	"log"
)

var Edb *elastic.Client

func InitEs(addr, dbname, username, password string) *elastic.Client {
	// 创建client连接ES
	if username == "" || password == "" {
		Edb, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(addr))
		if err != nil {
			log.Printf("[app.gstore] elastic error: %v", err.Error())
			panic(err)
		}
		return Edb
	} else {
		Edb, err := elastic.NewClient(
			elastic.SetSniff(false),
			// elasticsearch 服务地址，多个服务地址使用逗号分隔
			elastic.SetURL(addr),
			// 基于http base auth验证机制的账号和密码
			elastic.SetBasicAuth(username, password),
		)
		if err != nil {
			log.Printf("[app.gstore] elastic error: %v", err.Error())
			panic(err)
		}
		return Edb
	}

}
