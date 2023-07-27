package routes

import (
	controllers "gin-jwt-auth/src/controlllers"
	"gin-jwt-auth/src/middlewares"
)

func ProtectedRouter() {
	Router.GET("/protected/admin", middlewares.JWTMiddleware("admin"), controllers.ProtectedController)
	Router.GET("/protected/user", middlewares.JWTMiddleware("user"), controllers.ProtectedController)
	Router.GET("/protected/all", middlewares.JWTMiddleware("admin", "user"), controllers.ProtectedController)
}
