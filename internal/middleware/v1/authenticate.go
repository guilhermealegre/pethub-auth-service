package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain"
	"github.com/guilhermealegre/pethub-auth-service/api/v1/http"
	status "net/http"
	"strings"
)

type AuthenticateMiddleware struct {
	app domain.IApp
}

func NewAuthenticateMiddleware(app domain.IApp) domain.IMiddleware {
	return &AuthenticateMiddleware{
		app: app,
	}
}

func (c *AuthenticateMiddleware) RegisterMiddlewares() {
	http.GroupV1Auth.AddMiddleware(c)
}

func (c *AuthenticateMiddleware) GetHandlers() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		ValidateToken(c.app.Http().Config().JwtSecret),
	}
}

func ValidateToken(jwtSecret string) func(ctx *gin.Context) {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(status.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Check Bearer scheme
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(status.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			c.Abort()
			return
		}

		// Extract token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil {
			c.JSON(status.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if ok && token.Valid {
			c.Set("userUUID", claims["user_uuid"])
			c.Set("email", claims["email"])
		} else {
			c.JSON(status.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}
	}
}
