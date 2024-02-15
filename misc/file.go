package misc

import "os"

// resolve file path to absolute path
//
// eg:
//
//	~/config -> /home/<user>/config
//	config -> <workdir>/config
//	/var/log/a.log -> /var/log/a.log
func ResolveFilePath(original string) string {
	if original[0] == '/' {
		return original
	}
	if original[0] == '~' {
		home := os.Getenv("HOME")
		return home + original[1:]
	}

	wd, _ := os.Getwd()

	return wd + "/" + original
}
