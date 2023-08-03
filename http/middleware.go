package http

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Middleware struct {
	Cors gin.HandlerFunc
}

/*
Setup various middleware tools
*/
func (m *Middleware) setupMiddleware() {
	//setup various middleware
	m.setCors()
}

/*
Setup CORS middleware
*/
func (m *Middleware) setCors() {
	m.Cors = cors.Default()
}

/*
Use middleware in the router engine
*/
func UseMiddleware(r *gin.Engine) {
	var middleware *Middleware = &Middleware{}

	middleware.setupMiddleware()
	r.Use(middleware.Cors)
}
