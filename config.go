// Package config provides a way to to supply runtime configuration values in a loosely coupled way.
// You can use config to read the environment, or to provide an arbitrarily complex configuration mechanism,
// without binding your application tightly to that mechanism.
//
// As provided, config.Default() will return a Getter that provides access to values from the environment.
// Implement your own Loader to mutate the environment before other packages consume configuration values, or
// implement a custom Getter to use some other configuration scheme.
//
// Several packages under efixler are dependent on the config package. 
package config

import (
	"context"
	"log"
	"os"
	"strings"
	"sync"
)

var (
	loader      Loader
	loadLock    sync.Mutex
	defaultConf Getter
)

// Core interface for implementations providing configuration data to consumers.
type Getter interface {
	Get(string) string
	GetOrDefault(string, string) string
	GetStrings(string) []string
	MustGet(string) string
}

// Loader is a callback function that is be used to delegate configuration loading.
// Loader is expected to return a Getter that will be used by consumers to actually get the configuration 
// values.
// If the Loader copies configuration values into the environment, it can use
// 	Environment() 
// to return a Getter that gets from the environment. See the example for SetLoader.
//
// The Context argument defined to provide a hook for per-request mutated configs (which
// is not yet implemented). Calls to the Loader function can expect to receive a nil Context.
type Loader func(context.Context) Getter

// Pass a Loader that will be utilized to supply the Getter returned by Default(). 
// The Loader function will be called upon the first call to Default() after SetLoader().
// Each call to SetLoader will clear the default Getter. 
//
// A Loader can write values into the environment and return an Enviroment() Getter, or it can
// return anything that implements the Getter interface. The example shows the former case.
// 
// If you never call SetLoader() (or if you pass a nil value), the Default() Getter will pass through
// environment variables.
// 
// SetLoader() should be called early, preferably in an init() method as close as possible to the application's
// entry point, to ensure that consumers get the right configuration as they are initializing.
func SetLoader(cl Loader) {
	loader = cl
	defaultConf = nil
}

// Return the default configuration.
func Default() Getter {
	if defaultConf != nil {
		return defaultConf
	}
	loadLock.Lock()
	defer loadLock.Unlock()
	if defaultConf != nil {
		return defaultConf
	}
	if loader != nil {
		defaultConf = loader(nil)
	} else {
		defaultConf = Environment()
	}
	return defaultConf
}

// Env is a Getter implementation that reads from the environment.
type Env struct {
}

// Return a new Env Getter.
func Environment() Getter {
	return &Env{}
}

// Equivalent to os.Getenv(key). Note that other Get-ish methods in Env call Env.Get() (and not os.Getenv)
func (e *Env) Get(key string) string {
	return os.Getenv(key)
}

// If the requested key is not present or empty, return the dflt.
func (e *Env) GetOrDefault(key string, dflt string) string {
	rval := e.Get(key)
	if rval == "" {
		return dflt
	}
	return rval
}

// GetStrings will treat a comma-delimited config value as an []string, stripping whitespace around the commas.
func (e *Env) GetStrings(key string) []string {
	rval := strings.Split(e.Get(key), ",")
	for i, val := range rval {
		rval[i] = strings.TrimSpace(val)
	}
	return rval
}

// MustGet will panic if the key is not present or empty. Use this only when you really must get.
func (e *Env) MustGet(key string) string {
	v := e.Get(key)
	if v == "" {
		log.Panicf("%s environment variable not set.", key)
	}
	return v
}
