package main

import (
	"fmt"
	"main/controllers"
	"main/intializers"
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"

	}
	r.Run(":" + port)

}
