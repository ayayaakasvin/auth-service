package handlers

import (
	"net/http"
	"strings"

	"github.com/ayayaakasvin/auth-service/internal/ctx"
	"github.com/ayayaakasvin/auth-service/internal/libs/bcrypt"
	"github.com/ayayaakasvin/auth-service/internal/libs/bindjson"
	"github.com/ayayaakasvin/auth-service/internal/libs/validinput"
	"github.com/ayayaakasvin/auth-service/internal/models/request"
	"github.com/ayayaakasvin/auth-service/internal/models/response"
	"github.com/ayayaakasvin/auth-service/internal/models/token"
	"github.com/ayayaakasvin/auth-service/internal/repository/postgresql"
	"github.com/google/uuid"
)

const (
	AuthorizationHeader = "Authorization"
)

// LogIn handler
// @Summary      Log in
// @Description  Authenticates a user and returns access and refresh tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        payload  body      request.UserRequest  true  "Login payload"
// @Success      200      {object}  response.JsonResponse
// @Failure      400      {object}  response.JsonResponse
// @Failure      401      {object}  response.JsonResponse
// @Failure      500      {object}  response.JsonResponse
// @Router       /api/auth/login [post]
func (h *Handlers) LogIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var loginReq request.UserRequest
		if err := bindjson.BindJson(r.Body, &loginReq); err != nil {
			response.SendErrorJson(w, http.StatusBadRequest, "failed to bind request")
			return
		}

		userId, err := h.repo.AuthentificateUser(r.Context(), loginReq.Username, loginReq.Password)
		if err != nil {
			switch err.Error() {
			case postgresql.NotFound:
				response.SendErrorJson(w, http.StatusUnauthorized, "invalid credentials")
			case postgresql.UnAuthorized:
				response.SendErrorJson(w, http.StatusUnauthorized, "invalid credentials")
			}
			return
		}

		sessionId := uuid.New().String()

		accessToken, err := h.jwtM.GenerateToken(token.NewAccessTokenClaims(userId, sessionId, h.jwtM.AccessTokenTTL))
		if err != nil {
			response.SendErrorJson(w, http.StatusInternalServerError, "access token generation error")
			h.logger.WithError(err).Error("failed to generate token")
			return
		}

		refreshToken, err := h.jwtM.GenerateToken(token.NewRefreshTokenClaims(userId, h.jwtM.RefreshTokenTTL))
		if err != nil {
			response.SendErrorJson(w, http.StatusInternalServerError, "refresh token generation error")
			h.logger.WithError(err).Error("failed to generate token")
			return
		}

		data := response.NewData()
		data["access-token"] = accessToken
		data["refresh-token"] = refreshToken
		h.logger.Info(data)

		if err := h.cache.Set(r.Context(), sessionId, true, h.jwtM.AccessTokenTTL); err != nil {
			h.logger.WithField("session_id", sessionId).WithError(err).Error("failed to set session id")
			response.SendErrorJson(w, http.StatusInternalServerError, "cache error")
			return
		}

		response.SendSuccessJson(w, http.StatusOK, data)
	}
}

// Register handler
// @Summary      Register a new user
// @Description  Creates a new user account
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        payload  body      request.UserRequest  true  "Register payload"
// @Success      201      {object}  response.JsonResponse
// @Failure      400      {object}  response.JsonResponse
// @Failure      500      {object}  response.JsonResponse
// @Router       /api/auth/register [post]
func (h *Handlers) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var registerReq request.UserRequest
		if err := bindjson.BindJson(r.Body, &registerReq); err != nil {
			response.SendErrorJson(w, http.StatusBadRequest, "failed to bind request")
			return
		}

		if err := validinput.IsValidUsername(registerReq.Username); err != nil {
			response.SendErrorJson(w, http.StatusBadRequest, "invalid username for register: %s", err.Error())
			return
		}
		if err := validinput.IsValidPassword(registerReq.Password); err != nil {
			response.SendErrorJson(w, http.StatusBadRequest, "invalid password for register %s", err.Error())
			return
		}

		hashed, err := bcrypt.BcryptHashing(registerReq.Password)
		if err != nil {
			h.logger.WithError(err).Error("bcrypt hashing failed")
			response.SendErrorJson(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		if err := h.repo.RegisterUser(r.Context(), registerReq.Username, hashed); err != nil {
			h.logger.WithError(err).Error("register user failed")
			response.SendErrorJson(w, http.StatusInternalServerError, "failed to register")
			return
		}

		response.SendSuccessJson(w, http.StatusCreated, nil)
	}
}

// LogOut handler
// @Summary      Log out
// @Description  Invalidates the current session (deletes session id)
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200      {object}  response.JsonResponse
// @Failure      401      {object}  response.JsonResponse
// @Failure      500      {object}  response.JsonResponse
// @Router       /api/auth/logout [post]
func (h *Handlers) LogOut() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session_id, ok := r.Context().Value(ctx.CtxSessionIDKey).(string)
		if !ok {
			response.SendErrorJson(w, http.StatusUnauthorized, "missing session id")
			return
		}

		if err := h.cache.Del(r.Context(), session_id); err != nil {
			h.logger.Errorf("failed to delete session id: %v", err)
			response.SendErrorJson(w, http.StatusInternalServerError, "cache error")
			return
		}

		response.SendSuccessJson(w, http.StatusOK, nil)
	}
}

// RefreshTheToken handler
// @Summary      Refresh access token
// @Description  Exchanges a refresh token (sent in Authorization header) for a new access token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Refresh token in format: Bearer <token>"
// @Success      200            {object}  response.JsonResponse
// @Failure      400            {object}  response.JsonResponse
// @Failure      401            {object}  response.JsonResponse
// @Failure      500            {object}  response.JsonResponse
// @Router       /api/auth/refresh [post]
func (h *Handlers) RefreshTheToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get(AuthorizationHeader)
		if authHeader == "" {
			response.SendErrorJson(w, http.StatusUnauthorized, "authorization header missing")
			return
		}

		refreshTokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if refreshTokenString == authHeader {
			response.SendErrorJson(w, http.StatusUnauthorized, "authorization header missing")
			return
		}

		cl, err := h.jwtM.Validate(refreshTokenString, &token.RefreshTokenClaims{})
		if err != nil {
			response.SendErrorJson(w, http.StatusUnauthorized, "failed to validate jwt")
			return
		}

		fullClaims, ok := cl.(*token.RefreshTokenClaims)
		if !ok {
			response.SendErrorJson(w, http.StatusUnauthorized, "invalid claims")
			return
		}

		if fullClaims.UserID == 0 {
			response.SendErrorJson(w, http.StatusUnauthorized, "user_id is missing in refresh token")
			return
		}

		sessionId := uuid.New().String()

		accessToken, err := h.jwtM.GenerateToken(token.NewAccessTokenClaims(fullClaims.UserID, sessionId, h.jwtM.AccessTokenTTL))
		if err != nil {
			response.SendErrorJson(w, http.StatusInternalServerError, "access token generation error")
			h.logger.WithError(err).Error("failed to generate token")
			return
		}

		data := response.NewData()
		data["access-token"] = accessToken
		h.logger.Info(data)

		if err := h.cache.Set(r.Context(), sessionId, true, h.jwtM.AccessTokenTTL); err != nil {
			h.logger.Errorf("failed to set session id: %v", err)
			response.SendErrorJson(w, http.StatusInternalServerError, "cache error")
			return
		}

		response.SendSuccessJson(w, http.StatusOK, data)
	}
}
