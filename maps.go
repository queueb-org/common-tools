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
func TraverseMapString(in any, path string) any {
	if strings.Contains(path, ".") {
		idx := strings.Index(path, ".")
		left := path[:idx]
		right := path[idx+1:]
		if item, ok := in.(map[string]any); ok {
			return TraverseMapString(item[left], right)
		}
	}
	if item, ok := in.(map[string]any); ok {
		return item[path]
	}
	return (any)(nil)
}
