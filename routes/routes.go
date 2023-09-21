package routes

import (
	"blockchain/handlers"
	"blockchain/internal/types"

	"github.com/gin-gonic/gin"
)

type Router struct {
	Server       *gin.Engine
	AuthHandler  handlers.AuthHandlerInterface
	BlockHandler handlers.BlocksHandlerInterface
}

func NewRouter(server *gin.Engine,
	authHandler handlers.AuthHandlerInterface,
	blockHandler handlers.BlocksHandlerInterface,
) *Router {
	return &Router{
		server,
		authHandler,
		blockHandler,
	}
}

func (r *Router) RegisterRoutes(controllers ...types.Controller) {
	root := r.Server.Group("/")
	for _, ctr := range controllers {
		ctr.RegisterRoutes(root)
	}
}

func (r *Router) SetupRouter() {
	r.Server.Use(r.CORSMiddleware())
	r.RegisterRoutes()
}

func (r *Router) CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, id-token")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
