package cookie

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

var (
	ErrCookieNotSet = errors.New("cookie not set")
)


func TakeCookie(w http.ResponseWriter, r *http.Request) (string, error) {
	accessCookie, err := r.Cookie("access_token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return "", ErrCookieNotSet
		}
		return "", fmt.Errorf("error with token: %w", err)
	}

	return accessCookie.Value, nil
}


func SetCookie(w http.ResponseWriter, accessToken string) {
	cookie := &http.Cookie{
		Name: "access_token",
		Value: accessToken,
		Path: "/",
		HttpOnly: true,
		Secure: false,
	}
	http.SetCookie(w, cookie)
}


func DeleteCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name: "access_token",
		Value: "",
		Path: "/",
		HttpOnly: true,
		Secure: false,
		Expires: time.Now(),
	}
	
	http.SetCookie(w, cookie)
}