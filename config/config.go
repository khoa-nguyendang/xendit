package config

import (
	"errors"
	"log"
	"time"

	"github.com/spf13/viper"
)

// App config struct
type Config struct {
	Server                ServerConfig
	Mysql                 MysqlConfig
	MysqlReplica          MysqlConfig
	Redis                 RedisConfig
	Logger                Logger
	Jaeger                Jaeger
	AuthenticationService string
}

// Server config struct
type ServerConfig struct {
	AppVersion               string
	Port                     string
	Mode                     string
	JwtSecretKey             string
	JwtExpireInHour          int
	RefreshSecretKey         string
	RefreshTokenExpireInHour int
	CookieName               string
	ReadTimeout              time.Duration
	WriteTimeout             time.Duration
	CtxDefaultTimeout        time.Duration
	Debug                    bool
	MaxConnectionIdle        time.Duration
	Timeout                  time.Duration
	MaxConnectionAge         time.Duration
	Time                     time.Duration
	CacheExpiryShort         time.Duration // default 10 mins
	CacheExpiryMedium        time.Duration // default 1 hours
	CacheExpiryLong          time.Duration // default 4 hours
	CacheExpiryDayLong       time.Duration // default 1 day
	HashKey                  string
	PassKey                  string
	IvKey                    string
}

// Logger config
type Logger struct {
	Development       bool
	DisableCaller     bool
	DisableStacktrace bool
	Encoding          string
	Level             string
}

// MySQL config
type MysqlConfig struct {
	MysqlHost     string
	MysqlPort     string
	MysqlUser     string
	MysqlPassword string
	MysqlDbname   string
	MysqlDriver   string
}

// Redis config
type RedisConfig struct {
	RedisAddr      string
	RedisPassword  string
	RedisDB        string
	RedisDefaultdb string
	MinIdleConns   int
	PoolSize       int
	PoolTimeout    int
	Password       string
	DB             int
}

// Cookie config
type Cookie struct {
	Name     string
	MaxAge   int
	Secure   bool
	HTTPOnly bool
}

// Session config
type Session struct {
	Prefix string
	Name   string
	Expire int
}

// Metrics config
type Metrics struct {
	URL         string
	ServiceName string
}

// Jaeger
type Jaeger struct {
	Host        string
	ServiceName string
	LogSpans    bool
}

// Kafka config
type Kafka struct {
	Server   string
	Username string
	Password string
}

// Load config file from given path
func LoadConfig(filename string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigName(filename)
	v.AddConfigPath(".")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	return v, nil
}

// Parse config file
func ParseConfig(v *viper.Viper) (*Config, error) {
	var c Config

	err := v.Unmarshal(&c)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}

	return &c, nil
}

// Get config
func GetConfig(configPath string) (*Config, error) {
	cfgFile, err := LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	cfg, err := ParseConfig(cfgFile)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func GetConfigPath(configPath string) string {
	if configPath == "docker" {
		return "./config/config-docker"
	}
	return "./config/config.dev"
}
