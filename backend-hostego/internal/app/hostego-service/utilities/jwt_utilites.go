package utilities

import (
	"backend-hostego/internal/app/hostego-service/models"
	"backend-hostego/internal/pkg/logger"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var log = logger.GetLogger()
var jwtSecret = []byte("hanumat-gagan")

func ParseUserIdAndValidityFromJWTToken(ctx *gin.Context, tokenString string) (jwt.MapClaims, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		log.Errorln("ERROR : while executing ParseUnverified", err)
		return nil, errors.New("Error while parsing token.")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	} else {
		return nil, nil
	}
}

func GenerateJWT(user models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":    user.UserId,
		"first_name": user.FirstName,
		"email":      user.Email,
		"mobile":     user.MobileNumber,
		"exp":        time.Now().Add(24 * 30 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
