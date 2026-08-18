package tools

import (
	"os"
	"reflect"
	"testing"
	"time"
)

// E reflects environment variable (for testing purpose).
type E struct {
	Name  string
	Value string
}

// WithSetEnv sets environment variables in better/convenient form.
func WithSetEnv(t testing.TB, envs ...*E) {
	for _, e := range envs {
		t.Setenv(e.Name, e.Value)
	}
}

// WithUnsetEnv removes defined environment variables and restores them back
// once test is over.
func WithUnsetEnv(t testing.TB, envs ...*E) {
	orig := make([]*E, 0, len(envs))

	for _, e := range envs {
		if v, found := os.LookupEnv(e.Name); found {
			orig = append(orig, &E{Name: e.Name, Value: v})
		}
		_ = os.Unsetenv(e.Name)
	}

	t.Cleanup(func() {
		for _, e := range orig {
			_ = os.Setenv(e.Name, e.Value)
		}
	})
}

func NewE(name, value string) *E {
	return &E{Name: name, Value: value}
}

// testEnv testing helper
func testEnv[T EnvValue](t testing.TB, envName, envValue string, fallback, expected T) {
	t.Setenv(envName, envValue)

	// tests.WithSetEnv(t, tests.NewE(envName, envValue).WithUnset(true))

	result := Env(envName, fallback)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected: %v, got: %v", expected, result)
	}
}

// Tests

func TestEnv(t *testing.T) {
	t.Run("ok", func(in *testing.T) {
		expected := "expected value"
		envs := []*E{
			NewE("TEST_VARIABLE", expected),
		}
		WithSetEnv(in, envs...)

		if result := Env("TEST_VARIABLE", "value"); result != expected {
			in.Errorf("expected: %v, got: %v", expected, result)
		}

		if result := Env("TEST_NON_EXISTENT_VARIABLE", "value"); result != "value" {
			in.Errorf("fallback expected got: %v", result)
		}
	})

	t.Run("bool", func(in *testing.T) {
		testEnv(in, "TEST_INPUT_BOOL", "true", false, true)
		testEnv(in, "TEST_INPUT_BOOL", "false", true, false)
	})

	t.Run("ints", func(in *testing.T) {
		testEnv(in, "TEST_INPUT_INTS", "1", int(0), int(1))
		testEnv(in, "TEST_INPUT_INTS", "80", int8(0), int8(80))
		testEnv(in, "TEST_INPUT_INTS", "160", int16(0), int16(160))
		testEnv(in, "TEST_INPUT_INTS", "3200", int32(0), int32(3200))
		testEnv(in, "TEST_INPUT_INTS", "6400", int64(0), int64(6400))
	})

	t.Run("uints", func(in *testing.T) {
		testEnv(in, "TEST_INPUT_UINTS", "1", uint(0), uint(1))
		testEnv(in, "TEST_INPUT_UINTS", "80", uint8(0), uint8(80))
		testEnv(in, "TEST_INPUT_UINTS", "160", uint16(0), uint16(160))
		testEnv(in, "TEST_INPUT_UINTS", "3200", uint32(0), uint32(3200))
		testEnv(in, "TEST_INPUT_UINTS", "6400", uint64(0), uint64(6400))
	})

	t.Run("floats", func(in *testing.T) {
		testEnv(in, "TEST_INPUT_FLOATS", "1337", float32(0), float32(1337))
		testEnv(in, "TEST_INPUT_FLOATS", "1337", float64(0), float64(1337))
	})

	t.Run("[]string", func(in *testing.T) {
		testEnv(in, "TEST_INPUT_STRING_SLICE", "1,2,3,4", nil, []string{"1", "2", "3", "4"})
	})

	t.Run("[]ints", func(in *testing.T) {
		testEnv(in, "TEST_INPUT_INT_SLICE", "1,2,3,4", nil, []int{1, 2, 3, 4})
		testEnv(in, "TEST_INPUT_INT_SLICE", "1,2,3,4", nil, []int8{1, 2, 3, 4})
		testEnv(in, "TEST_INPUT_INT_SLICE", "1,2,3,4", nil, []int16{1, 2, 3, 4})
		testEnv(in, "TEST_INPUT_INT_SLICE", "1,2,3,4", nil, []int32{1, 2, 3, 4})
		testEnv(in, "TEST_INPUT_INT_SLICE", "1,2,3,4", nil, []int64{1, 2, 3, 4})
	})

	t.Run("map[string]string", func(in *testing.T) {
		testEnv(in, "TEST_INPUT_STRING_MAP", "1=2,2=3,3=4,4=5", nil, map[string]string{"1": "2", "2": "3", "3": "4", "4": "5"})
	})
}

// note, this test runs only code coverage, for real tests see [TestEnv].
func TestEnvP(t *testing.T) {
	EnvP(MakeEnv("APP"), "TEST", "value")
}

func testEnvIntSlice[T ints](test *testing.T, in string, expected any) {
	if result := EnvIntSlice[T](in, []T{0}); !reflect.DeepEqual(result, expected) {
		test.Errorf("expected: %v, got: %v", expected, result)
	}
}

func TestEnvIntSlice(t *testing.T) {
	input := "1, 2, 3,4 ,5, 6, test, me"
	env := "TEST_INPUT"
	WithSetEnv(t, NewE("TEST_INPUT", input))

	t.Run("int", func(in *testing.T) {
		testEnvIntSlice[int](in, env, []int{1, 2, 3, 4, 5, 6, 0, 0})
	})

	t.Run("int8", func(in *testing.T) {
		testEnvIntSlice[int8](in, env, []int8{1, 2, 3, 4, 5, 6, 0, 0})
	})

	t.Run("int16", func(in *testing.T) {
		testEnvIntSlice[int16](in, env, []int16{1, 2, 3, 4, 5, 6, 0, 0})
	})

	t.Run("int32", func(in *testing.T) {
		testEnvIntSlice[int32](in, env, []int32{1, 2, 3, 4, 5, 6, 0, 0})
	})

	t.Run("int64", func(in *testing.T) {
		testEnvIntSlice[int64](in, env, []int64{1, 2, 3, 4, 5, 6, 0, 0})
	})

	t.Run("fallback", func(in *testing.T) {
		testEnvIntSlice[int](in, "test", []int{0})
	})
}

func TestEnvInts(t *testing.T) {
	for _, entry := range []struct {
		name     string
		envSet   *E
		env      string
		expected int
	}{
		{
			name:     "ok",
			envSet:   NewE("TEST", "1337"),
			env:      "TEST",
			expected: 1337,
		},
		{
			name:     "fallback",
			env:      "TEST",
			expected: -1337,
		},
		{
			name:     "invalid-format",
			envSet:   NewE("TEST", "this is invalid int value"),
			env:      "TEST",
			expected: -1337,
		},
	} {
		t.Run(entry.name, func(in *testing.T) {
			if entry.envSet != nil {
				WithSetEnv(in, entry.envSet)
			}
			result := EnvInts[int](entry.env, -1337)
			if result != entry.expected {
				t.Errorf("expected: %v, got: %v", entry.expected, result)
			}
		})
	}
}

func TestEnvInt64(t *testing.T) {
	for _, entry := range []struct {
		name     string
		in       string
		env      func(testing.TB, ...*E)
		args     []*E
		expected int64
	}{
		{
			name:     "ok",
			in:       "TEST_INT_ENV",
			env:      WithSetEnv,
			args:     []*E{{Name: "TEST_INT_ENV", Value: "1337"}},
			expected: 1337,
		},
		{
			name:     "wrong-value",
			in:       "TEST_INT_ENV",
			env:      WithSetEnv,
			args:     []*E{{Name: "TEST_INT_ENV", Value: "string-1337"}},
			expected: -1,
		},
		{
			name:     "no-value",
			in:       "TEST_INT_ENV",
			env:      WithUnsetEnv,
			args:     []*E{{Name: "TEST_INT_ENV"}},
			expected: -1,
		},
	} {
		t.Run(entry.name, func(in *testing.T) {
			entry.env(in, entry.args...)

			if result := EnvInt64(entry.in, -1); result != entry.expected {
				in.Errorf("expected: %v, got: %v", entry.expected, result)
			}
		})
	}
}

func TestEnvInt_all(t *testing.T) {
	WithSetEnv(t,
		&E{Name: "TEST_UINT_ENV", Value: "101"},
		&E{Name: "TEST_UINT_FALLBACK_ENV", Value: "wrong-value"})

	EnvUint("TEST_UINT_ENV", 0)

	for _, entry := range []struct {
		name     string
		fallback interface{}
		call     func() interface{}
		expected interface{}
	}{
		{
			name:     "EnvInt8/ok",
			call:     func() interface{} { return EnvInt8("TEST_UINT_ENV", 1) },
			expected: int8(101),
		},
		{
			name:     "EnvInt8/fallback",
			call:     func() interface{} { return EnvInt8("TEST_UINT_FALLBACK_ENV", 101) },
			expected: int8(101),
		},
		{
			name:     "EnvInt16",
			call:     func() interface{} { return EnvInt16("TEST_UINT_ENV", 1) },
			expected: int16(101),
		},
		{
			name:     "EnvUint16/fallback",
			call:     func() interface{} { return EnvInt16("TEST_UINT_FALLBACK_ENV", 101) },
			expected: int16(101),
		},
		{
			name:     "EnvInt32",
			call:     func() interface{} { return EnvInt32("TEST_UINT_ENV", 1) },
			expected: int32(101),
		},
		{
			name:     "EnvInt32/fallback",
			call:     func() interface{} { return EnvInt32("TEST_UINT_FALLBACK_ENV", 101) },
			expected: int32(101),
		},
		{
			name:     "EnvInt",
			call:     func() interface{} { return EnvInt("TEST_UINT_ENV", 1) },
			expected: int(101),
		},
		{
			name:     "EnvInt/fallback",
			call:     func() interface{} { return EnvInt("TEST_UINT_FALLBACK_ENV", 101) },
			expected: 101,
		},
	} {
		t.Run(entry.name, func(in *testing.T) {
			if result := entry.call(); !reflect.DeepEqual(result, entry.expected) {
				in.Errorf("expected: `%v`, got: `%v`", entry.expected, result)
			}
		})
	}
}

func TestEnvUints(t *testing.T) {
	for _, entry := range []struct {
		name     string
		envSet   *E
		env      string
		expected uint
	}{
		{
			name:     "ok",
			envSet:   NewE("TEST", "1337"),
			env:      "TEST",
			expected: 1337,
		},
		{
			name:     "fallback",
			env:      "TEST",
			expected: 101,
		},
		{
			name:     "invalid-format",
			envSet:   NewE("TEST", "this is invalid int value"),
			env:      "TEST",
			expected: 101,
		},
	} {
		t.Run(entry.name, func(in *testing.T) {
			if entry.envSet != nil {
				WithSetEnv(in, entry.envSet)
			}
			result := EnvUints[uint](entry.env, 101)
			if result != entry.expected {
				t.Errorf("expected: %v, got: %v", entry.expected, result)
			}
		})
	}
}

func TestEnvUint64(t *testing.T) {
	for _, entry := range []struct {
		name     string
		in       string
		env      func(testing.TB, ...*E)
		args     []*E
		expected uint64
	}{
		{
			name:     "ok",
			in:       "TEST_UINT_ENV",
			env:      WithSetEnv,
			args:     []*E{{Name: "TEST_UINT_ENV", Value: "1337"}},
			expected: 1337,
		},
		{
			name:     "wrong-value",
			in:       "TEST_UINT_ENV",
			env:      WithSetEnv,
			args:     []*E{{Name: "TEST_UINT_ENV", Value: "string-1337"}},
			expected: 0,
		},
		{
			name:     "no-value",
			in:       "TEST_UINT_ENV",
			env:      WithSetEnv,
			args:     []*E{{Name: "TEST_UINT_ENV"}},
			expected: 0,
		},
	} {
		t.Run(entry.name, func(in *testing.T) {
			entry.env(in, entry.args...)
			if result := EnvUint64(entry.in, 0); result != entry.expected {
				in.Errorf("expected: %v, got: %v", entry.expected, result)
			}
		})
	}
}

func TestEnvUint_all(t *testing.T) {
	WithSetEnv(t,
		&E{Name: "TEST_UINT_ENV", Value: "101"},
		&E{Name: "TEST_UINT_FALLBACK_ENV", Value: "wrong-value"})

	EnvUint("TEST_UINT_ENV", 0)

	for _, entry := range []struct {
		name     string
		fallback interface{}
		call     func() interface{}
		expected interface{}
	}{
		{
			name:     "EnvUint8/ok",
			call:     func() interface{} { return EnvUint8("TEST_UINT_ENV", 1) },
			expected: uint8(101),
		},
		{
			name:     "EnvUint8/fallback",
			call:     func() interface{} { return EnvUint8("TEST_UINT_FALLBACK_ENV", 101) },
			expected: uint8(101),
		},
		{
			name:     "EnvUint16",
			call:     func() interface{} { return EnvUint16("TEST_UINT_ENV", 1) },
			expected: uint16(101),
		},
		{
			name:     "EnvUint16/fallback",
			call:     func() interface{} { return EnvUint16("TEST_UINT_FALLBACK_ENV", 101) },
			expected: uint16(101),
		},
		{
			name:     "EnvUint32",
			call:     func() interface{} { return EnvUint32("TEST_UINT_ENV", 1) },
			expected: uint32(101),
		},
		{
			name:     "EnvUint32/fallback",
			call:     func() interface{} { return EnvUint32("TEST_UINT_FALLBACK_ENV", 101) },
			expected: uint32(101),
		},
		{
			name:     "EnvUint",
			call:     func() interface{} { return EnvUint("TEST_UINT_ENV", 1) },
			expected: uint(101),
		},
		{
			name:     "EnvUint/fallback",
			call:     func() interface{} { return EnvUint("TEST_UINT_FALLBACK_ENV", 101) },
			expected: uint(101),
		},
	} {
		t.Run(entry.name, func(in *testing.T) {
			if result := entry.call(); !reflect.DeepEqual(result, entry.expected) {
				in.Errorf("expected: `%v`, got: `%v`", entry.expected, result)
			}
		})
	}
}

func TestEnvDuration(t *testing.T) {
	for _, entry := range []struct {
		name     string
		in       string
		env      func(testing.TB, ...*E)
		args     []*E
		expected time.Duration
	}{
		{
			name:     "ok",
			in:       "TEST_DURATION_ENV",
			env:      WithSetEnv,
			args:     []*E{{Name: "TEST_DURATION_ENV", Value: "1s"}},
			expected: time.Second * 1,
		},
		{
			name:     "ok/complex",
			in:       "TEST_DURATION_ENV",
			env:      WithSetEnv,
			args:     []*E{{Name: "TEST_DURATION_ENV", Value: "12h30m15s"}},
			expected: time.Hour*12 + time.Minute*30 + time.Second*15,
		},
		{
			name:     "wrong-value",
			in:       "TEST_DURATION_ENV",
			env:      WithSetEnv,
			args:     []*E{{Name: "TEST_DURATION_ENV", Value: "string-1337"}},
			expected: 1,
		},
		{
			name:     "no-value",
			in:       "TEST_DURATION_ENV",
			env:      WithUnsetEnv,
			args:     []*E{{Name: "TEST_DURATION_ENV"}},
			expected: 1,
		},
	} {
		t.Run(entry.name, func(in *testing.T) {
			entry.env(in, entry.args...)
			if result := EnvDuration(entry.in, 1); result != entry.expected {
				in.Errorf("expected: %v, got: %v", entry.expected, result)
			}
		})
	}
}

func TestEnvBool(t *testing.T) {
	for _, entry := range []struct {
		name     string
		in       string
		env      func(testing.TB, ...*E)
		args     []*E
		expected bool
	}{
		{
			name:     "ok",
			in:       "TEST_BOOL_ENV",
			env:      WithSetEnv,
			args:     []*E{{Name: "TEST_BOOL_ENV", Value: "true"}},
			expected: true,
		},
		{
			name:     "wrong-value",
			in:       "TEST_BOOL_ENV",
			env:      WithSetEnv,
			args:     []*E{{Name: "TEST_BOOL_ENV", Value: "string-1337"}},
			expected: false,
		},
		{
			name:     "no-value",
			in:       "TEST_BOOL_ENV",
			env:      WithUnsetEnv,
			args:     []*E{{Name: "TEST_BOOL_ENV"}},
			expected: false,
		},
	} {
		t.Run(entry.name, func(in *testing.T) {
			entry.env(in, entry.args...)
			if result := EnvBool(entry.in, false); result != entry.expected {
				in.Errorf("expected: %v, got: %v", entry.expected, result)
			}
		})
	}
}

func TestEnvSting(t *testing.T) {
	for _, entry := range []struct {
		name     string
		in       string
		env      func(testing.TB, ...*E)
		args     []*E
		expected string
	}{
		{
			name:     "ok",
			in:       "TEST_STRING_ENV",
			env:      WithSetEnv,
			args:     []*E{{Name: "TEST_STRING_ENV", Value: "true"}},
			expected: "true",
		},
		{
			name:     "no-value",
			in:       "TEST_STRING_ENV",
			env:      WithUnsetEnv,
			args:     []*E{{Name: "TEST_STRING_ENV"}},
			expected: "fallback",
		},
	} {
		t.Run(entry.name, func(in *testing.T) {
			entry.env(in, entry.args...)
			if result := EnvString(entry.in, "fallback"); result != entry.expected {
				in.Errorf("expected: %v, got: %v", entry.expected, result)
			}
		})
	}
}

func TestEnvStingSlice(t *testing.T) {
	for _, entry := range []struct {
		name     string
		in       string
		env      func(testing.TB, ...*E)
		args     []*E
		expected []string
	}{
		{
			name:     "ok",
			in:       "TEST_STRING_ENV",
			env:      WithSetEnv,
			args:     []*E{{Name: "TEST_STRING_ENV", Value: "this,is, the , test "}},
			expected: []string{"this", "is", "the", "test"},
		},
		{
			name:     "no-value",
			in:       "TEST_STRING_ENV",
			env:      WithUnsetEnv,
			args:     []*E{{Name: "TEST_STRING_ENV"}},
			expected: []string{"fallback"},
		},
	} {
		t.Run(entry.name, func(in *testing.T) {
			entry.env(in, entry.args...)
			result := EnvStringSlice(entry.in, []string{"fallback"})
			if !reflect.DeepEqual(result, entry.expected) {
				in.Errorf("expected: %v, got: %v", entry.expected, result)
			}
		})
	}
}

func TestEnvStringMap(t *testing.T) {
	for _, entry := range []struct {
		name     string
		in       string
		env      func(testing.TB, ...*E)
		args     []*E
		expected map[string]string
	}{
		{
			name: "ok",
			in:   "TEST_STRING_ENV",
			env:  WithSetEnv,
			args: []*E{
				{
					Name:  "TEST_STRING_ENV",
					Value: "X-TEST=value,  X-FLAG=1337, X-DOUBLE-VALUE=value=1337, unsupported stuff",
				},
			},
			expected: map[string]string{
				"X-TEST": "value", "X-FLAG": "1337", "X-DOUBLE-VALUE": "value=1337",
			},
		},
		{
			name:     "no-value",
			in:       "TEST_STRING_ENV",
			env:      WithUnsetEnv,
			args:     []*E{{Name: "TEST_STRING_ENV"}},
			expected: map[string]string{"fallback": "yes"},
		},
	} {
		t.Run(entry.name, func(in *testing.T) {
			entry.env(in, entry.args...)
			result := EnvStringMap(entry.in, map[string]string{"fallback": "yes"})
			if !reflect.DeepEqual(result, entry.expected) {
				in.Errorf("expected: %v, got: %v", entry.expected, result)
			}
		})
	}
}

func TestEnvFloats(t *testing.T) {
	for _, entry := range []struct {
		name     string
		envSet   *E
		env      string
		expected float64
	}{
		{
			name:     "ok",
			envSet:   NewE("TEST", "1337"),
			env:      "TEST",
			expected: 1337,
		},
		{
			name:     "fallback",
			env:      "TEST",
			expected: -1337,
		},
		{
			name:     "invalid-format",
			envSet:   NewE("TEST", "this is invalid int value"),
			env:      "TEST",
			expected: -1337,
		},
	} {
		t.Run(entry.name, func(in *testing.T) {
			if entry.envSet != nil {
				WithSetEnv(in, entry.envSet)
			}
			result := EnvFloats[float64](entry.env, -1337)
			if result != entry.expected {
				t.Errorf("expected: %v, got: %v", entry.expected, result)
			}
		})
	}
}

func TestEnvFloat(t *testing.T) {
	for _, entry := range []struct {
		name     string
		in       string
		env      func(testing.TB, ...*E)
		args     []*E
		expected float64
	}{
		{
			name:     "ok",
			in:       "TEST_FLOAT_ENV",
			env:      WithSetEnv,
			args:     []*E{{Name: "TEST_FLOAT_ENV", Value: "1337"}},
			expected: 1337,
		},
		{
			name:     "wrong-value",
			in:       "TEST_FLOAT_ENV",
			env:      WithSetEnv,
			args:     []*E{{Name: "TEST_FLOAT_ENV", Value: "string-1337"}},
			expected: -1,
		},
		{
			name:     "no-value",
			in:       "TEST_FLOAT_ENV",
			env:      WithUnsetEnv,
			args:     []*E{{Name: "TEST_FLOAT_ENV"}},
			expected: -1,
		},
	} {
		t.Run(entry.name, func(in *testing.T) {
			entry.env(in, entry.args...)
			if result := EnvFloat(entry.in, -1); result != entry.expected {
				in.Errorf("expected: %v, got: %v", entry.expected, result)
			}
		})
	}
}

func TestEnvUsageP(t *testing.T) {
	expected := "test\n[env: APP_ME]"
	result := EnvUsageP(MakeEnv("APP"), "ME", "test")
	if result != expected {
		t.Errorf("expected: %v, got: %v", expected, result)
	}
}
