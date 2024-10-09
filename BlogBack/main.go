package main

import (
	"BlogBack/api"
	"BlogBack/routes"
)

func main() {
	//引入数据路
	api.InitDb()
	//引入路由组件
	routes.InitRouter()

}
