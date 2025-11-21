package endpoints

import "strings"

// joinPaths safely joins multiple URL path segments
func JoinPaths(base string, paths ...string) string {
	// Remove trailing slash from base if present
	base = strings.TrimRight(base, "/")

	// Trim slashes from each path segment and join with "/"
	for i, p := range paths {
		paths[i] = strings.Trim(p, "/")
	}

	fullPath := strings.Join(paths, "/")
	return base + "/" + fullPath
}
