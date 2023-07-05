package api

import (
	"errors"
	"net/http"

	"dz1/internal/domain"
)

var errInvalidCreds = errors.New("invalid username or password")

func (b *Builder) login(w http.ResponseWriter, r *http.Request) {
	var req domain.UserAuthReq

	if err := readRequestBody(r.Body, &req); err != nil {
		writeErrorResponse(w, http.StatusUnauthorized, err)
	}

	if err := b.validator.Struct(req); err != nil {
		writeErrorResponse(w, http.StatusUnauthorized, errInvalidCreds)
		return
	}

	_, err := b.service.AuthenticateUser(b.ctx, req.ID, req.Password)
	if err != nil {
		writeErrorResponse(w, http.StatusUnauthorized, errInvalidCreds)
		return
	}

	_, tokenString, _ := b.tokenAuth.Encode(map[string]interface{}{"user_id": req.ID})

	writeJSONResponse(w, http.StatusOK, &domain.UserAuthResp{AccessToken: tokenString})
}
