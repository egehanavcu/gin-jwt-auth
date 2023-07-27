package main

import (
	"gin-jwt-auth/src/config"
	"gin-jwt-auth/src/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.DebugMode)
	routes.Init()
	routes.Router.Run(config.Get("Port").(string))
}
