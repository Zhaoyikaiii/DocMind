package middleware

import (
	"github.com/gin-gonic/gin"
)

type Middleware struct {
	Auth     gin.HandlerFunc
	CORS     gin.HandlerFunc
	Logger   gin.HandlerFunc
	Recovery gin.HandlerFunc
}

func NewMiddleware() *Middleware {
	return &Middleware{
		Auth:     AuthMiddleware(),
		CORS:     CORSMiddleware(),
		Logger:   RequestLogger(),
		Recovery: Recovery(),
	}
}

// ApplyMiddleware applies all middleware to the router
func ApplyMiddleware(r *gin.Engine) {
	m := NewMiddleware()

	r.Use(m.Recovery)
	r.Use(m.Logger)
	r.Use(m.CORS)

}
