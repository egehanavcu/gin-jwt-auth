package routes

import controllers "gin-jwt-auth/src/controlllers"

func AuthRouter() {
	authGroup := Router.Group("/auth")
	{
		authGroup.POST("/login", controllers.LoginHandler)
		authGroup.POST("/register", controllers.RegisterHandler)
		authGroup.GET("/logout", controllers.LogoutHandler)
	}
}
