package utils

import (
	"fmt"
	"path/filepath"
	"os"
	"github.com/go-ini/ini"
)

const (
	ProdEnv = "production"
	DevEnv  = "development"
)

func GetEnv() string {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return DevEnv
	}
	appMode := cfg.Section("server").Key("AppMode").String()
	fmt.Printf("Config AppMode: %s\n", appMode)

	// 严格匹配，忽略空格和大小写
	if appMode == "release" {
		return ProdEnv
	}
	return DevEnv
}

func GetImageBaseDir() string {
	var dir string
	if GetEnv() == ProdEnv {
		dir = "/www/server/go_project/cspona/dist/Pictures"
	} else {
		devDir, _ := filepath.Abs(filepath.Join("..", "BlogFont", "public", "Pictures"))
		dir = devDir
	}
	fmt.Println("Image storage directory:", dir)

	// 确保目录存在
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("Failed to create directory: %v\n", err)
	}
	return dir
}
