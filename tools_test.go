package tools

import (
	"reflect"
	"testing"
	"time"
)

var (
	testBlankString = ""
	testString      = "test value"
	testInt         = 1337
	testBool        = false
)

type value struct {
	value string
}

func (t *value) IsEmpty() bool {
	return t.value == ""
}

type ClientFunc = func()

var (
	noFunc ClientFunc
	okFunc = func() {
		//lint:ignore S1023, it's intended.
		return
	}
)

func TestOr(t *testing.T) {
	type blank struct{}
	chanInt := make(chan int)

	for _, entry := range []struct {
		name           string
		in             []interface{}
		expected       interface{}
		checkOnPointer bool
	}{
		{name: "nil", in: []interface{}{nil, nil, nil}, expected: nil},
		{name: "blank", in: []interface{}{}, expected: nil},
		{
			name:     "ints/all-blanks",
			in:       []interface{}{int(0), int16(0), int32(0), int64(0)},
			expected: int64(0),
		},
		{
			name:     "uints/all-blanks",
			in:       []interface{}{uint(0), uint16(0), uint32(0), uint64(0)},
			expected: uint64(0),
		},
		{name: "strings/no-blanks", in: []interface{}{"test", "me", "up"}, expected: "test"},
		{name: "strings/blanks", in: []interface{}{"", "me", "", "", "test"}, expected: "me"},
		{name: "strings/all-blanks", in: []interface{}{"", "", "", "", ""}, expected: ""},
		{name: "mixed/no-blanks", in: []interface{}{1337, true, 3.1, []string{"test"}}, expected: 1337},
		{
			name:     "mixed/blanks",
			in:       []interface{}{0, false, time.Duration(0), 0.0, []string{"test"}, true},
			expected: []string{"test"},
		},
		{
			name: "mixed/all-blanks",
			in: []interface{}{
				0, false, 0.0, complex64(complex(0, 0)), complex(0, 0), []string{}, []byte{},
			},
			expected: []byte{},
		},
		{
			name:     "slices/all-blanks",
			in:       []interface{}{[]string{}, []bool{}, []float32{}, []int16{}, []int32{}},
			expected: []int32{},
		},
		{
			name: "maps/all-blanks",
			in: []interface{}{
				map[string]struct{}{}, map[bool]struct{}{}, map[int]int{}, map[bool]bool{},
			},
			expected: map[bool]bool{},
		},
		{
			name:     "pointers/no-blanks",
			in:       []interface{}{&testBlankString, &testString, &testInt, &testBool},
			expected: &testString,
		},
		{
			name:     "pointers/blanks",
			in:       []interface{}{new(""), new(int32(0)), new(int64(0)), new(false)},
			expected: new(false),
		},
		{
			name:     "structs/ptr-blanks",
			in:       []interface{}{&blank{}, &blank{}, blank{}, &value{"me"}, blank{}},
			expected: &value{"me"},
		},
		{
			name:     "structs/blanks",
			in:       []interface{}{blank{}, &blank{}, blank{}, value{"me"}, blank{}},
			expected: blank{},
		},
		{
			name:           "funcs",
			in:             []interface{}{noFunc, okFunc},
			expected:       okFunc,
			checkOnPointer: true,
		},
		{
			name:     "not-supported/chan-1",
			in:       []interface{}{chanInt, 1},
			expected: 1,
		},
		{
			name:     "not-supported/chan-2",
			in:       []interface{}{0, 0, false, chanInt},
			expected: chanInt,
		},
	} {
		t.Run(entry.name, func(in *testing.T) {
			result := Or(entry.in...)
			switch entry.checkOnPointer {
			case true:
				if reflect.ValueOf(result).Pointer() != reflect.ValueOf(entry.expected).Pointer() {
					in.Errorf("pointers aren't match: %#+v != %#+v", result, entry.expected)
				}
			default:
				if !reflect.DeepEqual(result, entry.expected) {
					in.Errorf("expected: `%#+v`, got: `%#+v`", entry.expected, result)
				}
			}
		})
	}
}

func TestAnd(t *testing.T) {
	type blank struct{}
	chanInt := make(chan int)

	for _, entry := range []struct {
		name     string
		in       []interface{}
		expected interface{}
	}{
		{"nil", []interface{}{nil, nil, nil}, nil},
		{"blank", []interface{}{}, nil},
		{
			"ints/all-blanks", []interface{}{int(0), int16(0), int32(0), int64(0)},
			int(0),
		},
		{
			"uints/all-blanks", []interface{}{uint(0), uint16(0), uint32(0), uint64(0)},
			uint(0),
		},
		{"strings/no-blanks", []interface{}{"test", "me", "up"}, "up"},
		{"strings/blanks", []interface{}{"me", "test", "", "and", "this"}, ""},
		{"strings/all-blanks", []interface{}{"", "", "", "", ""}, ""},
		{"mixed/no-blanks", []interface{}{1337, true, 3.1, []string{"test"}}, []string{"test"}},
		{"mixed/blanks", []interface{}{1, true, 1.3, []byte{}, []string{"test"}, true}, []byte{}},
		{
			"mixed/all-blanks",
			[]interface{}{
				0, false, 0.0, complex64(complex(0, 0)), complex(0, 0), []string{}, []byte{},
			},
			0,
		},
		{"slices/all-blanks", []interface{}{[]string{}, []bool{}, []float32{}, []int16{}, []int32{}}, []string{}},
		{
			"maps/all-blanks",
			[]interface{}{
				map[string]struct{}{}, map[bool]struct{}{}, map[int]int{}, map[bool]bool{},
			},
			map[string]struct{}{},
		},
		{
			"pointers/no-blanks",
			[]interface{}{&testString, &testInt, &testBool},
			&testBool,
		},
		{
			"pointers/blanks",
			[]interface{}{new(""), new(int32(0)), new(int64(0)), new(false)},
			new(""),
		},
		{
			"structs/ptr-blanks",
			[]interface{}{&blank{}, &blank{}, blank{}, &value{"me"}, blank{}},
			&blank{},
		},
		{
			"structs/blanks",
			[]interface{}{blank{}, &blank{}, blank{}, value{"me"}, blank{}},
			blank{},
		},
		{
			"not-supported/chan-1",
			[]interface{}{1, 2, 3, 1337, true, chanInt, true, "test", 0, false},
			chanInt,
		},
	} {
		t.Run(entry.name, func(in *testing.T) {
			if result := And(entry.in...); !reflect.DeepEqual(result, entry.expected) {
				in.Errorf("expected: %#+v, got: %#+v", entry.expected, result)
			}
		})
	}
}

func TestPick(t *testing.T) {
	for _, entry := range []struct {
		name     string
		in       []string
		expected string
	}{
		{"defaults", []string{}, ""},
		{"first-non-zero", []string{"", "", "", "non-zero", "", "last"}, "non-zero"},
	} {
		t.Run(entry.name, func(in *testing.T) {
			if result := Pick(entry.in...); result != entry.expected {
				in.Errorf("expected: %#+v, got: %#+v", entry.expected, result)
			}
		})
	}
}
