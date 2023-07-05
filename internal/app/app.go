package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"dz1/internal/domain"
	"dz1/internal/infrastructure/api"
	"dz1/internal/infrastructure/config"
	"dz1/internal/infrastructure/repository"
)

type App struct {
	cfg *config.Config
}

func New(cfg *config.Config) *App {
	return &App{cfg: cfg}
}

func (a App) Run(ctx context.Context, logger *zap.Logger) error {
	databaseUrl := fmt.Sprintf("postgres://%s:%s@%s/%s",
		a.cfg.Database.Username,
		a.cfg.Database.Password,
		a.cfg.Database.Addr,
		a.cfg.Database.Database,
	)
	poolConn, err := pgxpool.New(context.Background(), databaseUrl)
	if err != nil {
		return fmt.Errorf("Unable to connect to database: %v\n", err)
	}
	defer poolConn.Close()
	logger.Info("database successfully connected")

	userRepository := repository.NewUser(poolConn)
	userService := domain.NewUserService(userRepository)
	tokenAuth := jwtauth.New("HS256", []byte(a.cfg.App.Auth), nil)

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{ /*"https://*", "http://*", */ "*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods: []string{"GET", "POST", "PATCH", "DELETE"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token",
			"Origin", "DNT", "X-CustomHeader", "Keep-Alive", "User-Agent",
			"X-Requested-With", "If-Modified-Since", "Cache-Control", "Content-Range", "Range"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers

	}))
	router.Use(middleware.Logger)
	apiBuilder := api.NewBuilder()
	apiBuilder.
		Service(userService).
		Logger(logger).
		Context(ctx).
		Auth(tokenAuth).
		RegisterRoutes(router)

	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", a.cfg.App.Listen),
		Handler: router,
	}

	go func(shutdown func(ctx context.Context) error) {
		<-ctx.Done()
		logger.Info("context canceled")
		_ = shutdown(ctx)
	}(srv.Shutdown)

	logger.Info("http server",
		zap.Int("listen", a.cfg.App.Listen),
	)

	return srv.ListenAndServe()

}
