package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/atdevten/peace/internal/application/usecases"
	"github.com/atdevten/peace/internal/domain/repositories"
	jwtinfra "github.com/atdevten/peace/internal/infrastructure/auth/jwt"
	"github.com/atdevten/peace/internal/infrastructure/config"
	redisclient "github.com/atdevten/peace/internal/infrastructure/database/redis"
	"github.com/atdevten/peace/internal/infrastructure/database/redis/repository"
	httpmiddleware "github.com/atdevten/peace/internal/interfaces/http/middleware"
	websocketHandlers "github.com/atdevten/peace/internal/interfaces/websocket/handlers"
	"github.com/gin-gonic/gin"
)

// WebSocketServer wires infrastructure, application and interface layers, and runs WebSocket server
type WebSocketServer struct {
	cfg                  *config.Config
	userOnlineStatusRepo repositories.UserOnlineStatusRepository
	userOnlineStatusUC   usecases.UserOnlineStatusUseCase
	engine               *gin.Engine
	httpServer           *http.Server
}

// NewWebSocketServer initializes dependencies, registers routes, and returns a ready server instance
func NewWebSocketServer(cfg *config.Config) (*WebSocketServer, error) {
	// Initialize Redis client (real)
	addr := cfg.GetRedisAddr()
	password := cfg.Database.Redis.Password
	db := cfg.Database.Redis.DB
	redisCli := redisclient.NewRealClient(addr, password, db)
	if err := redisCli.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis at %s: %w", addr, err)
	}
	userOnlineStatusRepo := repository.NewRedisUserOnlineStatusRepository(redisCli)

	// JWT service and middleware
	jwtSvc := jwtinfra.NewService(cfg.Auth.JWT.Secret, cfg.Auth.JWT.Expiration, cfg.Auth.JWT.RefreshExpiration)
	authMW := httpmiddleware.NewAuthMiddleware(jwtSvc)

	// Use cases
	userOnlineStatusUC := usecases.NewUserOnlineStatusUseCase(userOnlineStatusRepo)

	// Handlers
	onlineStatusHandler := websocketHandlers.NewOnlineStatusHandler(userOnlineStatusUC, jwtSvc)

	// Gin engine
	engine := gin.Default()

	// Health endpoints
	engine.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// WebSocket routes under /ws
	wsGroup := engine.Group("/ws")
	{
		// Optional auth to allow Sec-WebSocket-Protocol fallback
		wsGroup.Use(authMW.OptionalAuth())

		// Single WebSocket endpoint for all online status operations
		wsGroup.GET("", onlineStatusHandler.HandleWebSocket)

		// Health check for WebSocket server
		wsGroup.GET("/health", func(c *gin.Context) {
			onlineUsers, err := userOnlineStatusUC.GetOnlineUsers(c.Request.Context())
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  "error",
					"message": "Failed to get online users",
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"status":            "healthy",
				"online_users":      len(onlineUsers),
				"total_connections": len(onlineUsers),
			})
		})
	}

	s := &WebSocketServer{
		cfg:                  cfg,
		userOnlineStatusRepo: userOnlineStatusRepo,
		userOnlineStatusUC:   userOnlineStatusUC,
		engine:               engine,
	}

	// Prepare http.Server with timeouts
	wsPort := "8081" // Default WebSocket port
	if cfg.Server.WebSocketPort != "" {
		wsPort = cfg.Server.WebSocketPort
	}

	s.httpServer = &http.Server{
		Addr:              fmt.Sprintf(":%s", wsPort),
		Handler:           engine,
		ReadTimeout:       cfg.Server.ReadTimeout,
		ReadHeaderTimeout: cfg.Server.ReadTimeout,
		WriteTimeout:      cfg.Server.WriteTimeout,
		IdleTimeout:       cfg.Server.IdleTimeout,
	}

	return s, nil
}

// Run starts the WebSocket server and blocks until it stops
func (s *WebSocketServer) Run() error {
	if s.httpServer == nil {
		return fmt.Errorf("websocket server is not initialized")
	}

	fmt.Printf("WebSocket server starting on %s\n", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the WebSocket server and closes connections
func (s *WebSocketServer) Shutdown(ctx context.Context) error {
	var firstErr error

	// Shutdown HTTP server
	if s.httpServer != nil {
		if err := s.httpServer.Shutdown(ctx); err != nil {
			firstErr = fmt.Errorf("websocket shutdown: %w", err)
		}
	}

	return firstErr
}

// Engine exposes the underlying Gin engine (useful for tests)
func (s *WebSocketServer) Engine() *gin.Engine {
	return s.engine
}
