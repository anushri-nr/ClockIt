package services

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
)

var jwtSecret []byte
var jwtExpiry = time.Hour * 24

const ContextEmployeeID = "employee_id"
const ContextEmployeeRole = "employee_role"

func init() {
	// try to load .env
	if err := godotenv.Load(); err != nil {
		log.Printf("auth init: could not load .env: %v", err)
	}

	s := os.Getenv("SECRET_KEY")
	if s == "" {
		// fail fast — require a secret key to be configured
		log.Fatalf("SECRET_KEY environment variable is not set.")
	}
	jwtSecret = []byte(s)
}

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing Authorization header"})
			return
		}

		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid Authorization header format"})
			return
		}

		tokenStr := parts[1]
		claims, err := ParseToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		c.Set(ContextEmployeeID, claims.EmployeeID)
		c.Set(ContextEmployeeRole, claims.Role)
		c.Next()
	}
}

type JWTClaims struct {
	EmployeeID uint   `json:"employee_id"`
	Role       string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken creates a signed JWT for the given employee id and role
func GenerateToken(employeeID uint, role string) (string, error) {
	claims := JWTClaims{
		EmployeeID: employeeID,
		Role:       role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseToken validates the token string and returns the claims
func ParseToken(tokenStr string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrTokenInvalidClaims
}
