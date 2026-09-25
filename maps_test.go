package tools

import "testing"

func TestTraverseMapString(t *testing.T) {
	data := map[string]any{
		"test": map[string]any{
			"one":   "a string",
			"two":   true,
			"three": 1337,
			"four":  1337.33,
			//: and so forth including next level of map[string]interface{}
		},
	}

	for _, entry := range []struct {
		path     string
		expected any
	}{
		{"test.one", "a string"},
		{"test.two", true},
		{"test.three", 1337},
		{"test.four", 1337.33},
		{"test.nonexistent", nil},
		{"1.2.3.4", nil},
	} {
		t.Run(entry.path, func(in *testing.T) {
			if result := TraverseMapString(data, entry.path); result != entry.expected {
				t.Errorf("expected: %v, got: %v", entry.expected, result)
			}
		})
	}
}
