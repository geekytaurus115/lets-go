package handler

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func generateToken(username, userType string) (string, error) {
	secretKey := []byte(SECRET_KEY)

	// define the claims
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 1).Unix(), // token expiration time
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// sign the token with the secret key
	generatedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return generatedToken, nil
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		secretKey := []byte(SECRET_KEY)

		tokenString := c.GetHeader("Authorization")
		if strings.TrimSpace(tokenString) == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		log.Println("tokenString--> ", tokenString)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, _ := token.Claims.(jwt.MapClaims)
		username := fmt.Sprintf("%v", claims["username"])
		log.Println("Token to username--> ", username)

		// fetch user-type
		userType, err := GetUserTypeByUserName(username)
		if err != nil {
			log.Println("Error fetching user-type", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			c.Abort()
			return
		}

		// passing user-type through context
		c.Set("user_type", userType)

		c.Next()
	}
}
