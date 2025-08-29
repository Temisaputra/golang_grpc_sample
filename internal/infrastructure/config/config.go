package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	HTTPPort          string `envconfig:"HTTP_PORT"`
	DBType            string `envconfig:"DB_TYPE"`
	DBHost            string `envconfig:"DB_HOST"`
	DBPort            string `envconfig:"DB_PORT"`
	DBUsername        string `envconfig:"DB_USERNAME"`
	DBPassword        string `envconfig:"DB_PASSWORD"`
	DBName            string `envconfig:"DB_NAME"`
	DBMaxIdleConns    int    `envconfig:"DB_MAX_IDLE_CONNS"`
	DBMaxOpenConns    int    `envconfig:"DB_MAX_OPEN_CONNS"`
	DBMaxIdleTime     int    `envconfig:"DB_MAX_IDLE_TIME"`
	DBMaxLifetime     int    `envconfig:"DB_MAX_LIFETIME"`
	AppName           string `envconfig:"APP_NAME" default:"be-inventory"`
	JWTSecret         string `envconfig:"JWT_SECRET"`
	JWTExpiration     string `envconfig:"JWT_EXPIRATION" default:"1h"`
	DisableStacktrace bool   `envconfig:"DISABLE_STACKTRACE" default:"false"`
	Env               string `envconfig:"ENV" default:"development"`
	GRPCPort          string `envconfig:"GRPC_PORT" default:"50051"`
}

func Get() *Config {
	cfg := Config{}
	envconfig.MustProcess("", &cfg)
	return &cfg
}
