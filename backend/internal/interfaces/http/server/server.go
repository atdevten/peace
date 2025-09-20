package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/atdevten/peace/internal/application/services/google"
	appjwt "github.com/atdevten/peace/internal/application/services/jwt"
	appUsecases "github.com/atdevten/peace/internal/application/usecases"
	infraJWT "github.com/atdevten/peace/internal/infrastructure/auth/jwt"
	infraConfig "github.com/atdevten/peace/internal/infrastructure/config"
	infraDB "github.com/atdevten/peace/internal/infrastructure/database"
	pgRepo "github.com/atdevten/peace/internal/infrastructure/database/postgres/repository"
	infraLogging "github.com/atdevten/peace/internal/infrastructure/logging"
	infraTracing "github.com/atdevten/peace/internal/infrastructure/tracing"
	httpHandlers "github.com/atdevten/peace/internal/interfaces/http/handlers"
	httpMiddleware "github.com/atdevten/peace/internal/interfaces/http/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// HTTPServer wires infrastructure, application and interface layers, and runs Gin HTTP server
type HTTPServer struct {
	cfg        *infraConfig.Config
	dbManager  *infraDB.DatabaseManager
	engine     *gin.Engine
	httpServer *http.Server
	logger     *infraLogging.Logger
	tracer     *infraTracing.Tracer
}

// NewHTTPServer initializes dependencies, registers routes, and returns a ready server instance
func NewHTTPServer(cfg *infraConfig.Config) (*HTTPServer, error) {
	// Initialize logging
	logger := infraLogging.NewLogger(&infraLogging.Config{
		Level:      cfg.Log.Level,
		Format:     cfg.Log.Format,
		TimeFormat: cfg.Log.TimeFormat,
		Caller:     cfg.Log.Caller,
		CallerSkip: cfg.Log.CallerSkip,
		FilePath:   cfg.Log.FilePath,
	})

	// Initialize tracing
	tracer, err := infraTracing.NewTracer(&infraTracing.Config{
		Enabled:          cfg.ELK.Enabled,
		ServiceName:      cfg.ELK.ServiceName,
		Environment:      cfg.ELK.Environment,
		APMServerURL:     cfg.ELK.APMServerURL,
		ElasticsearchURL: cfg.ELK.ElasticsearchURL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize tracing: %w", err)
	}

	// Initialize database(s)
	dbManager, err := infraDB.NewDatabaseManager(cfg)
	if err != nil {
		return nil, fmt.Errorf("infraDB.NewDatabaseManager: %w", err)
	}

	// Repositories
	userRepo := pgRepo.NewPostgreSQLUserRepository(dbManager.Postgres)
	recordRepo := pgRepo.NewPostgreSQLMentalHealthRecordRepository(dbManager.Postgres)
	quoteRepo := pgRepo.NewPostgreSQLQuoteRepository(dbManager.Postgres)
	tagRepo := pgRepo.NewTagRepository(dbManager.Postgres)

	// Services (infrastructure implementation for application port)
	var jwtService appjwt.Service = infraJWT.NewService(
		cfg.Auth.JWT.Secret,
		cfg.Auth.JWT.Expiration,
		cfg.Auth.JWT.RefreshExpiration,
	)

	// Google OAuth service
	googleService := google.NewService(
		cfg.Auth.Google.ClientID,
		cfg.Auth.Google.ClientSecret,
		cfg.Auth.Google.RedirectURI,
	)

	// Use cases
	authUC := appUsecases.NewAuthUseCase(userRepo, jwtService, googleService)
	userUC := appUsecases.NewUserUseCase(userRepo)
	recordUC := appUsecases.NewMentalHealthRecordUseCase(recordRepo)
	quoteUC := appUsecases.NewQuoteUseCase(quoteRepo)
	tagUC := appUsecases.NewTagUseCase(tagRepo, quoteRepo)

	// Handlers
	authHandler := httpHandlers.NewAuthHandler(authUC)
	userHandler := httpHandlers.NewUserHandler(userUC)
	recordHandler := httpHandlers.NewMentalHealthRecordHandler(recordUC)
	quoteHandler := httpHandlers.NewQuoteHandler(quoteUC)
	tagHandler := httpHandlers.NewTagHandler(tagUC)

	// Middleware
	authMW := httpMiddleware.NewAuthMiddleware(jwtService)
	loggingMW := httpMiddleware.NewLoggingMiddleware(logger)

	// Gin engine
	engine := gin.Default()

	// Add OpenTelemetry middleware
	if tracer.IsEnabled() {
		engine.Use(otelgin.Middleware(cfg.ELK.ServiceName))
	}

	// Add logging middleware
	engine.Use(loggingMW.RequestLogger())

	// CORS middleware
	corsConfig := cors.Config{
		AllowOrigins:     cfg.App.CORS.AllowedOrigins,
		AllowMethods:     cfg.App.CORS.AllowedMethods,
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-Requested-With"},
		AllowCredentials: true,
	}
	if len(corsConfig.AllowOrigins) == 0 {
		// Fallback: allow all origins in dev if not configured
		corsConfig = cors.DefaultConfig()
		corsConfig.AllowAllOrigins = true
		corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
		corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-Requested-With"}
		corsConfig.AllowCredentials = true
	}
	engine.Use(cors.New(corsConfig))

	// Public health endpoints
	engine.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// API routes under /api
	api := engine.Group("/api")

	// Auth routes (public)
	authGroup := api.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/refresh", authHandler.Refresh)
		authGroup.GET("/google/url", authHandler.GetGoogleAuthURL)
		authGroup.POST("/google/login", authHandler.GoogleLogin)
	}

	// User routes (protected)
	userGroup := api.Group("/user")
	userGroup.Use(authMW.RequireAuth())
	{
		userGroup.GET("/me", userHandler.Me)
		userGroup.PUT("/profile", userHandler.UpdateProfile)
		userGroup.PUT("/password", userHandler.UpdatePassword)
		userGroup.POST("/deactivate", userHandler.Deactivate)
		userGroup.DELETE("/account", userHandler.DeleteAccount)
	}

	// Mental health records (protected)
	recordGroup := api.Group("/records")
	recordGroup.Use(authMW.RequireAuth())
	{
		recordGroup.POST("", recordHandler.Create)
		recordGroup.GET("", recordHandler.GetByCondition)
		recordGroup.GET("/heatmap", recordHandler.GetHeatmap)
		recordGroup.GET("/streak", recordHandler.GetStreak)
		recordGroup.GET("/:id", recordHandler.GetByID)
		recordGroup.PUT("/:id", recordHandler.Update)
		recordGroup.DELETE("/:id", recordHandler.Delete)
	}

	// Quotes (public)
	quotesGroup := api.Group("/quotes")
	{
		quotesGroup.GET("", quoteHandler.GetAllQuotes)
		quotesGroup.GET("/random", quoteHandler.GetRandomQuote)
		quotesGroup.GET("/:id", quoteHandler.GetByID)
		quotesGroup.POST("", quoteHandler.CreateQuote)
		quotesGroup.PUT("/:id", quoteHandler.UpdateQuote)
		quotesGroup.DELETE("/:id", quoteHandler.DeleteQuote)

		// Quote tags
		quotesGroup.GET("/:id/tags", tagHandler.GetTagsByQuoteID)
		quotesGroup.POST("/:id/tags", tagHandler.AddTagToQuote)
		quotesGroup.DELETE("/:id/tags", tagHandler.RemoveTagFromQuote)
	}

	// Tags (public)
	tagsGroup := api.Group("/tags")
	{
		tagsGroup.GET("", tagHandler.GetAllTags)
		tagsGroup.POST("", tagHandler.CreateTag)
		tagsGroup.PUT("/:id", tagHandler.UpdateTag)
		tagsGroup.DELETE("/:id", tagHandler.DeleteTag)
	}

	s := &HTTPServer{
		cfg:       cfg,
		dbManager: dbManager,
		engine:    engine,
		logger:    logger,
		tracer:    tracer,
	}

	// Prepare http.Server with timeouts
	s.httpServer = &http.Server{
		Addr:              cfg.GetServerAddr(),
		Handler:           engine,
		ReadTimeout:       cfg.Server.ReadTimeout,
		ReadHeaderTimeout: cfg.Server.ReadTimeout,
		WriteTimeout:      cfg.Server.WriteTimeout,
		IdleTimeout:       cfg.Server.IdleTimeout,
	}

	return s, nil
}

// Run starts the HTTP server and blocks until it stops
func (s *HTTPServer) Run() error {
	if s.httpServer == nil {
		return fmt.Errorf("http server is not initialized")
	}
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the HTTP server and closes database connections
func (s *HTTPServer) Shutdown(ctx context.Context) error {
	var firstErr error
	if s.httpServer != nil {
		if err := s.httpServer.Shutdown(ctx); err != nil {
			firstErr = fmt.Errorf("http shutdown: %w", err)
		}
	}
	// Close DB connections
	if s.dbManager != nil {
		s.dbManager.Close()
	}
	// Shutdown tracing
	if s.tracer != nil {
		if err := s.tracer.Shutdown(ctx); err != nil {
			if firstErr == nil {
				firstErr = fmt.Errorf("tracer shutdown: %w", err)
			}
		}
	}
	return firstErr
}

// Engine exposes the underlying Gin engine (useful for tests)
func (s *HTTPServer) Engine() *gin.Engine { return s.engine }
