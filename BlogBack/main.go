package main

import (
	"BlogBack/api"
	"BlogBack/routes"
	"BlogBack/utils"
)

func main() {

	//引入数据路
	api.InitDb()
	//初始化
	r := routes.InitRouter()
	// 在main函数中添加
	r.Static("/Pictures", utils.GetImageBaseDir())


	r.SetTrustedProxies([]string{utils.DbHost})
	//启动服务器
	r.Run(utils.HttpPort)

}
