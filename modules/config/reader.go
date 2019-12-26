// Package config ...
package config

import (
	"os"
	"path/filepath"

	"github.com/jinzhu/configor"
)

// BasicAuth ...
type BasicAuth struct {
	Name string `toml:"name" default:"Secure Portal"`
}

// Cookies ...
type Cookies struct {
	Redirect string `toml:"redirect" default:""`
	Auth     string `toml:"auth" default:"Auth-Portal"`
}

// Source ...
type Source struct {
	Host string `toml:"host" required:"true"`
}

// Auth ...
type Auth struct {
	Port      int       `toml:"port" default:"80"`
	RPCPort   int       `toml:"rpcport" default:"50051"`
	Source    Source    `toml:"source"`
	Cookies   Cookies   `toml:"headers"`
	BasicAuth BasicAuth `toml:"basicauth"`
}

// Database ...
type Database struct {
	Name     string `toml:"name" default:"app"`
	User     string `toml:"user" default:"app"`
	Password string `toml:"password" default:"app"`
	Host     string `toml:"host" default:"mysql"`
	Port     int    `toml:"port" default:"3306"`
}

// Vars ...
var Vars = struct {
	Debug    bool     `toml:"debug" default:"false"`
	Auth     Auth     `toml:"auth"`
	Database Database `toml:"database"`
}{}

// ReadDefault ...
func ReadDefault() error {
	file, err := filepath.Abs("./config.toml")
	if err != nil {
		return err
	}
	return Read(file)
}

// IsDebug ...
func IsDebug() bool {
	return os.Getenv("DEBUG") == "1" || os.Getenv("DEBUG") == "true"
}

// IsVerbose ...
func IsVerbose() bool {
	return os.Getenv("VERBOSE") == "1" || os.Getenv("VERBOSE") == "true"
}

// Read ...
func Read(file string) error {
	return configor.New(&configor.Config{ENVPrefix: "SP", Debug: IsDebug(), Verbose: IsVerbose()}).Load(&Vars, file)
}
