package tools

import (
	"testing"
)

func TestMakeEnv(t *testing.T) {
	for _, entry := range []struct {
		name     string
		envs     []string
		in       string
		expected string
	}{
		{
			name:     "simple",
			envs:     []string{"APP"},
			in:       "TEST",
			expected: "APP_TEST",
		},
		{
			name:     "complex",
			envs:     []string{"THIS", "IS"},
			in:       "COMPLEX_ENV_VALUE",
			expected: "THIS_IS_COMPLEX_ENV_VALUE",
		},
	} {
		t.Run(entry.name, func(in *testing.T) {
			if result := MakeEnv(entry.envs...)(entry.in); result != entry.expected {
				in.Errorf("expected: %v, got: %v", entry.expected, result)
			}
		})
	}
}

func TestEnvJoin(t *testing.T) {
	for _, entry := range []struct {
		name     string
		in       string
		envs     []string
		expected string
	}{
		{
			name:     "simple",
			envs:     []string{"APP"},
			in:       "TEST",
			expected: "APP_TEST",
		},
		{
			name:     "complex",
			envs:     []string{"THIS", "IS"},
			in:       "COMPLEX_ENV_VALUE",
			expected: "THIS_IS_COMPLEX_ENV_VALUE",
		},
	} {
		t.Run(entry.name, func(in *testing.T) {
			if result := EnvJoin(entry.in, entry.envs...); result != entry.expected {
				in.Errorf("expected: %v, got: %v", entry.expected, result)
			}
		})
	}
}

func TestSensitive(t *testing.T) {
	t.Run("ok", func(in *testing.T) {
		p := &Sensitive{}
		if err := p.Set("my-custom-password"); err != nil {
			in.Errorf("got error: %v", err)
		}

		expected := "[hidden]"
		if result := p.String(); result != expected {
			in.Errorf("expected: `%v`, got: `%v`", expected, result)
		}

		expected = "my-custom-password"
		if result := p.Reveal(); result != expected {
			in.Errorf("expected: `%v`, got: `%v`", expected, result)
		}

		expected = "Sensitive"
		if result := p.Type(); result != expected {
			in.Errorf("expected: `%v`, got: `%v`", expected, result)
		}
	})
}
