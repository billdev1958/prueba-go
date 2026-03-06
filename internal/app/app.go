package app

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	v1 "prueba-go/internal/infrastructure/http/v1"
	"prueba-go/internal/infrastructure/http/v1/audit"
	"prueba-go/internal/infrastructure/http/v1/comercios"
	"prueba-go/internal/infrastructure/http/v1/reports"
	"prueba-go/internal/infrastructure/http/v1/transactions"
	"prueba-go/internal/infrastructure/postgres"
	auditRepo "prueba-go/internal/infrastructure/postgres/audit"
	comercioRepo "prueba-go/internal/infrastructure/postgres/comercio"
	reportRepo "prueba-go/internal/infrastructure/postgres/report"
	transactionRepo "prueba-go/internal/infrastructure/postgres/transaction"
	auditUC "prueba-go/internal/usecases/audit"
	comercioUC "prueba-go/internal/usecases/commerce"
	reportUC "prueba-go/internal/usecases/report"
	transactionUC "prueba-go/internal/usecases/transaction"
	"prueba-go/pkg/util/logger"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	cors "github.com/rs/cors/wrapper/gin"
)

type App struct {
	pool   *pgxpool.Pool
	port   string
	router *gin.Engine
}

func NewApp() (*App, error) {
	// 1. Get configuration from environment
	port := os.Getenv("PORT")

	dsn := os.Getenv("DATABASE_URL")

	// 2. Setup database
	pool, err := setupDatabase(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to setup database: %w", err)
	}

	// 3. Initialize router (Gin)
	router := gin.New()
	router.Use(gin.Recovery())

	return &App{
		pool:   pool,
		port:   port,
		router: router,
	}, nil
}

func (a *App) Run() error {
	// Logger configuration using slog
	l := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(l)

	// Custom Logger Middleware
	a.router.Use(logger.GinLogger(l))

	// CORS
	a.router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPatch, http.MethodOptions, http.MethodPut},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Authorization", "Accept", "X-User-Id"},
		AllowCredentials: true,
	}))

	// Initializing services/modules
	if err := a.initializeModules(); err != nil {
		return fmt.Errorf("failed to initialize modules: %w", err)
	}

	// Server config
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", a.port),
		Handler: a.router,
	}

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-quit
		log.Println("server is shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Server forced to shutdown: %v", err)
		}
		a.pool.Close()
	}()

	log.Printf("starting server on port %s...", a.port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to listen and serve: %w", err)
	}

	return nil
}

func (a *App) initializeModules() error {
	// Repositories
	baseRepo := postgres.NewPgxRepository(a.pool)
	commRepo := comercioRepo.NewComercioRepository(a.pool)
	transRepo := transactionRepo.NewTransactionRepository(a.pool)
	repRepo := reportRepo.NewReportRepository(a.pool)
	audRepo := auditRepo.NewPgxAuditRepository(baseRepo)

	// Usecases
	audUsecase := auditUC.NewUsecases(audRepo)
	commUsecase := comercioUC.NewUsecases(commRepo, audRepo)
	transUsecase := transactionUC.NewUsecases(transRepo, audRepo)
	repUsecase := reportUC.NewUsecases(repRepo, audRepo)

	// Handlers
	commHandler := comercios.NewComercioHandler(commUsecase)
	transHandler := transactions.NewTransactionHandler(transUsecase)
	repHandler := reports.NewReportHandler(repUsecase)
	audHandler := audit.NewAuditHandler(audUsecase)

	// Route registration
	v1.RegisterRoutes(a.router, v1.RouterConfig{
		CommerceHandler:    commHandler,
		ReportHandler:      repHandler,
		TransactionHandler: transHandler,
		AuditHandler:       audHandler,
	})

	return nil
}

func setupDatabase(dsn string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbPool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	return dbPool, nil
}
