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
	r.POST("/user", controllers.Create)
	r.POST("/notes", controllers.CreateNote)
	r.GET("/notes", controllers.GetNotes)
	r.GET("/notes/:id", controllers.GetNotesbyid)
	r.PUT("/notes/:id", controllers.UpdateNote)
	r.DELETE("/notes/:id", controllers.DeleteNote)
	r.POST("/signup", controllers.SignUp)
	r.POST("/login", controllers.Login)
	r.GET("/validate", middleware.RequireAuth, controllers.Validate)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"

	}
	r.Run(":" + port)

}
