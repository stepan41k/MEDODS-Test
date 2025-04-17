package jwt

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)


const (
	accessDuration = 15 * time.Minute
	refreshDuration =  43200 * time.Minute
	diffDuration = 2591100
)

var (
	ErrUserUnauthorized = errors.New("user unauthorized")
)


func NewAccesssTokens(guid string, key string, ip string) (string, error) {
	accessToken := jwt.New(jwt.SigningMethodHS512)

	claims := accessToken.Claims.(jwt.MapClaims)
	claims["guid"] = guid
	claims["ip"] = ip
	claims["key"] = key
	claims["exp"] = time.Now().Add(accessDuration).Unix()

	accessTokenString, err := accessToken.SignedString([]byte(os.Getenv("SECRET_FOR_ACCESS")))
	if err != nil {
		return "", err
	}

	return accessTokenString, nil
}


func NewRefreshToken(key string, ip string) (string, error) {
	refreshToken := jwt.New(jwt.SigningMethodHS512)

	claims := refreshToken.Claims.(jwt.MapClaims)
	claims["ip"] = ip
	claims["key"] = key
	claims["exp"] = time.Now().Add(refreshDuration).Unix()

	refreshTokenString, err := refreshToken.SignedString([]byte(os.Getenv("SECRET_FOR_REFRESH")))
	if err != nil {
		return "", err
	}

	return refreshTokenString, nil
}


func CheckMatching(accessToken string, refreshToken string) (bool, error) {
	accesstoken, _ := jwt.Parse(accessToken, func(token *jwt.Token) (any, error) {
		return []byte(os.Getenv("SECRET_FOR_ACCESS")), nil
	})

	refreshtoken, _ := jwt.Parse(refreshToken, func(token *jwt.Token) (any, error) {
		return []byte(os.Getenv("SECRET_FOR_REFRESH")), nil
	})

	accessClaims := accesstoken.Claims.(jwt.MapClaims)
	refreshClaims := refreshtoken.Claims.(jwt.MapClaims)

	if accessClaims["key"].(string) != refreshClaims["key"].(string) {
		return false, fmt.Errorf("tokens dont match")
	}

	return true, nil
}


func CheckRefresh(refreshToken string) (bool, error) {
	refreshtoken, err := jwt.Parse(refreshToken, func(token *jwt.Token) (any, error) {
		return []byte(os.Getenv("SECRET_FOR_REFRESH")), nil
	})

	if err != nil {
		return false, err
	}

	if !refreshtoken.Valid {
		return false, ErrUserUnauthorized
	}

	return true, nil
}


func GetGUID(accessToken string) []byte {
	accesstoken, _ := jwt.Parse(string(accessToken), func(token *jwt.Token) (any, error) {
		return []byte(os.Getenv("SECRET_FOR_ACCESS")), nil
	})

	claims := accesstoken.Claims.(jwt.MapClaims)
	guid := claims["guid"].(string)

	return []byte(guid)
}


func CheckIP(newIP string, refreshToken string) error {
	refreshtoken, _ := jwt.Parse(refreshToken, func(token *jwt.Token) (any, error) {
		return []byte(os.Getenv("SECRET_FOR_REFRESH")), nil
	})

	accessClaims := refreshtoken.Claims.(jwt.MapClaims)
	accessIp := accessClaims["ip"].(string)

	if newIP != accessIp {
		fmt.Println("ip changed")
		return fmt.Errorf("ip changed")
		//
	}

	return nil
}