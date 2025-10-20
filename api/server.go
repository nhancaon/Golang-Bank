package api

import (
	"fmt"
	db "simple-bank/db/sqlc"
	"simple-bank/token"
	"simple-bank/util"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	config     util.Config
	router     *gin.Engine
	tokenMaker token.Maker
	store      db.Store
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

// NewServer creates a new HTTP server and setup routing
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}
	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	r := gin.Default()
	r.POST("/accounts", server.createAccount)
	r.GET("/accounts/:id", server.getAccount)
	r.GET("/accounts", server.listAccounts)

	r.POST("/transfers", server.createTransfer)

	r.POST("/users", server.createUser)
	r.POST("/users/login", server.loginUser)

	server.router = r
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
