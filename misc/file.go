package misc

import (
	"os"
	"path"
	"strings"
)

// resolve file path to absolute path
//
// eg:
//
//	~/config -> /home/<user>/config
//	config -> <workdir>/config
//	/var/log/a.log -> /var/log/a.log
func ResolveFilePath(original string) string {
	original = strings.TrimSpace(original)
	ol := len(original)
	if ol == 0 {
		return original
	}

	// 当前程序运行的工作目录
	wd, _ := os.Getwd()

	if original[0] == '/' || (ol > 1 && original[1] == ':') {
		return original
	} else if ol > 1 && original[0:2] == "~/" {
		home := os.Getenv("HOME")
		return home + original[1:]
	} else if ol > 1 && original[0:2] == "./" {
		return wd + original[1:]
	} else if ol > 2 && original[0:3] == "../" {
		// 获取当前工作目录的上级路径
		wd = path.Dir(wd)
		return wd + original[2:]
	}

	// 默认是工作目录下的相对路径
	return wd + "/" + original
}
