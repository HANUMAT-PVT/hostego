package middlewares

import (
	"backend-hostego/internal/app/hostego-service/constants"
	"backend-hostego/internal/app/hostego-service/dto"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

func VerifyUserAuthCookie(c fiber.Ctx) (string, error) {
	cookie := c.Cookies("auth_token")
	fmt.Print("cookie", cookie)
	if cookie == "" {
		return "", fmt.Errorf("no auth token found")
	}
	token, err := jwt.Parse(cookie, func(t *jwt.Token) (interface{}, error) {
		return []byte("hanumat"), nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		user_id, exists := claims["user_id"].(string) // Ensure correct type
		if !exists {
			return "", fmt.Errorf("user_id not found in token")
		}
		return user_id, nil
	}

	return "", fmt.Errorf("invalid token")
}

var jwtSecret = []byte("hanumat-gagan") // Ensure this is consistent across your system

// Middleware to validate user authentication token
func InternalTokenVerifivationORAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		reqCtx := dto.ReqCtx{}

		err := ValidateToken(ctx, &reqCtx)
		err2 := ValidateUserAuthToken(ctx, &reqCtx)
		if err != nil && err2 != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authentication token"})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

func ValidateToken(ctx *gin.Context, rCtx *dto.ReqCtx) error {
	tokenString := ctx.Request.Header.Get("X-Access-Token")
	if tokenString != "HJBJKNJNJNBJNBBNNB" {
		return errors.New("missing authentication token")
	}
	return nil
}

func UserAuthTokenValidation() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		reqCtx := dto.ReqCtx{}

		err := ValidateUserAuthToken(ctx, &reqCtx)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authentication token"})
			ctx.Abort()
			return
		}

		ctx.Set(string(constants.UserIdCtxKey), reqCtx.UserId)
		ctx.Next()
	}
}

// Validate the JWT token and extract user ID
func ValidateUserAuthToken(ctx *gin.Context, rCtx *dto.ReqCtx) error {
	tokenString := ctx.Request.Header.Get("X-Access-Token")
	if tokenString == "" {
		return errors.New("missing authentication token")
	}

	userId, err := getUserIdFromToken(tokenString)
	if err != nil {
		return err
	}

	rCtx.UserId = userId
	return nil
}

// Extract user_id from JWT claims securely
func getUserIdFromToken(accessToken string) (int, error) {
	claims, err := ParseUserIdAndValidityFromJWTToken(accessToken)
	if err != nil {
		return 0, err
	}
	userIdStr, ok := claims["user_id"].(string)
	if !ok {
		return 0, errors.New("invalid user ID format in token")
	}
	var userId int
	_, err = fmt.Sscanf(userIdStr, "%d", &userId)
	if err != nil {
		return 0, errors.New("error converting user_id to integer")
	}
	return userId, nil
}

func ParseUserIdAndValidityFromJWTToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok || token.Method.Alg() != "HS512" {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token claims")
}
