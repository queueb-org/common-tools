package tools

import (
	"flag"
	"fmt"
	"strings"
)

// MakeEnv provides a helper to join environment variable prefixes together
// and return one argument function as a helper.
func MakeEnv(envPrefixes ...string) func(string) string {
	return func(flag string) string {
		prefix := strings.Join(envPrefixes, "_")
		return fmt.Sprintf("%s_%s", prefix, flagNormalize(flag))
	}
}

// EnvJoin adds tail to the end of envPrefixes and calls Env
func EnvJoin(tail string, envPrefixes ...string) string {
	new := make([]string, 0, len(envPrefixes)+1)
	new = append(new, envPrefixes...)
	new = append(new, tail)
	return strings.Join(new, "_")
}

func flagNormalize(rawFlag string) string {
	return strings.ToUpper(strings.ReplaceAll(rawFlag, "-", "_"))
}

// should comply with github.com/spf13/pflag.Value interface
var _ flag.Value = &Sensitive{}

// Sensitive might be used to hide real value in representation.
type Sensitive struct {
	container string
}

// String implements [flag.Value] inteface, returns redacted value instead the real one.
func (v *Sensitive) String() string {
	return "[hidden]"
}

// Set implements [flag.Value] interface, sets sensitive string data.
func (v *Sensitive) Set(in string) (err error) {
	v.container = in
	return
}

// Type implements [pflag.Value], returns internal type name.
func (v *Sensitive) Type() string {
	return "Sensitive"
}

// Reveal returns unmodified value.
func (v *Sensitive) Reveal() string {
	return v.container
}
