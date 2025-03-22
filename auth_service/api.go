package authservice

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type APIServer struct {
	addr     string
	dbClient *mongo.Client
}

func NewAPIServer(addr string, dbClient *mongo.Client) *APIServer {
	return &APIServer{
		addr:     addr,
		dbClient: dbClient,
	}
}

func (s *APIServer) Start() error {
	// set up gin router
	router := gin.Default()
	vi := router.Group("/api/v1")
	mongoDb := s.dbClient.Database("authservice")
	userStore := NewUserStore(mongoDb)

	// set up gin router
	return nil
}
