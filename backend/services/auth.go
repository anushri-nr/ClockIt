package services

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"clockit/backend/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret []byte
var jwtExpiry = time.Hour * 24

const ContextEmployeeID = "employee_id"
const ContextEmployeeRole = "employee_role"

var ErrAuthNotInitialized = errors.New("auth not initialized; call services.InitAuth(secret) before using auth functions")

// InitAuth initializes package-level auth state
func InitAuth(secret string) error {
	if secret == "" {
		return ErrAuthNotInitialized
	}
	jwtSecret = []byte(secret)
	return nil
}

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if jwtSecret == nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "auth not initialized"})
			return
		}
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

// SupervisorAuthorizationMiddleware ensures the caller is a supervisor and if the route
// contains :employee_id, that it matches the authenticated employee ID.
func SupervisorAuthorizationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// role check
		r, ok := c.Get(ContextEmployeeRole)
		rs, _ := r.(string)
		if !ok || rs != string(models.RoleSupervisor) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden: supervisor access required"})
			return
		}

		// optional path param check
		pid := c.Param("employee_id")
		if pid != "" {
			p64, err := strconv.ParseUint(pid, 10, 64)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid employee_id"})
				return
			}
			authIDVal, ok := c.Get(ContextEmployeeID)
			if !ok {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "missing authenticated employee id"})
				return
			}
			authID, _ := authIDVal.(uint)
			if uint(p64) != authID {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden: cannot access shifts for this employee_id"})
				return
			}
		}

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
	if jwtSecret == nil {
		return "", ErrAuthNotInitialized
	}

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
	if jwtSecret == nil {
		return nil, ErrAuthNotInitialized
	}

	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Ensure that only HS256-signed tokens are accepted
		if token.Method != jwt.SigningMethodHS256 {
			return nil, jwt.ErrTokenSignatureInvalid
		}
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
