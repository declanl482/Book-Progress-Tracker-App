package api

import (
	"github.com/declanl482/go-book-tracker-app/storage"
	"github.com/gin-gonic/gin"
)

type Server struct {
	ListenAddress string
	Storer        storage.Storage
	router        *gin.Engine
}

func NewServer(listenAddress string, storer storage.Storage) *Server {
	router := gin.Default()
	return &Server{
		ListenAddress: listenAddress,
		Storer:        storer,
		router:        router,
	}
}

func (s *Server) Start() error {
	// Create a new Gin router.

	s.router.Use(s.CORSMiddleware())
	s.router.Use(s.DBConnectionMiddleware())

	s.RegisterAuthHandlers()
	s.router.Use(s.RequireValidAccessToken())
	s.RegisterUserHandlers()
	s.RegisterBookHandlers()
	return s.router.Run(s.ListenAddress)
}

func (s *Server) RegisterAuthHandlers() {
	// Register the auth handlers.
	s.router.POST("/auth/login", s.handleLoginUser)
	s.router.POST("/auth/register", s.handleCreateUser)
}

func (s *Server) RegisterUserHandlers() {
	// Register the user handlers.
	s.router.GET("/users/:id", s.handleGetUser)
	s.router.PATCH("/users/:id", s.handleUpdateUser)
	s.router.DELETE("/users/:id", s.handleDeleteUser)
}

func (s *Server) RegisterBookHandlers() {
	// Register the book handlers.
	s.router.POST("/books/", s.handleCreateBook)
	s.router.GET("/books/", s.handleGetBooks)
	s.router.GET("/books/:id", s.handleGetBook)
	s.router.PATCH("/books/:id", s.handleUpdateBook)
	s.router.DELETE("/books/:id", s.handleDeleteBook)
}
