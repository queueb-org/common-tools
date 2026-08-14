package tools

import (
	"reflect"
	"time"
)

// IsEmptier checks if object is empty. It's almost like reflect.Value.IsZero,
// however it checks only contents of the object
// In comparison to IsZero, an empty slice is treated like not zero entry, but
// it's empty (blank).
type IsEmptier interface {
	IsEmpty() bool
}

// Pick provides picking first given non zero argument.
func Pick[T comparable](args ...T) T {
	def := *new(T)

	for _, arg := range args {
		if arg == def {
			continue
		} else {
			return arg
		}
	}

	return def
}

// Or simulates python-like approach with evaluating left part of expression
// if right operands and or operator are not blank:
// in python {} or True gives True:
//
//	Or(0, true) -> true
func Or(parts ...interface{}) interface{} {
	if len(parts) == 0 {
		return nil
	}

	for _, entry := range parts {
		if !isEmpty(entry) {
			return entry
		}
	}

	//: use the latest argument (or should it be the first?)
	return parts[len(parts)-1]
}

// And simulates python-like approach with evaluating left part of expression
// if right operands and or operator are not blank:
// in python: ["test"] and True or {} gives {}:
//
//	And([]string{"test}, true, map[string]string{}) -> map[string]string{}
func And(parts ...interface{}) interface{} {
	if len(parts) == 0 {
		return nil
	}

	for _, entry := range parts {
		if isEmpty(entry) {
			return entry
		}
	}

	return parts[len(parts)-1]
}

func isEmpty(entry interface{}) bool {
	//: on nil explicit true (i.e. blank)
	if entry == nil {
		return true
	}

	switch typed := entry.(type) {
	case string:
		return typed == ""
	// NOTE, you can't unite all int/uint types with typed == 0 expression
	// despite the compiler treats it like valid expression, the result will
	// be always false, even if it's 0.
	case int:
		return typed == 0
	case int16:
		return typed == 0
	case int32:
		return typed == 0
	case int64:
		return typed == 0
	case uint:
		return typed == 0
	case uint16:
		return typed == 0
	case uint32:
		return typed == 0
	case uint64:
		return typed == 0
	case float32, float64:
		return typed == 0.0
	case bool:
		return !typed
	case complex64:
		return real(typed) == 0 && imag(typed) == 0
	case complex128:
		return real(typed) == 0 && imag(typed) == 0
	case time.Duration:
		return typed == time.Duration(0)
	default:
		return isComplexObjEmpty(entry)
	}
}

func isComplexObjEmpty(entry interface{}) bool {
	fieldType := reflect.TypeOf(entry)
	field := reflect.ValueOf(entry)

	switch fieldType.Kind() {
	case reflect.Ptr:
		if field.Elem().Kind() == reflect.Struct {
			if i, ok := field.Interface().(IsEmptier); ok {
				return i.IsEmpty()
			}
			return true
		}
		val := field.Elem().Interface()
		return isEmpty(val)
	case reflect.Slice, reflect.Map:
		return field.Len() == 0
	case reflect.Struct:
		return true
	case reflect.Func:
		return field.IsNil()
	default:
		return true
	}
}
