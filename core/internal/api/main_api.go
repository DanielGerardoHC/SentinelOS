package api

import (
	"log"
	"net/http"
	"sentinelos/core/internal/api/handlers"
	"sentinelos/core/internal/auth"
)

func StartAPIServer() {

	// cargar usuarios
	users, err := auth.LoadUsers("/srv/sentinelos/core/internal/auth/users.yml")
	if err != nil {
		log.Fatalf("error loading users: %v", err)
	}

	// crear authservice
	authService := auth.NewAuthService(users)

	// router
	mux := http.NewServeMux()

	// ruta de login
	mux.HandleFunc("/api/login", handlers.LoginHandler(authService))

	// levantar servidor
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalf("server error: %v", err)
	}
}
