package controllers

import (
	"main/intializers"
	"main/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Create(c *gin.Context) {
	var NewUser models.User

	if err := c.ShouldBind(&NewUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	postuser := models.User{
		Email: NewUser.Email,
		ID:    NewUser.ID,
		Name:  NewUser.Name,
		Notes: NewUser.Notes,
	}

	intializers.DB.Create(&postuser)
	c.JSON(http.StatusOK, gin.H{
		"user": postuser,
	})

}
func CreateNote(c *gin.Context) {
	var note models.Note

	// Bind JSON
	if err := c.ShouldBindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Make sure UserID is provided and valid
	if note.UserID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "UserID is required"})
		return
	}

	// Save note to DB
	if err := intializers.DB.Create(&note).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"note": note})
}

func GetNotes(c *gin.Context) {
	var notes []models.Note
	if err := intializers.DB.Find(&notes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"notes": notes})
}

func GetNotesbyid(c *gin.Context) {
	id := c.Param("id")
	idint, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}
	var notes []models.Note

	if err := intializers.DB.Find(&notes, idint).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"notes": notes})
}
func UpdateNote(c *gin.Context) {
	id := c.Param("id")
	idint, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
		return
	}

	// Bind JSON from request
	var input struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find the note
	var note models.Note
	if err := intializers.DB.First(&note, idint).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}

	// Update fields
	note.Title = input.Title
	note.Content = input.Content
	intializers.DB.Save(&note)

	c.JSON(http.StatusOK, gin.H{"note": note})
}
func DeleteNote(c *gin.Context) {
	id := c.Param("id")
	idint, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
		return

	}
	var note models.Note
	if err := intializers.DB.First(&note, idint).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}
	intializers.DB.Delete(&note)
	c.JSON(http.StatusOK, gin.H{"note": note})

}
func DeletingUser(c *gin.Context) {
	id := c.Param("id")
	idint, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
		return

	}
	var users models.User
	if err := intializers.DB.First(&users, idint).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}
	intializers.DB.Delete(&users)
	c.JSON(http.StatusOK, gin.H{"note": users})

}
