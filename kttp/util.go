package kttp

import (
	"net/http"
	"regexp"
	"strings"
)

// extract file name from header Content-Disposition
func ExtractFileName(header http.Header) string {
	contentDispositions := header.Values("Content-Disposition")
	if len(contentDispositions) == 0 {
		return ""
	}
	for _, cd := range contentDispositions {
		idx := strings.Index(strings.ToLower(cd), "filename=")
		if idx != -1 {
			return cd[idx+9:]
		}
	}
	return ""
}

// type for path variables
type PathVar map[string]string

var pathVarReg = regexp.MustCompile(`:\w+`)

// fil path varibles
//
//	FillPathVariables("http://localhost:8080/api/user/:id", PathVar{"id": "123"}) => "http://localhost:8080/api/user/123"
//	FillPathVariables("http://localhost:8080/api/user/:id/:action", PathVar{"id": "123", "action": "edit"}) => "http://localhost:8080/api/user/123/edit"
//	FillPathVariables("http://localhost:8080/api/user/:id/:action", PathVar{"id": "123"}) => "http://localhost:8080/api/user/123/:action"
func FillPathVariables(url string, vars PathVar) string {
	if len(vars) == 0 {
		return url
	}

	return pathVarReg.ReplaceAllStringFunc(url, func(match string) string {
		v, ok := vars[match[1:]]
		if ok {
			return v
		} else {
			return match
		}
	})
}
