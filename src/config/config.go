package config

type Config struct {
	Port      string
	JWTSecret string
}

func New() *Config {
	return &Config{
		Port:      ":8080",
		JWTSecret: "your-256-bit-secret",
	}
}
