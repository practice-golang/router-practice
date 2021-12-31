package auth

import (
	"log"
	"net/http"
)

func SetCookieHeader(w http.ResponseWriter, r *http.Request) {
	c1 := http.Cookie{
		Name:     "name",
		Value:    "user",
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
	}
	c2 := http.Cookie{
		Name:     "session",
		Value:    "hello world",
		SameSite: http.SameSiteStrictMode,
		MaxAge:   60 * 60 * 24 * 1,
		HttpOnly: true,
	}
	w.Header().Set("Set-Cookie", c1.String())
	w.Header().Add("Set-Cookie", c2.String())
}

func GetCookie(r *http.Request) {
	session, err := r.Cookie("session")
	if err != nil {
		log.Println("GetCookie:", err)
	}
	cookies := r.Cookies()

	log.Println(session)
	log.Println(cookies)
}
