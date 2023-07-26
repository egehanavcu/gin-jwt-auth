package main

import (
	"gin-jwt-auth/src/config"
	"gin-jwt-auth/src/routes"

	"github.com/gin-gonic/gin"
)

var cfg *config.Config

func main() {
	gin.SetMode(gin.DebugMode)

	cfg := config.New()

	routes.Init()
	routes.Router.Run(cfg.Port)
}
