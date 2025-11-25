package endpoints

import "strings"

// joinPaths safely joins multiple URL path segments
func JoinPaths(base string, paths ...string) string {
	base = strings.TrimRight(base, "/")

	for i, p := range paths {
		paths[i] = strings.Trim(p, "/")
	}

	fullPath := strings.Join(paths, "/")
	return base + "/" + fullPath
}
