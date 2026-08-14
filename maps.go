package tools

import "strings"

// TraverseMapString follows map[string]interface{} instances over
// fields stored in path and separated with `.`
// Example:
//
//	in := map[string]interface{}{
//		"test": map[string]interface{}{
//			"one": "test",
//			"two": 1337,
//		}
//	}
//	TraverseMapString(in, "test.one").(string)  # should be "test"
//	TraverseMapString(in, "test.two").(string)  # should be 1337
//	TraverseMapString(in, "test.nonexistent")   # nil
func TraverseMapString(in interface{}, path string) interface{} {
	if strings.Contains(path, ".") {
		idx := strings.Index(path, ".")
		left := path[:idx]
		right := path[idx+1:]
		if item, ok := in.(map[string]interface{}); ok {
			return TraverseMapString(item[left], right)
		}
	}
	if item, ok := in.(map[string]interface{}); ok {
		return item[path]
	}
	return (interface{})(nil)
}
