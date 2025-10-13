package api

import (
	db "simple-bank/db/sqlc"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
	store  *db.Store
}

func ErrorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

// NewServer creates a new HTTP server and setup routing
func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	r := gin.Default()

	r.POST("/accounts", server.createAccount)
	r.GET("/accounts/:id", server.getAccount)
	r.GET("/accounts", server.listAccounts)
	server.router = r

	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
