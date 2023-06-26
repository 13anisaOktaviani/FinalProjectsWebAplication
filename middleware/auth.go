package middleware

import (
	"a21hc3NpZ25tZW50/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cookie, err := ctx.Request.Cookie("session_token")
		if err != nil || cookie.Value == "" {
			if ctx.GetHeader("Content-Type") == "application/json" {
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
				ctx.Abort()
				return
			} else {
				ctx.Redirect(http.StatusSeeOther, "/login")
				ctx.Abort()
				return
			}
		}

		tokenString := cookie.Value
		claims := &model.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return model.JwtKey, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
				ctx.Abort()
				return
			}
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
			ctx.Abort()
			return
		}

		if !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			ctx.Abort()
			return
		}

		ctx.Set("email", claims.Email)

		ctx.Next()
	}
}
