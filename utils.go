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
