package middleware

import (
	"log"
	"main/intializers"
	"main/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(c *gin.Context) {
	// get the cookies off req

	tokenString, err := c.Cookie("refresh_token")

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	//Decode/validat it

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {

		return []byte(os.Getenv("SECRET")), nil
	}, 
	jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		log.Fatal(err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {

		//check the exp
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		//find the User with token sub

		var user models.User

		intializers.DB.First(&user, claims["sub"])

		if user.ID == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		//Attach to req
		c.Set("user", user)

		//Continue
		c.Next()

	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

}
