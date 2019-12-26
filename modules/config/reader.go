// Package config ...
package config

import (
	"path/filepath"

	"github.com/jinzhu/configor"
)

// Server ...
type Server struct {
	Port    int `toml:"port" default:"8000"`
	RPCPort int `toml:"rpc_port" default:"50051"`
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
	Name     string   `toml:"name" default:"App"`
	Debug    bool     `toml:"debug" default:"false"`
	Version  string   `toml:"version" default:"latest"`
	Server   Server   `toml:"server"`
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

// Read ...
func Read(file string) error {
	return configor.New(&configor.Config{ENVPrefix: "APP", Debug: false, Verbose: false}).Load(&Vars, file)
}
