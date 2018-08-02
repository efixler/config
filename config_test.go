package config

import (
	"context"
	"os"
	"strings"
	"testing"
)

func TestBaseDefault(t *testing.T) {
	c := Default()
	if c.Get("CONFIG_TEST") != "1" {
		t.Errorf("Error checking CONFIG_TEST key; expected '1', got '%s'", c.Get("CONFIG_TEST"))
	}
}

func TestEnvironment(t *testing.T) {
	c := Environment()
	if c.Get("CONFIG_TEST") != "1" {
		t.Errorf("Error checking CONFIG_TEST key; expected '1', got '%s'", c.Get("CONFIG_TEST"))
	}
}

func TestCustomEnvWriter(t *testing.T) {
	f := func(context.Context) Getter {
		os.Setenv("LOADED_STUFF", "1")
		return Environment()
	}
	SetLoader(f)
	if Default().Get("LOADED_STUFF") != "1" {
		t.Error("Custom loader failed; env changes not picked up")
	}
}

type fullyCustomProvider struct {
	*Env
}

func (*fullyCustomProvider) Get(key string) string {
	return strings.ToLower(key)
}

func TestFullyCustomImplementation(t *testing.T) {
	f := func(context.Context) Getter {
		return &fullyCustomProvider{}
	}
	SetLoader(f)
	c := Default()
	if c.Get("HELLO") != "hello" {
		t.Errorf("Custom provider: expected 'hello', got '%s' -- type should be *fullyCustomProvider, is %T",
			c.Get("HELLO"), c)
	}
}

func TestMain(m *testing.M) {
	os.Setenv("CONFIG_TEST", "1")
	os.Exit(m.Run())
}
