package config

import (
	"time"

	"github.com/spf13/viper"
)

// Default values of configuration options
const (
	// config file name
	configFileName = "conf"
)

var (
	// defApplicationName defines default application name
	defApplicationName = "API Server"

	// defLoggingLevel holds default Logging level
	// See `godoc.org/github.com/op/go-logging` for the full format specification
	// See `golang.org/pkg/time/` for time format specification
	defLoggingLevel = "INFO"

	// defLoggingFormat holds default format of the Logger output
	defLoggingFormat = "%{color}%{time:15:04:05.000} ▶ %{level:.4s} %{message} from=%{shortpkg}/%{shortfunc}%{color:reset}:"

	// defLoggingFilePath holds default log file path
	defLoggingFilePath = "orderbot.log"

	// defServerBind holds default API server binding address
	defServerBind = "0.0.0.0:9095"

	// defServerDomain holds default API server domain address
	defServerDomain = "0.0.0.0:9095"

	// default set of timeouts for the server
	defReadTimeout     = 2
	defWriteTimeout    = 15
	defIdleTimeout     = 1
	defHeaderTimeout   = 1
	defResolverTimeout = 30

	// defMongoUrl holds default MongoDB connection string
	defMongoUrl = "mongodb://127.0.0.1:27017"

	// defMongoDatabase holds the default name of the API persistent database
	defMongoDatabase = "dongleapi"

	// defCacheEvictionTime holds default time for in-memory eviction periods
	defCacheEvictionTime = 15 * time.Minute

	// defCacheMax size represents the default max size of the cache in MB
	defCacheMaxSize = 4096

	// defCorsAllowOrigins represents the cors allow orginins
	defCorsAllowOrigins = []string{"*"}
)

// applyDefaults sets default values for configuration options.
func applyDefaults(cfg *viper.Viper) {
	cfg.SetDefault(keyAppName, defApplicationName)

	cfg.SetDefault(keyLoggingLevel, defLoggingLevel)
	cfg.SetDefault(keyLoggingFormat, defLoggingFormat)
	cfg.SetDefault(keyLoggingFilePath, defLoggingFilePath)

	cfg.SetDefault(keyBindAddress, defServerBind)
	cfg.SetDefault(keyDomainAddress, defServerDomain)

	cfg.SetDefault(keyCorsAllowOrigins, defCorsAllowOrigins)

	// server timeouts
	cfg.SetDefault(keyTimeoutRead, defReadTimeout)
	cfg.SetDefault(keyTimeoutWrite, defWriteTimeout)
	cfg.SetDefault(keyTimeoutHeader, defHeaderTimeout)
	cfg.SetDefault(keyTimeoutIdle, defIdleTimeout)
	cfg.SetDefault(keyTimeoutResolver, defResolverTimeout)

	// DB Settings
	cfg.SetDefault(keyMongoUrl, defMongoUrl)
	cfg.SetDefault(keyMongoDatabase, defMongoDatabase)

	// in-memory cache
	cfg.SetDefault(keyCacheEvictionTime, defCacheEvictionTime)
	cfg.SetDefault(keyCacheMaxSize, defCacheMaxSize)

}
