package endpoints

// TODO: remove this, all URL building should be on existing URL types
func BuildURL(baseURL, path string) string {
	return baseURL + path
}
