package handlers

import (
	"encoding/json"
	"net/http"

	"sentinelos/core/internal/api/middleware"
)

func MeHandler(w http.ResponseWriter, r *http.Request) {

//	claims, ok := r.Context().Value("user").(*auth.Claims)
	claims, ok := middleware.UserFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	resp := map[string]interface{}{
		"username": claims.Username,
		"role":     claims.Role,
		"expires":  claims.ExpiresAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
