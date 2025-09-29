package controllers

import (
	"main/intializers"
	"main/models"
	"main/utils"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *gin.Context) {
	//take the info from the request body
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
//bind the data from the request body
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
	// 3. Create access token (short-lived)
	accessToken, err := utils.CreateJWT(user.ID, 15*time.Minute)
	if err != nil {
		c.JSON(500, gin.H{"error": "Could not create access token"})
		return
	}

	// 4. Create refresh token (long-lived JWT)
	refreshToken, err := utils.CreateRefreshJWT(user.ID, 7*24*time.Hour)
	if err != nil {
		c.JSON(500, gin.H{"error": "Could not create refresh token"})
		return
	}

	// 5. Set refresh token in HttpOnly cookie
	c.SetCookie("refresh_token", refreshToken, int((7 * 24 * time.Hour).Seconds()), "/", "localhost", false, true)

	c.JSON(200, gin.H{
		"access token": accessToken,
	})

}

func Logout(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(400, gin.H{
			"error": "unauthorised refrsh token has been given by the cookies",
		})
		return
	}

	// delete token from the DB
	intializers.DB.Where("token = ?", refreshToken).Delete(&models.RefreshToken{})

	//remove the cookies
	c.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	c.JSON(200, gin.H{
		"message": "logout has been successfully",
	})

}

func Validate(c *gin.Context) {

	c.JSON(200, gin.H{
		"message": "i am logged in",
	})

}
