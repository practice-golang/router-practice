package handler

import (
	"log"
	"net/http"

	"router-practice/auth"
	"router-practice/internal/router"
)

func HelloGlobalMiddleware1(next router.Handler) router.Handler {
	return func(c *router.Context) {
		log.Println("HelloGlobalMiddleware1")

		c.Params = append(c.Params, "middleware global1")
		c.Text(http.StatusOK, "middleware global1")

		next(c)
	}
}

func HelloGlobalMiddleware2(next router.Handler) router.Handler {
	return func(c *router.Context) {
		log.Println("HelloGlobalMiddleware2")

		c.Params = append(c.Params, "middleware global2")
		c.Text(http.StatusOK, "middleware global2")

		next(c)
	}
}

func HelloMiddleware(next router.Handler) router.Handler {
	return func(c *router.Context) {
		c.Text(http.StatusInternalServerError, "Middle ware test error")

		// Parameter transfer test
		// c.Request = c.WithContext(context.WithValue(c.Context(), router.ContextKey("say"), "hello middleware"))

		// next(c)
	}
}

func AuthMiddleware(next router.Handler) router.Handler {
	return func(c *router.Context) {
		claim, err := auth.GetClaim(*c.Request, "cookie")
		if err != nil {
			auth.ExpireCookie(c.ResponseWriter)

			c.Text(http.StatusUnauthorized, "Auth error")

			return
		}

		c.AuthInfo = claim

		next(c)
	}
}

func AuthApiMiddleware(next router.Handler) router.Handler {
	return func(c *router.Context) {
		c.Request.Header.Get("Authorization")
		claim, err := auth.GetClaim(*c.Request, "header")
		if err != nil {
			log.Println("AuthApiMiddleware:", err)
			c.Text(http.StatusUnauthorized, "Auth error")

			return
		}

		c.AuthInfo = claim

		next(c)
	}
}
