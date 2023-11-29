package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	env "github.com/edwinbustillos/api-go-gin/utils"
	"github.com/gin-gonic/gin"
)

func GetToken(c *gin.Context) {
	env.Load()
	secret := env.GetEnv("SECRET", "")

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  "user-id",
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	})

	token, err := claims.SignedString([]byte(secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func RefreshToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	refreshToken := strings.TrimPrefix(authHeader, "Bearer ")

	env.Load()
	secret := env.GetEnv("SECRET", "")

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		newClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":  claims["id"],
			"exp": time.Now().Add(time.Hour * 72).Unix(),
		})

		newToken, err := newClaims.SignedString([]byte(secret))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": newToken})
	}
}
