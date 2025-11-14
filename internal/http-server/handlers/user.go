package handlers

import (
	"net/http"

	"github.com/ayayaakasvin/auth-service/internal/ctx"
	"github.com/ayayaakasvin/auth-service/internal/models/response"
)

const userInfoKey = "user_id"

// PublicUserInfo returns public info about a user by query param
// @Summary      Get Public User Info
// @Description  Returns a user info by ID query parameter (public data only)
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user_id   query     int  true  "User ID"
// @Success      200  {object}  models.User
// @Failure      400  {object}  response.JsonResponse
// @Failure      404  {object}  response.JsonResponse
// @Router       /api/public/user [get]
func (h *Handlers) PublicUserInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDString := r.URL.Query().Get(userInfoKey)

		userID, err := h.jwtM.FetchUserID(userIDString)
		if err != nil {
			response.SendErrorJson(w, http.StatusBadRequest, "invalid user_id param")
			return
		}

		user, err := h.repo.GetPublicUserInfo(r.Context(), userID)
		if err != nil {
			response.SendErrorJson(w, http.StatusNotFound, "user not found")
			return
		}

		data := response.NewData()
		data["user"] = user

		response.SendSuccessJson(w, http.StatusOK, data)
	}
}

// PrivateUserInfo returns private info of the authenticated user
// @Summary      Get Private User Info
// @Description  Returns a user info (including sensitive fields) by JWT
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.User
// @Failure      400  {object}  response.JsonResponse
// @Failure      500  {object}  response.JsonResponse
// @Security     ApiKeyAuth
// @Router       /api/me [get]
func (h *Handlers) PrivateUserInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDUint := r.Context().Value(ctx.CtxUserIDKey).(uint)

		user, err := h.repo.GetPrivateUserInfo(r.Context(), userIDUint)
		if err != nil {
			response.SendErrorJson(w, http.StatusInternalServerError, "failed to get record")
			h.logger.Errorf("Failed to find record in db: %s", err.Error())
			return
		}

		data := response.NewData()
		data["user"] = user

		response.SendSuccessJson(w, http.StatusOK, data)
	}
}
