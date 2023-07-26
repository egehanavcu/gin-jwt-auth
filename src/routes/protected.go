package routes

import (
	controllers "gin-jwt-auth/src/controlllers"
)

func ProtectedRouter() {
	Router.GET("/protected", controllers.ProtectedController)
}
