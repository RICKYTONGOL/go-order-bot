// Package config handles API server configuration binding and loading.
package config

import "time"

// Config defines configuration options structure for server.
type Config struct {
	// AppName holds the name of the application
	AppName string `mapstructure:"app_name"`

	// Logger configuration
	Log Log `mapstructure:"log"`

	// Server configuration
	Server Server `mapstructure:"server"`

	// Database configuration
	Db Database `mapstructure:"db"`

	// Websocket Server configuration
	WSServer WSServer `mapstructure:"ws"`

	// Cache configuration
	Cache Cache `mapstructure:"cache"`

	// Twilio access configuration
	TwilioInfo TwilioInfo `mapstructure:"twilioinfo"`
}

// Log represents the logger configuration
type Log struct {
	Level    string `mapstructure:"level"`
	Format   string `mapstructure:"format"`
	FilePath string `mapstructure:"filepath"`
}

// Server represents the GraphQL server configuration
type Server struct {
	BindAddress     string   `mapstructure:"bind"`
	DomainAddress   string   `mapstructure:"domain"`
	Origin          string   `mapstructure:"origin"`
	Peers           []string `mapstructure:"peers"`
	CorsOrigin      []string `mapstructure:"cors_origins"`
	ReadTimeout     int64    `mapstructure:"read_timeout"`
	WriteTimeout    int64    `mapstructure:"write_timeout"`
	IdleTimeout     int64    `mapstructure:"idle_timeout"`
	HeaderTimeout   int64    `mapstructure:"header_timeout"`
	ResolverTimeout int64    `mapstructure:"resolver_timeout"`
}

type WSServer struct {
	BindAddress string `mapstructure:"bind"`
}

// Database represents the database access configuration.
type Database struct {
	Url    string `mapstructure:"url"`
	DbName string `mapstructure:"db"`
}

// Cache represents the cache sub-system configuration.
type Cache struct {
	Eviction time.Duration `mapstructure:"eviction"`
	MaxSize  int           `mapstructure:"size"`
}

// TwilioInfo represents the twilio access information.
type TwilioInfo struct {
	AccountSid string `mapstructure:"accountsid"`
	AuthToken  string `mapstructure:"authtoken"`
}
