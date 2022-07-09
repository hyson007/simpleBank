package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/hyson007/simpleBank/db/sqlc"
	"github.com/hyson007/simpleBank/token"
	"github.com/hyson007/simpleBank/util"
)

type Server struct {
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
	config     util.Config
}

// func NewServer(store db.Store) *Server {
// 	server := &Server{store: &store}
// 	router := gin.Default()

// 	router.POST("/accounts", server.createAccount)

// 	return server
// }

func NewServer(config util.Config, store db.Store) (*Server, error) {
	// tokenMaker, err := token.NewJWTMaker(config.TokenSymKey)
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymKey)
	if err != nil {
		return nil, fmt.Errorf("can not create token maker: %w", err)
	}
	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
		config:     config,
	}

	router := gin.Default()
	server.router = router

	//create customized validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	//we assume createUser no need authorization and same for login user
	// but create account does require authorziation, user can only create
	// account that belongs to that specific user
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	// the / means all routes in this group
	authRoutes := router.Group("/").Use(authMiddleWare(server.tokenMaker))
	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts", server.listAccount)
	authRoutes.POST("/transfers", server.createTransfer)

	return server, nil
}

func (s *Server) Start(addr string) error {
	fmt.Println("starting server")
	return s.router.Run(addr)
}

//gin.H is just a short cut for map[string]interface{}
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
