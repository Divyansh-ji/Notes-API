package main

import (
	"fmt"

	"main/controllers"
	"main/intializers"
	"main/middleware"
	"os"

	"github.com/gin-gonic/gin"
)

func init() {
	fmt.Println("Initializing database")
	intializers.LoadEnvVariables()
	intializers.ConnectToDB()

	intializers.SyncDataBase()

}
func main() {
	fmt.Println("welcome to notes app")
	r := gin.Default()
	//r.POST("/user", controllers.Create)
	//r.POST("/notes", controllers.CreateNote)
	//r.GET("/notes", controllers.GetNotes)
	//r.GET("/notes/:id", controllers.GetNotesbyid)
	//r.PUT("/notes/:id", controllers.UpdateNote)
	//r.DELETE("/notes/:id", controllers.DeleteNote)
	//r.DELETE("/user/:id", controllers.DeletingUser)
	r.POST("/signup", controllers.SignUp)
	r.POST("/login", controllers.Login)
	r.POST("/refreshToken", controllers.RefreshToken)
	r.GET("/logout", controllers.Logout)
	// Public endpoints (if any)
	r.POST("/user", controllers.Create)

	authGroup := r.Group("/api")
	authGroup.Use(middleware.RequireAuth)
	{
		//authGroup.POST("/user", controllers.Create)
		authGroup.POST("/notes", controllers.CreateNote)
		authGroup.GET("/notes", controllers.GetNotes)
		authGroup.GET("/notes/:id", controllers.GetNotesbyid)
		authGroup.PUT("/notes/:id", controllers.UpdateNote)
		authGroup.DELETE("/notes/:id", controllers.DeleteNote)
		authGroup.DELETE("/user/:id", controllers.DeletingUser)

	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"

	}
	r.Run(":" + port)

}
