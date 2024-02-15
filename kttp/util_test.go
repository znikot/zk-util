package kttp

import (
	"log"
	"regexp"
	"strings"
	"testing"

	"github.com/znikot/zk-util/misc"
)

func TestFillPathVaribles(t *testing.T) {
	fillPathVars("http://abc.com/:name/:id", PathVar{"name": "foo", "id": "123"})
	fillPathVars("http://abc.com/:name/:id/:action", PathVar{"name": "foo", "id": "123", "action": "del"})
	fillPathVars("//abc.com/:name/:id/:action", PathVar{"name": "foo", "id": "123", "action": "del"})
	fillPathVars("://abc.com/:name/:id/:action", PathVar{"name": "foo", "id": "123", "action": "del"})
	fillPathVars("://abc.com/:type/:name/:id/:action#hello", PathVar{"name": "foo", "id": "123"})
	fillPathVars("https://abc.com/:type/:name/:i_d/:action#hello", PathVar{"name": "foo", "i_d": "123"})
}

func fillPathVars(u string, v PathVar) {
	log.Printf("%s + %s -> %s", u, misc.ToJSON(v), FillPathVariables(u, v))
}

var urls = []string{
	"http://abc.com/:name/:id/:action",
	"http://abc.com/:name/:id/test",
	":id/name/:action",
}

func TestReg(t *testing.T) {
	reg := regexp.MustCompile(`(:\w+)`)

	for _, u := range urls {
		vars := reg.FindAllString(u, -1)
		// vars := reg.FindStringSubmatch(u)
		if len(vars) > 0 {
			log.Printf("match: %s -> %s", u, strings.Join(vars, " "))
		} else {
			log.Printf("no match: %s", u)
		}
	}
}
