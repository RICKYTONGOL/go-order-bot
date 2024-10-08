// Package config handles server configuration binding and loading.
package config

// default configuration elements and keys
const (
	// configuration options
	keyAppName        = "app_name"
	keyConfigFilePath = "confpath"

	// logging related options
	keyLoggingLevel    = "log.level"
	keyLoggingFormat   = "log.format"
	keyLoggingFilePath = "log.filepath"

	// server related keys
	keyBindAddress      = "server.bind"
	keyDomainAddress    = "server.domain"
	keyApiPeers         = "server.peers"
	keyApiStateOrigin   = "server.origin"
	keyCorsAllowOrigins = "server.cors_origins"

	// server time out related keys
	keyTimeoutRead     = "server.read_timeout"
	keyTimeoutWrite    = "server.write_timeout"
	keyTimeoutIdle     = "server.idle_timeout"
	keyTimeoutHeader   = "server.header_timeout"
	keyTimeoutResolver = "server.resolver_timeout"

	// off-chain database related options
	keyMongoUrl      = "db.url"
	keyMongoDatabase = "db.db"

	// cache related options
	keyCacheEvictionTime = "cache.eviction"
	keyCacheMaxSize      = "cache.size"
)
