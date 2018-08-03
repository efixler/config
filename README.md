# config
Go package to provide loosely coupled configuration values to your app. 

[![Go Documentation](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)][godocs]

[godocs]: https://godoc.org/github.com/efixler/config

Out-of-the-box, config provides a simple passthrough for environment variables. With a little bit of code you can do custom configuration settings based on the deploy env, or arbitrarily complex implementations. In all of these cases, packages that consume configuration only need to import the config package, and not the concrete data provider.

## Installation

`go get github.com/efixler/config`

## Usage

````
import (
	"github.com/efixler/config"
)

cfg := config.Default()
apiHost := cfg.GetOrDefault("API_HOST", "https://api.foobar.com")

 ````

See the [Godoc](https://godoc.org/github.com/efixler/config) for details and more examples. 
