package handlers

import (
	"encoding/json"
	"net/http"

	"sentinelos/core/internal/auth"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token     string `json:"token"`
	ExpiresIn int64  `json:"expires_in"`
}

func LoginHandler(authService *auth.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}

		user, err := authService.Authenticate(req.Username, req.Password)
		if err != nil {
			http.Error(w, "invalid credentials"+err.Error(), http.StatusUnauthorized)
			return
		}

		token, expires, err := auth.GenerateJWT(user)
		if err != nil {
			http.Error(w, "token error", http.StatusInternalServerError)
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
