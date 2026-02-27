package handlers

import (
	"net/http"

	"sentinelos/core/internal/system/config_engine"
)

func BeginConfigHandler(w http.ResponseWriter, r *http.Request) {

	//username := r.Context().Value("username").(string)
    username := "admin" // TODO: obtener del contexto de autenticación
	err := config_engine.BeginConfig(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	

	w.Write([]byte("config session started"))
}