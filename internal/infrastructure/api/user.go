package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"

	"dz1/internal/domain"
	"dz1/internal/infrastructure"
)

func (b *Builder) register(w http.ResponseWriter, r *http.Request) {
	var userReq domain.UserCreateReq

	if err := readRequestBody(r.Body, &userReq); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	if err := b.validator.Struct(userReq); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	userID, err := b.service.CreateUser(b.ctx, userReq.ToUser())
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err)
	}

	writeJSONResponse(w, http.StatusCreated, userID)
}

func readRequestBody[T any](r io.ReadCloser, body *T) error {
	defer r.Close()
	return json.NewDecoder(r).Decode(&body)
}

func (b *Builder) getUserByID(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")

	if err := b.validator.Var(userID, "uuid"); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	user, err := b.service.GetUser(b.ctx, userID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, infrastructure.ErrUserNotFound) {
			statusCode = http.StatusNotFound
		}

		writeErrorResponse(w, statusCode, err)
		return
	}

	writeJSONResponse(w, http.StatusOK, user)
}
