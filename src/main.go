package main

import (
	"fmt"
	"gin-admin/global"
	"gin-admin/initialize"
)

func main() {
	initialize.Init()
	fmt.Printf("服务初始化成功，数据库名称：%s", global.Config.DB.DBName)
}
