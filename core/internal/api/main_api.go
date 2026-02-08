package api

import (
	"log"
	"net/http"

	"sentinelos/core/internal/api/handlers"
	"sentinelos/core/internal/api/middleware"
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

	// ruta login
	mux.HandleFunc("/api/login", handlers.LoginHandler(authService))

	// ruta protegida /api/me
	mux.Handle(
		"/api/me",
		middleware.JWTMiddleware(
			http.HandlerFunc(handlers.MeHandler),
		),
	)

	mux.Handle(
	"/api/status",
	middleware.JWTMiddleware(
		http.HandlerFunc(handlers.StatusHandler),
	),
    )
    
	mux.Handle(
	"/api/interfaces",
	middleware.JWTMiddleware(
		http.HandlerFunc(handlers.InterfacesHandler),
	),
    )

	log.Println("SentinelOS API listening on :8080")

	// levantar servidor
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalf("server error: %v", err)
	}
}
