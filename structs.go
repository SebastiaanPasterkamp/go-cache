package cache

// Configuration offers settings for all cache.Repository implementations.
type Configuration struct {
	*InMemorySettings
	*RedisSettings
}

// InMemorySettings contains all settings for the in-memory cache implementation.
type InMemorySettings struct {
}

// RedisSettings contains all settings for the redis based cache implementation.
type RedisSettings struct {
	Address  string `json:"address" arg:"--redis-address,env:REDIS_ADDRESS" placeholder:"host:port" help:"The address of the redis cache service."`
	Password string `json:"password" arg:"--redis-password,env:REDIS_PASSWORD" help:"The redis password to use."`
	Database int    `json:"database" arg:"--redis-database,env:REDIS_DATABASE" placeholder:"#" help:"The redis database to use."`
}
