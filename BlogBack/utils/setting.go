package utils

import (
	"fmt"
	"github.com/go-ini/ini"
)

var (
	AppMode    string
	HttpPort   string
	Db         string
	DbHost     string
	DbPort     string
	DbUser     string
	DbPassWord string
	DbName     string
)

func init() {
	file, err := ini.Load("config/config.ini")
	if err != nil {
		fmt.Println("配置文件读取错误，请检查文件路径：", err)
	}
	LoadServer(file)
	LoadData(file)
}

func LoadServer(file *ini.File) {
	if file == nil {
		fmt.Println("配置文件未加载成功")
		return
	}
	AppMode = file.Section("server").Key("AppMode").MustString("debug")
	HttpPort = file.Section("server").Key("HttpPort").MustString(":3000")
}

func LoadData(file *ini.File) {
	if file == nil {
		fmt.Println("配置文件未加载成功")
		return
	}
	Db = file.Section("database").Key("Db").MustString("debug")
	DbHost = file.Section("database").Key("DbHost").MustString("localhost")
	DbPort = file.Section("database").Key("DbPort").MustString("3306")
	DbUser = file.Section("database").Key("DbUser").MustString("svb")
	DbPassWord = file.Section("database").Key("DbPassWord").MustString("svb2861")
	DbName = file.Section("database").Key("DbName").MustString("base")
}
