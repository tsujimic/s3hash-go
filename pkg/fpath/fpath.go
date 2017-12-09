package fpath

import "strings"

// SplitPath ...
func SplitPath(path string) (string, string) {
	path = strings.TrimLeft(path, "/")

	i := 0
	for i < len(path) && path[i] != '/' {
		i++
	}

	return path[:i], strings.TrimLeft(path[i:], "/")
}

// SplitName ...
func SplitName(path string) string {
	path = strings.TrimRight(path, "/")

	i := len(path) - 1
	for i >= 0 && path[i] != '/' {
		i--
	}

	return path[i+1:]
}

// func splitPath(path string) (bucket, object string) {
// 	path = strings.TrimLeft(path, "/")

// 	i := 0
// 	for i < len(path) && path[i] != '/' {
// 		i++
// 	}

// 	return path[:i], strings.TrimLeft(path[i:], "/")
// }

// func splitName(path string) (name string) {
// 	path = strings.TrimRight(path, "/")

// 	i := len(path) - 1
// 	for i >= 0 && path[i] != '/' {
// 		i--
// 	}

// 	return path[i+1:]
// }

// func splitPath(path string) (container, blob string) {
// 	path = strings.TrimLeft(path, "/")

// 	i := 0
// 	for i < len(path) && path[i] != '/' {
// 		i++
// 	}

// 	return path[:i], strings.TrimLeft(path[i:], "/")
// }

// func splitName(path string) (name string) {
// 	path = strings.TrimRight(path, "/")

// 	i := len(path) - 1
// 	for i >= 0 && path[i] != '/' {
// 		i--
// 	}

// 	return path[i+1:]
// }
