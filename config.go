// The config package provides a basic interface for a runtime configuration reader,
// a basic implementation of same that reads for the environment, and a mechanism to
// load custom configuration providers without injecting project-specific dependencies 
// into generically useful packages.
// 
// Complete documentation to follow
//
package config

import (
	"context"
	"log"
	"os"
	"strings"
	"sync"
)

var (
	loader      ConfigLoader
	loadLock    sync.Mutex
	defaultConf Getter
)

type Getter interface {
	Get(string) string
	GetOrDefault(string, string) string
	GetStrings(string) []string
	MustGet(string) string
}

type ConfigLoader func(context.Context) Getter

func SetLoader(cl ConfigLoader) {
	loader = cl
	defaultConf = nil
}

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

type Env struct {
}

func Environment() Getter {
	return &Env{}
}

func (e *Env) Get(key string) string {
	return os.Getenv(key)
}

func (e *Env) GetOrDefault(key string, dflt string) string {
	rval := e.Get(key)
	if rval == "" {
		return dflt
	}
	return rval
}

func (e *Env) GetStrings(key string) []string {
	rval := strings.Split(e.Get(key), ",")
	for i, val := range rval {
		rval[i] = strings.TrimSpace(val)
	}
	return rval
}

func (e *Env) MustGet(key string) string {
	v := e.Get(key)
	if v == "" {
		log.Panicf("%s environment variable not set.", key)
	}
	return v
}
