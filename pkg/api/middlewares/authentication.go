// Package middlwares contains the middlewares used in the server.
package middlewares

import (
	"net/http"
	"time"

	"github.com/codescalersinternships/secretnote-api-spa-eyadhussein/pkg/storage"
	"github.com/codescalersinternships/secretnote-api-spa-eyadhussein/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JwtAuthMiddleware is a middleware that verifies the user's token and sets the user in the context.
func JwtAuthMiddleware(store storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		username, err := c.Cookie("user")
		if err != nil {
			c.JSON(http.StatusUnauthorized, util.NewResponseError(
				util.ErrUnauthorized, http.StatusUnauthorized,
			))
			c.Abort()
			return
		}
		user, err := store.GetUserByUsername(username)
		if err != nil {
			c.JSON(http.StatusNotFound, util.NewResponseError(
				err, http.StatusNotFound,
			))
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

// CreateToken creates a new token with the given username and duration.
func CreateToken(username string, duration time.Duration, secretKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(duration).Unix(),
		})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// VerifyToken verifies the user's token.
// @Summary Verify token
// @Description Verify token
// @Tags auth
// @Success 200 {object} api.SuccessResponse
// @Failure 401 {object} util.ResponseError "Unauthorized"
// @Router /auth/verify-token [post]
// @Security Token
func VerifyToken(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenCookie, err := c.Cookie("token")

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token cookie"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenCookie, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		c.Next()
	}
}
