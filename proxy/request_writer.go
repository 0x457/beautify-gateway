package proxy

import (
	"strings"
)

// DesliceValues is used to collapse single value string slices from map values.
func desliceValues(slice map[string][]string) map[string]interface{} {
	desliced := make(map[string]interface{})
	for k, v := range slice {
		if len(v) == 1 {
			desliced[k] = v[0]
		} else {
			desliced[k] = v
		}
	}
	return desliced
}

//MuxRouterPath  url params
//MuxRouterPath("/api/{key}","key","123")=/api/123
func MuxRouterPath(url string, vars map[string]string) string {
	paths := strings.Split(url, "/")
	for index := 0; index < len(paths); index++ {
		p := paths[index]
		if len(p) <= 1 {
			continue
		}
		if p[0] == '{' && p[len(p)-1] == '}' {
			key := p[1 : len(p)-1]
			if v, ok := vars[key]; ok {
				paths[index] = v
			} else {
				return ""
			}
		}
	}
	return strings.Join(paths, "/")
}
