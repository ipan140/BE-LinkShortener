package config

import (
	"time"
)

type AppConfig struct {
	App      ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Worker   WorkerConfig
}

type ServerConfig struct {
	Name        string
	Port        string
	Env         string
	BaseURL     string
	ReadTimeout time.Duration
}

type DatabaseConfig struct {
	Driver          string // "postgres" or "sqlite"
	DSN             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type RedisConfig struct {
	Host         string
	Port         string
	Password     string
	DB           int
	PoolSize     int
	MinIdleConns int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	StreamKey    string
	QueueKey     string
}

type JWTConfig struct {
	Secret        string
	ExpireMinutes int
}

type WorkerConfig struct {
	BatchSize     int
	FlushInterval time.Duration
	WorkerCount   int
}

// LoadConfig initializes application configuration from environment variables
func LoadConfig() *AppConfig {
	LoadEnv()

	return &AppConfig{
		App: ServerConfig{
			Name:        GetEnv("APP_NAME", "UltraFast-Redirector"),
			Port:        GetEnv("APP_PORT", "8080"),
			Env:         GetEnv("APP_ENV", "development"),
			BaseURL:     GetEnv("BASE_URL", "http://localhost:8080"),
			ReadTimeout: time.Duration(GetEnvInt("SERVER_READ_TIMEOUT_SEC", 5)) * time.Second,
		},
		Database: DatabaseConfig{
			Driver:          GetEnv("DB_DRIVER", "postgres"),
			DSN:             GetEnv("DB_DSN", "host=localhost user=postgres password=Cebes1234 dbname=shortener_db port=5432 sslmode=disable TimeZone=Asia/Jakarta"),
			MaxOpenConns:    GetEnvInt("DB_MAX_OPEN_CONNS", 50),
			MaxIdleConns:    GetEnvInt("DB_MAX_IDLE_CONNS", 25),
			ConnMaxLifetime: time.Duration(GetEnvInt("DB_CONN_MAX_LIFETIME_MIN", 30)) * time.Minute,
		},
		Redis: RedisConfig{
			Host:         GetEnv("REDIS_HOST", "localhost"),
			Port:         GetEnv("REDIS_PORT", "6379"),
			Password:     GetEnv("REDIS_PASSWORD", ""),
			DB:           GetEnvInt("REDIS_DB", 0),
			PoolSize:     GetEnvInt("REDIS_POOL_SIZE", 200),
			MinIdleConns: GetEnvInt("REDIS_MIN_IDLE_CONNS", 50),
			DialTimeout:  time.Duration(GetEnvInt("REDIS_DIAL_TIMEOUT_MS", 2000)) * time.Millisecond,
			ReadTimeout:  time.Duration(GetEnvInt("REDIS_READ_TIMEOUT_MS", 1000)) * time.Millisecond,
			WriteTimeout: time.Duration(GetEnvInt("REDIS_WRITE_TIMEOUT_MS", 1000)) * time.Millisecond,
			StreamKey:    GetEnv("REDIS_STREAM_KEY", "stream:clicks"),
			QueueKey:     GetEnv("REDIS_QUEUE_KEY", "queue:clicks"),
		},
		JWT: JWTConfig{
			Secret:        GetEnv("JWT_SECRET", "super-secret-high-performance-key-change-in-prod"),
			ExpireMinutes: GetEnvInt("JWT_EXPIRE_MINUTES", 1440), // 24 hours
		},
		Worker: WorkerConfig{
			BatchSize:     GetEnvInt("WORKER_BATCH_SIZE", 100),
			FlushInterval: time.Duration(GetEnvInt("WORKER_FLUSH_INTERVAL_MS", 1000)) * time.Millisecond,
			WorkerCount:   GetEnvInt("WORKER_COUNT", 4),
		},
	}
}
