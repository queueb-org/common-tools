package tools

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// EnvValue supported values for [Env] helper.
type EnvValue interface {
	// bases
	ints | uints | floats | string | bool |
		// slices:
		[]string | []int | []int8 | []int16 | []int32 | []int64 |
		// maps:
		map[string]string
}

// PrefixFunc is a function that should generate environment name
// prefixes.
// Example:
//
//	prefixer := MakeEnv("MY_APP")
//	// would expect "MY_APP_CUSTOM_FLAG" environment variable.
//	_ = EnvP(prefixer, "custom-flag", "test")
type PrefixFunc = func(in string) string

// EnvP works as [Env] but environment prefixer function is given as first argument.
func EnvP[T EnvValue](prefixer PrefixFunc, name string, fallback T) T {
	return Env(prefixer(name), fallback)
}

// Env reads environment variable and if it's not blank returns it,
// otherwise returns fallback one.
func Env[T EnvValue](name string, fallback T) T {
	result := os.Getenv(name)
	if result == "" {
		return fallback
	}

	var ret T
	switch v := any(fallback).(type) {
	case bool:
		ret = any(EnvBool(name, v)).(T)
	case string:
		ret = any(EnvString(name, v)).(T)
	case int:
		ret = any(EnvInts(name, v)).(T)
	case int8:
		ret = any(EnvInts(name, v)).(T)
	case int16:
		ret = any(EnvInts(name, v)).(T)
	case int32:
		ret = any(EnvInts(name, v)).(T)
	case int64:
		ret = any(EnvInts(name, v)).(T)
	case uint:
		ret = any(EnvUints(name, v)).(T)
	case uint8:
		ret = any(EnvUints(name, v)).(T)
	case uint16:
		ret = any(EnvUints(name, v)).(T)
	case uint32:
		ret = any(EnvUints(name, v)).(T)
	case uint64:
		ret = any(EnvUints(name, v)).(T)
	case float32:
		ret = any(EnvFloats(name, v)).(T)
	case float64:
		ret = any(EnvFloats(name, v)).(T)
	case []string:
		ret = any(EnvStringSlice(name, v)).(T)
	case []int:
		ret = any(EnvIntSlice(name, v)).(T)
	case []int8:
		ret = any(EnvIntSlice(name, v)).(T)
	case []int16:
		ret = any(EnvIntSlice(name, v)).(T)
	case []int32:
		ret = any(EnvIntSlice(name, v)).(T)
	case []int64:
		ret = any(EnvIntSlice(name, v)).(T)
	case map[string]string:
		ret = any(EnvStringMap(name, v)).(T)
	}

	return ret
}

// EnvString reads environment variable and if it's not blank returns it
// otherwise returns fallback value
func EnvString(name string, fallback string) string {
	result := os.Getenv(name)
	if result == "" {
		return fallback
	}
	return result
}

// EnvStringSlice reads environment variable (comma separated) and parses them
// into a string slice (commas can not be escaped).
func EnvStringSlice(name string, fallback []string) []string {
	result := os.Getenv(name)
	if result == "" {
		return fallback
	}
	raw := strings.Split(result, ",")
	out := make([]string, 0, len(raw))
	for _, item := range raw {
		out = append(out, strings.Trim(item, "\r\n "))
	}
	return out
}

// EnvStringMap reads environment variable (with comma separated values) and parses them
// into a map[string]string.
//
//	X-TEST=me, X-Flag=true, Host: FQDN-Value =>
//	map[string]string{"X-TEST": "me", "X-Flag": "true", "Host": "FQDN-Value"}
func EnvStringMap(name string, fallback map[string]string) map[string]string {
	result := os.Getenv(name)
	if result == "" {
		return fallback
	}

	out := make(map[string]string)
	values := strings.SplitSeq(result, ",") // X-TEST=ME, HOST=test.fqdn, etc=val
	for value := range values {
		pair := strings.SplitN(strings.Trim(value, "\r\n "), "=", 2)
		if len(pair) < 2 {
			//: not applicable, skipping silently
			continue
		}
		out[pair[0]] = pair[1]
	}

	return out
}

type ints interface {
	int | int8 | int16 | int32 | int64
}

func safeInt64(input string, fallback int64) int64 {
	value, err := strconv.ParseInt(input, 10, 64)
	if err == nil {
		return value
	}
	return fallback
}

func safeInt[T ints](input string, fallback T) T {
	value, err := strconv.ParseInt(input, 10, 64)
	if err == nil {
		return T(value)
	}
	return fallback
}

// EnvIntSlice reads environment variable (comma separated) and parses them
// into a int64 slice (commas can not be escaped).
// Note, you might use the following format: "1, 13, test, -15", instead of test
// you will receive 0
func EnvIntSlice[T ints](name string, fallback []T) []T {
	result := os.Getenv(name)
	if result == "" {
		return fallback
	}
	raw := strings.Split(result, ",")
	out := make([]T, 0, len(raw))
	for _, item := range raw {
		out = append(out, safeInt(strings.Trim(item, "\r\n "), T(0)))
	}
	return out
}

// EnvInts reads environment variable, if the value can be parsed into
// int value (int8, int16, int32, int64, int), returns it, otherwise returns the fallback value.
func EnvInts[T ints](name string, fallback T) T {
	raw := os.Getenv(name)
	if raw == "" {
		return fallback
	}
	return safeInt(raw, fallback)
}

// EnvInt64 reads environment variable, if the value can be parsed into
// int64, returns it, otherwise returns the fallback value
//
// Deprecated: use EnvInts instead.
func EnvInt64(name string, fallback int64) int64 {
	raw := os.Getenv(name)
	if raw == "" {
		return fallback
	}
	return safeInt64(raw, fallback)
}

// EnvInt reads environment variable, if the value can be parsed into
// int, returns it, otherwise returns the fallback value
//
// Deprecated: use EnvInts instead.
func EnvInt(name string, fallback int) (result int) {
	if result = int(EnvInt64(name, -1)); result == -1 {
		result = fallback
	}
	return
}

// EnvInt8 reads environment variable, if the value can be parsed into
// int8, returns it, otherwise returns the fallback value
//
// Deprecated: use EnvInts instead.
func EnvInt8(name string, fallback int8) (result int8) {
	if result = int8(EnvInt64(name, -1)); result == -1 {
		result = fallback
	}
	return
}

// EnvInt16 reads environment variable, if the value can be parsed into
// int16, returns it, otherwise returns the fallback value
//
// Deprecated: use EnvInts instead.
func EnvInt16(name string, fallback int16) (result int16) {
	if result = int16(EnvInt64(name, -1)); result == -1 {
		result = fallback
	}
	return
}

// EnvInt32 reads environment variable, if the value can be parsed into
// int32, returns it, otherwise returns the fallback value
//
// Deprecated: use EnvInts instead.
func EnvInt32(name string, fallback int32) (result int32) {
	if result = int32(EnvInt64(name, -1)); result == -1 {
		result = fallback
	}
	return
}

type uints interface {
	uint | uint8 | uint16 | uint32 | uint64
}

// EnvUints reads environment variable, if the value can be parsed into
// uint values, returns it, otherwise returns the fallback value
func EnvUints[T uints](name string, fallback T) T {
	raw := os.Getenv(name)
	if raw == "" {
		return fallback
	}

	value, err := strconv.ParseUint(raw, 10, 64)
	if err == nil {
		return T(value)
	}
	return fallback
}

// EnvUint64 reads environment variable, if the value can be parsed into
// uint64, returns it, otherwise returns the fallback value
//
// Deprecated: use EnvUints instead.
func EnvUint64(name string, fallback uint64) uint64 {
	raw := os.Getenv(name)
	if raw == "" {
		return fallback
	}

	value, err := strconv.ParseUint(raw, 10, 64)
	if err == nil {
		return value
	}
	return fallback
}

// EnvUint32 reads environment variable, if the value can be parsed into
// uint32, returns it, otherwise returns the fallback value
//
// Deprecated: use EnvUints instead.
func EnvUint32(name string, fallback uint32) (result uint32) {
	if result = uint32(EnvUint64(name, 0)); result == 0 {
		result = fallback
	}
	return
}

// EnvUint16 reads environment variable, if the value can be parsed into
// uint16, returns it, otherwise returns the fallback value
//
// Deprecated: use EnvUints instead.
func EnvUint16(name string, fallback uint16) (result uint16) {
	if result = uint16(EnvUint64(name, 0)); result == 0 {
		result = fallback
	}
	return
}

// EnvUint8 reads environment variable, if the value can be parsed into
// uint8, returns it, otherwise returns the fallback value
//
// Deprecated: use EnvUints instead.
func EnvUint8(name string, fallback uint8) (result uint8) {
	if result = uint8(EnvUint64(name, 0)); result == 0 {
		result = fallback
	}
	return
}

// EnvUint reads environment variable, if the value can be parsed into
// uint, returns it, otherwise returns the fallback value
//
// Deprecated: use EnvUints instead.
func EnvUint(name string, fallback uint) (result uint) {
	if result = uint(EnvUint64(name, 0)); result == 0 {
		result = fallback
	}
	return
}

// EnvBool reads environment variable, if the value can be parsed into
// bool, returns it, otherwise returns the fallback value
func EnvBool(name string, fallback bool) bool {
	result, err := strconv.ParseBool(os.Getenv(name))
	if err != nil {
		return fallback
	}
	return result
}

type floats interface {
	float32 | float64
}

// EnvFloats reads environment variable, if the value can be parsed into
// float values, returns it, otherwise returns the fallback value
func EnvFloats[T floats](name string, fallback T) T {
	raw := os.Getenv(name)
	if raw == "" {
		return fallback
	}

	value, err := strconv.ParseFloat(raw, 64)
	if err == nil {
		return T(value)
	}
	return fallback
}

// EnvFloat reads environment variable, if the value can be parsed into
// float64, returns it, otherwise returns the fallback value
//
// Deprecated: use EnvFloats instead.
func EnvFloat(name string, fallback float64) float64 {
	raw := os.Getenv(name)
	if raw == "" {
		return fallback
	}

	value, err := strconv.ParseFloat(raw, 64)
	if err == nil {
		return value
	}
	return fallback
}

// EnvDuration reads environment variable, if the value can be parsed into
// time.Duration, returns it, otherwise returns the fallback value
func EnvDuration(name string, fallback time.Duration) time.Duration {
	rawValue := os.Getenv(name)
	if rawValue == "" {
		return fallback
	}
	if value, err := time.ParseDuration(rawValue); err == nil {
		return value
	}
	return fallback
}

// EnvUsage modifies original help string with environment usage
func EnvUsage(env string, original string) string {
	return fmt.Sprintf("%s\n[env: %s]", original, env)
}

// EnvUsageP same as [EnvUsage] but takes environment prefix helper function as first argument.
func EnvUsageP(prefixer PrefixFunc, env string, original string) string {
	return EnvUsage(prefixer(env), original)
}
