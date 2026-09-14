package utils

import (
	"path"
	"strings"
)

const ImagePublicPrefix = "/Pictures/"

// ResolveImageURL 把库里的文件名/路径统一成站点相对路径 /Pictures/...
// 已是 http(s) 或 /Pictures/ 则规范化后返回；空值用 fallback。
func ResolveImageURL(raw string, fallback string) string {
	s := strings.TrimSpace(raw)
	if s == "" {
		s = strings.TrimSpace(fallback)
	}
	if s == "" {
		return ""
	}

	if strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://") {
		if idx := strings.Index(s, "/Pictures/"); idx >= 0 {
			return normalizePicturesPath(s[idx+len("/Pictures/"):])
		}
		return s
	}

	s = strings.ReplaceAll(s, "\\", "/")
	s = strings.TrimPrefix(s, "/")
	if strings.HasPrefix(strings.ToLower(s), "pictures/") {
		s = s[len("pictures/"):]
	}
	return normalizePicturesPath(s)
}

func normalizePicturesPath(name string) string {
	name = strings.TrimPrefix(name, "/")
	name = path.Clean("/" + name)
	name = strings.TrimPrefix(name, "/")
	if name == "" || name == "." {
		return ""
	}
	return ImagePublicPrefix + name
}
