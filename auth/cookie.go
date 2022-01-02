package auth

import (
	"log"
	"net/http"
)

func SetCookieHeader(w http.ResponseWriter, token string, duration int64) {
	session := http.Cookie{
		Name:     "token",
		Value:    token,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(duration),
		HttpOnly: true,
	}

	w.Header().Add("Set-Cookie", session.String())
}

func ExpireCookie(w http.ResponseWriter) {
	session := http.Cookie{
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
		HttpOnly: true,
	}

	w.Header().Add("Set-Cookie", session.String())
}

// GetCookie - Not use
func GetCookie(r *http.Request) {
	token, err := r.Cookie("token")
	if err != nil {
		log.Println("GetCookie:", err)
	}
	cookies := r.Cookies()

	log.Println(token)
	log.Println(cookies)
}
