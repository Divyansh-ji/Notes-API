package controllers

import (
	"main/intializers"
	"main/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *gin.Context) {
	var users struct {
		Email    string
		Password string
	}
	if err := c.ShouldBind(&users); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(users.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}
	user := models.User{
		Email:    users.Email,
		Password: string(hashed),
	}
	if err := intializers.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}
	c.JSON(http.StatusCreated, gin.H{"user": user})

}
func Login(c *gin.Context) {

	var users struct {
		Email    string
		Password string
	}
	if err := c.ShouldBind(&users); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}
	// look up the user
	var user models.User

	intializers.DB.First(&user, "email = ?", users.Email)
	if user.ID == 0 {
		c.JSON(404, gin.H{
			"error": "Invalid user /loginx credentialas",
		})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(users.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(404, gin.H{
			"error": "Invalid user /login credentialas",
		})
		return

	}
	//send it back
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)

	c.JSON(200, gin.H{})

}
func Validate(c *gin.Context) {

	c.JSON(200, gin.H{
		"message": "i am logged in",
	})

}
