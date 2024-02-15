package misc

import (
	"fmt"
	"testing"
)

func TestResolveFilePath(t *testing.T) {

	// envs := os.Environ()
	// for _, env := range envs {
	// 	fmt.Println(env)
	// }

	p := "/user/kinz/test.txt"
	fmt.Printf("path %s: %s\n", p, ResolveFilePath(p))

	p = "user/kinz/test.txt"
	fmt.Printf("path %s: %s\n", p, ResolveFilePath(p))

	p = "~/user/kinz/test.txt"
	fmt.Printf("path %s: %s\n", p, ResolveFilePath(p))
}
