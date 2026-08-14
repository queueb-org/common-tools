package tools

import "os"

// Exists check location exists
func Exists(location string) bool {
	if _, err := os.Stat(location); err == nil {
		return true
	}
	return false
}

// IsDir checks if the location belongs a directory
func IsDir(location string) bool {
	if stat, err := os.Stat(location); err == nil {
		return stat.IsDir()
	}
	return false
}

// IsFile checks if the location belongs a file
func IsFile(location string) bool {
	// NOTE, this is not `repeat-yourself` case, due to
	// error will always be treated as false (don't use here !IsDir(location))
	if stat, err := os.Stat(location); err == nil {
		return !stat.IsDir()
	}
	return false
}
