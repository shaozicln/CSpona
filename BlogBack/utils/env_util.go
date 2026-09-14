package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	appMode := strings.TrimSpace(cfg.Section("server").Key("AppMode").String())
	fmt.Printf("Config AppMode: %s\n", appMode)

	if strings.EqualFold(appMode, "release") {
		return ProdEnv
	}
	return DevEnv
}

// GetImageBaseDir 图片存储目录（磁盘路径，不是 URL）
// 优先级：config.ini [upload] Dir > 环境默认
// 生产默认：/www/server/go_project/cspona/dist/Pictures
// 开发默认：../BlogFont/public/Pictures
func GetImageBaseDir() string {
	dir := ""

	if cfg, err := ini.Load("config.ini"); err == nil {
		dir = strings.TrimSpace(cfg.Section("upload").Key("Dir").String())
	}

	if dir == "" {
		if GetEnv() == ProdEnv {
			dir = "/www/server/go_project/cspona/dist/Pictures"
		} else {
			devDir, _ := filepath.Abs(filepath.Join("..", "BlogFont", "public", "Pictures"))
			dir = devDir
		}
	} else if !filepath.IsAbs(dir) {
		abs, err := filepath.Abs(dir)
		if err == nil {
			dir = abs
		}
	}

	fmt.Println("Image storage directory:", dir)
	fmt.Printf("App env for images: %s\n", GetEnv())

	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("Failed to create directory: %v\n", err)
	}
	return dir
}

// GetAboutMePath 首页「关于我」Markdown 文件路径
// 优先级：config.ini [content] AboutMe > 默认 content/about-me.md
func GetAboutMePath() string {
	path := ""
	if cfg, err := ini.Load("config.ini"); err == nil {
		path = strings.TrimSpace(cfg.Section("content").Key("AboutMe").String())
	}
	if path == "" {
		path = filepath.Join("content", "about-me.md")
	}
	if !filepath.IsAbs(path) {
		if abs, err := filepath.Abs(path); err == nil {
			path = abs
		}
	}
	return path
}
