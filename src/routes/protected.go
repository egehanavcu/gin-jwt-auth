package routes

import (
	controllers "gin-jwt-auth/src/controlllers"
	"gin-jwt-auth/src/middlewares"
)

func ProtectedRouter() {
	Router.GET("/protected/admin", middlewares.AuthMiddleware("admin"), controllers.ProtectedController)
	Router.GET("/protected/user", middlewares.AuthMiddleware("user"), controllers.ProtectedController)
	Router.GET("/protected/all", middlewares.AuthMiddleware("admin", "user"), controllers.ProtectedController)
}
