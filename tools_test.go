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
		in             []any
		expected       any
		checkOnPointer bool
	}{
		{name: "nil", in: []any{nil, nil, nil}, expected: nil},
		{name: "blank", in: []any{}, expected: nil},
		{
			name:     "ints/all-blanks",
			in:       []any{int(0), int16(0), int32(0), int64(0)},
			expected: int64(0),
		},
		{
			name:     "uints/all-blanks",
			in:       []any{uint(0), uint16(0), uint32(0), uint64(0)},
			expected: uint64(0),
		},
		{name: "strings/no-blanks", in: []any{"test", "me", "up"}, expected: "test"},
		{name: "strings/blanks", in: []any{"", "me", "", "", "test"}, expected: "me"},
		{name: "strings/all-blanks", in: []any{"", "", "", "", ""}, expected: ""},
		{name: "mixed/no-blanks", in: []any{1337, true, 3.1, []string{"test"}}, expected: 1337},
		{
			name:     "mixed/blanks",
			in:       []any{0, false, time.Duration(0), 0.0, []string{"test"}, true},
			expected: []string{"test"},
		},
		{
			name: "mixed/all-blanks",
			in: []any{
				0, false, 0.0, complex64(complex(0, 0)), complex(0, 0), []string{}, []byte{},
			},
			expected: []byte{},
		},
		{
			name:     "slices/all-blanks",
			in:       []any{[]string{}, []bool{}, []float32{}, []int16{}, []int32{}},
			expected: []int32{},
		},
		{
			name: "maps/all-blanks",
			in: []any{
				map[string]struct{}{}, map[bool]struct{}{}, map[int]int{}, map[bool]bool{},
			},
			expected: map[bool]bool{},
		},
		{
			name:     "pointers/no-blanks",
			in:       []any{&testBlankString, &testString, &testInt, &testBool},
			expected: &testString,
		},
		{
			name:     "pointers/blanks",
			in:       []any{new(""), new(int32(0)), new(int64(0)), new(false)},
			expected: new(false),
		},
		{
			name:     "structs/ptr-blanks",
			in:       []any{&blank{}, &blank{}, blank{}, &value{"me"}, blank{}},
			expected: &value{"me"},
		},
		{
			name:     "structs/blanks",
			in:       []any{blank{}, &blank{}, blank{}, value{"me"}, blank{}},
			expected: blank{},
		},
		{
			name:           "funcs",
			in:             []any{noFunc, okFunc},
			expected:       okFunc,
			checkOnPointer: true,
		},
		{
			name:     "not-supported/chan-1",
			in:       []any{chanInt, 1},
			expected: 1,
		},
		{
			name:     "not-supported/chan-2",
			in:       []any{0, 0, false, chanInt},
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
		in       []any
		expected any
	}{
		{"nil", []any{nil, nil, nil}, nil},
		{"blank", []any{}, nil},
		{
			"ints/all-blanks", []any{int(0), int16(0), int32(0), int64(0)},
			int(0),
		},
		{
			"uints/all-blanks", []any{uint(0), uint16(0), uint32(0), uint64(0)},
			uint(0),
		},
		{"strings/no-blanks", []any{"test", "me", "up"}, "up"},
		{"strings/blanks", []any{"me", "test", "", "and", "this"}, ""},
		{"strings/all-blanks", []any{"", "", "", "", ""}, ""},
		{"mixed/no-blanks", []any{1337, true, 3.1, []string{"test"}}, []string{"test"}},
		{"mixed/blanks", []any{1, true, 1.3, []byte{}, []string{"test"}, true}, []byte{}},
		{
			"mixed/all-blanks",
			[]any{
				0, false, 0.0, complex64(complex(0, 0)), complex(0, 0), []string{}, []byte{},
			},
			0,
		},
		{"slices/all-blanks", []any{[]string{}, []bool{}, []float32{}, []int16{}, []int32{}}, []string{}},
		{
			"maps/all-blanks",
			[]any{
				map[string]struct{}{}, map[bool]struct{}{}, map[int]int{}, map[bool]bool{},
			},
			map[string]struct{}{},
		},
		{
			"pointers/no-blanks",
			[]any{&testString, &testInt, &testBool},
			&testBool,
		},
		{
			"pointers/blanks",
			[]any{new(""), new(int32(0)), new(int64(0)), new(false)},
			new(""),
		},
		{
			"structs/ptr-blanks",
			[]any{&blank{}, &blank{}, blank{}, &value{"me"}, blank{}},
			&blank{},
		},
		{
			"structs/blanks",
			[]any{blank{}, &blank{}, blank{}, value{"me"}, blank{}},
			blank{},
		},
		{
			"not-supported/chan-1",
			[]any{1, 2, 3, 1337, true, chanInt, true, "test", 0, false},
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
