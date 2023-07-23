package api

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"dz1/internal/domain"
	"dz1/internal/infrastructure/api/validators"
)

type Builder struct {
	service   domain.UserService
	validator *validator.Validate
	logger    *zap.Logger
	ctx       context.Context
	tokenAuth *jwtauth.JWTAuth
}

func NewBuilder() *Builder {
	v := validator.New()
	validators.RegisterCustomValidators(v)

	return &Builder{
		validator: v,
	}
}

func (b *Builder) Service(service domain.UserService) *Builder {
	b.service = service
	return b
}

func (b *Builder) Logger(logger *zap.Logger) *Builder {
	b.logger = logger
	return b
}

func (b *Builder) Context(ctx context.Context) *Builder {
	b.ctx = ctx
	return b
}

func (b *Builder) Auth(tokenAuth *jwtauth.JWTAuth) *Builder {
	b.tokenAuth = tokenAuth
	return b
}

func (b *Builder) RegisterRoutes(router chi.Router) {
	router.Post("/login", b.login)
	router.Route("/user", func(r chi.Router) {
		r.Post("/register", b.register)

		// protected
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(b.tokenAuth))
			r.Use(jwtauth.Authenticator)

			r.Get("/get/{id}", b.getUserByID)
			r.Get("/search", b.getSearchUsers)
		})

	})
}
