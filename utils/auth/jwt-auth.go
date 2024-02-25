package auth

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type JWTClaim struct {
	Authorized bool
	ID         int    `json:"id"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	jwt.StandardClaims
}

func GenerateToken(id int, email string, username string) (string, error) {

	expirationTime, err := strconv.Atoi(os.Getenv("TOKEN_HOUR_LIFESPAN"))

	if err != nil {
		return "", err
	}

	claims := &JWTClaim{
		Authorized: true,
		ID:         id,
		Email:      email,
		Username:   username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(expirationTime)).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("API_SECRET")))

}

func TokenValidation(inputToken string) (int, string, error) {

	token, err := jwt.ParseWithClaims(
		inputToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte([]byte(os.Getenv("API_SECRET"))), nil
		},
	)
	if err != nil {
		err = errors.New("invalid access token")
		return 0, "", err
	}
	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		err = errors.New("couldn't parse claims")
		return 0, "", err
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("token expired")
		return 0, "", err
	}
	return claims.ID, claims.Email, nil
}

func ExtractToken(c *gin.Context) string {
	token := c.Query("token")
	if token != "" {
		return token
	}
	bearerToken := c.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

func ExtractTokenID(c *gin.Context) (int, string, error) {

	tokenString := ExtractToken(c)
	id, email, err := TokenValidation(tokenString)

	if err != nil {
		return 0, "", err
	}

	return id, email, nil
}
