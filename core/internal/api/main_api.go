package api

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"sentinelos/core/internal/api/handlers"
	"sentinelos/core/internal/api/middleware"
	"sentinelos/core/internal/auth"
)

func StartAPIServer() {

	users, err := auth.LoadUsers("/srv/sentinelos/core/internal/auth/users.yml")
	if err != nil {
		log.Fatalf("error loading users: %v", err)
	}

	authService := auth.NewAuthService(users)

	r := chi.NewRouter()

	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.StripSlashes)

	// ruta principal del API
	r.Route("/api", func(api chi.Router) {

		// public
		api.Post("/login", handlers.LoginHandler(authService))

		// protected
		api.Group(func(protected chi.Router) {

			protected.Use(middleware.JWTMiddleware)

			protected.Get("/me", handlers.MeHandler)
			protected.Get("/status", handlers.StatusHandler)

			// interfaces
			protected.Route("/interfaces", func(rt chi.Router) {
				rt.Get("/", handlers.InterfacesHandler)
				rt.Put("/{name}", handlers.EditInterfaceHandler)
			})

			// routes
			protected.Route("/routes", func(rt chi.Router) {
				rt.Get("/", handlers.RoutesHandler)
				rt.Post("/", handlers.CreateRouteHandler)
				rt.Put("/{id}", handlers.EditRouteHandler)
				rt.Delete("/{id}", handlers.DeleteRouteHandler)
			})

			// vlans

			protected.Route("/vlans", func(rt chi.Router) {
				rt.Get("/", handlers.VlansHandler)
				rt.Post("/", handlers.CreateVlanHandler)
				rt.Put("/{name}", handlers.EditVlanHandler)
			})

			protected.Get("/policies", handlers.PoliciesHandler)
			protected.Get("/zones", handlers.ZonesHandler)
			protected.Get("/dhcp", handlers.DhcpHandler)
			protected.Get("/nat", handlers.NatHandler)

			// config engine
			protected.Post("/config/begin", handlers.BeginConfigHandler)
			protected.Post("/config/commit", handlers.CommitHandler)
		})
	})

	log.Println("SentinelOS API listening on :8080")

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalf("server error: %v", err)
	}
}
