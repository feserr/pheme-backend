// Package models config models.
package models

// ServerConfig server configuration.
type ServerConfig struct {
	IP   string `yaml:"ip"`
	Port int    `yaml:"port"`
}

// DbConfig database configuration.
type DbConfig struct {
	Host   string `yaml:"host"`
	User   string `yaml:"user"`
	Dbname string `yaml:"dbname"`
}

// Network configuration.
type Network struct {
	Server ServerConfig `yaml:"server"`
	Db     DbConfig     `yaml:"db"`
}
