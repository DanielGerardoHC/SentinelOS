package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"sentinelos/core/internal/auth"
	"sentinelos/core/pkg/utils"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token     string `json:"token"`
	ExpiresIn int64  `json:"expires_in"`
}

// LoginHandler godoc
// @Summary Admin Login
// @Description Authenticate an admin and receive a JWT token.
// @Tags authentication
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} LoginResponse "Successful login"
// @Failure 400 {object} utils.APIError "ERR_NET_1004 Invalid JSON"
// @Failure 401 {object} utils.APIError "ERR_SEC_5005 Invalid credentials"
// @Failure 405 {object} utils.APIError "ERR_NET_1009 Method not allowed"
// @Failure 500 {object} utils.APIError "ERR_SEC_5006 Failed to generate token"
// @Router /api/login [post]
func LoginHandler(authService *auth.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			utils.SendError(w, http.StatusMethodNotAllowed, "ERR_NET_1009", "Method not allowed", "use POST")
			return
		}

		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.SendError(w, http.StatusBadRequest, "ERR_NET_1004", "Invalid JSON payload", err.Error())
			return
		}

		user, err := authService.Authenticate(req.Username, req.Password)
		if err != nil {
			var apiErr *utils.APIError
			if errors.As(err, &apiErr) {
				utils.SendError(w, http.StatusUnauthorized, apiErr.Code, apiErr.Message, apiErr.Details)
			} else {
				utils.SendError(w, http.StatusUnauthorized, "ERR_SEC_5005", "Invalid credentials", err.Error())
			}
			return
		}

		token, expires, err := auth.GenerateJWT(user)
		if err != nil {
			utils.SendError(w, http.StatusInternalServerError, "ERR_SEC_5006", "Failed to generate token", err.Error())
			return
		}

		resp := LoginResponse{
			Token:     token,
			ExpiresIn: expires,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
