package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		secretKey := []byte(SECRET_KEY)

		tokenString := r.Header.Get("Authorization")
		if strings.TrimSpace(tokenString) == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		log.Println("tokenString", tokenString)

		jwttoken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})

		log.Println("token--> ", jwttoken)
		claims, _ := jwttoken.Claims.(jwt.MapClaims)
		userName := fmt.Sprintf("%v", claims["username"])
		log.Println("Token to username---> ", userName)

		// fetch user-type
		userType, err := GetUserTypeByUserName(userName)
		if err != nil {
			log.Println("Error fetching user-type", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// pass user-type through context
		ctx := context.WithValue(r.Context(), "user_type", userType)

		if err != nil || !jwttoken.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

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
