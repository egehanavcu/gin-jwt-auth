package config

import "time"

var Config = map[string]interface{}{
	"Port":                     ":8080",
	"JWTSecret":                "your-256-bit-secret",
	"AccessTokenIdentifierKey": "email",
	"AccessTokenDuration":      time.Minute * 1,
	"RefreshTokenDuration":     time.Hour * 24,
}

func Get(key string) interface{} {
	return Config[key]
}
