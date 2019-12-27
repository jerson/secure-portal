// Package config ...
package config

import (
	"os"
	"path/filepath"

	"github.com/jinzhu/configor"
)

// BasicAuth ...
type BasicAuth struct {
	Name     string `toml:"name" default:"Secure Portal" required:"true"`
	Username string `toml:"username" default:"admin"`
	Password string `toml:"password" default:"admin"`
}

// Cookies ...
type Cookies struct {
	Redirect string `toml:"redirect" default:"Redirect" required:"true"`
	Auth     string `toml:"auth" default:"Auth-Portal" required:"true"`
}

// AuthSource ...
type AuthSource struct {
	Host string `toml:"host" required:"true"`
}

// AuthPath ...
type AuthPath struct {
	Login             string   `toml:"login" default:"/sp/login"  required:"true"`
	Logout            string   `toml:"logout" default:"/sp/logout"  required:"true"`
	Register          string   `toml:"register" default:"/sp/register"  required:"true"`
	RegisterWhitelist []string `toml:"register_whitelist" default:"127.0.0.1"`
	LogoutRedirect    string   `toml:"logout_redirect" default:"/"  required:"true"`
}

// Auth ...
type Auth struct {
	Port      int        `toml:"port" default:"80" required:"true"`
	RPCPort   int        `toml:"rpc_port" default:"50051"`
	Type      string     `toml:"type" default:"basicauth"`
	Source    AuthSource `toml:"source"`
	Path      AuthPath   `toml:"path"`
	Cookies   Cookies    `toml:"headers"`
	BasicAuth BasicAuth  `toml:"basic_auth"`
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
