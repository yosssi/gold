package gold

import "strings"

// Path constructs a path and returns it.
func Path(baseDir, path string) string {
	if CurrentDirectoryBasedPath(path) || AbsolutePath(path) {
		return path
	}
	if baseDir != "" {
		return baseDir + "/" + path
	}
	return path
}

// CurrentDirectoryBasedPath checks whether the path is based on a current
// directory or not and returns the result.
func CurrentDirectoryBasedPath(path string) bool {
	return strings.HasPrefix(path, "./") || strings.HasPrefix(path, "../")
}

// AbsolutePath checks whether the path is an absolute path or not and
// returns the result.
func AbsolutePath(path string) bool {
	return strings.HasPrefix(path, "/")
}

// assetPath constructs an asset path and returns it.
func assetPath(path, baseDir string) string {
	if strings.HasPrefix(path, baseDir) {
		path = strings.Replace(path, baseDir, "", 1)
		if strings.HasPrefix(path, "/") {
			path = strings.Replace(path, "/", "", 1)
		}
	}
	if strings.HasPrefix(path, "./") {
		path = strings.Replace(path, "./", "", 1)
	}
	path = strings.Replace(path, "/./", "/", -1)
	var tokens []string
	for _, s := range strings.Split(path, "/") {
		if s == ".." {
			if l := len(tokens); l > 0 {
				tokens = tokens[0 : l-1]
			}
		} else {
			tokens = append(tokens, s)
		}
	}
	return strings.Join(tokens, "/")
}
