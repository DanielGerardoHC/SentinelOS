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

			protected.Route("/interfaces", func(rt chi.Router) {
				rt.Get("/", handlers.InterfacesHandler)
				rt.Put("/{name}", handlers.EditInterfaceHandler)
			})

			protected.Route("/routes", func(rt chi.Router) {
				rt.Get("/", handlers.RoutesHandler)
				rt.Post("/", handlers.CreateRouteHandler)
				rt.Put("/{id}", handlers.EditRouteHandler)
				rt.Delete("/{id}", handlers.DeleteRouteHandler)
			})

			protected.Route("/vlans", func(rt chi.Router) {
				rt.Get("/", handlers.VlansHandler)
				rt.Post("/", handlers.CreateVlanHandler)
				rt.Put("/{name}", handlers.EditVlanHandler)
				rt.Delete("/{name}", handlers.DeleteVlanHandler)
			})

			protected.Route("/dhcp", func(rt chi.Router) {
				rt.Get("/", handlers.DhcpHandler)
				rt.Post("/", handlers.CreateDhcpHandler)
				rt.Put("/{interface}", handlers.EditDhcpHandler)
				rt.Delete("/{interface}", handlers.DeleteDhcpHandler)
			})

			protected.Route("/nat", func(rt chi.Router) {
				rt.Get("/", handlers.NatHandler)
				rt.Post("/", handlers.CreateNatHandler)
				rt.Post("/{id}/move", handlers.MoveNatHandler)
				rt.Put("/{id}", handlers.EditNatHandler)
				rt.Delete("/{id}", handlers.DeleteNatHandler)
			})

			protected.Route("/zones", func(rt chi.Router) {
				rt.Get("/", handlers.ZonesHandler)
				rt.Post("/", handlers.CreateZoneHandler)
				rt.Put("/{name}", handlers.EditZoneHandler)
				rt.Delete("/{name}", handlers.DeleteZoneHandler)
			})

			protected.Route("/services", func(rt chi.Router) {
				rt.Get("/", handlers.ServicesHandler)
				rt.Post("/", handlers.CreateServiceHandler)
				rt.Put("/{name}", handlers.EditServiceHandler)
				rt.Delete("/{name}", handlers.DeleteServiceHandler)
			})

			protected.Route("/addresses", func(rt chi.Router) {
				rt.Get("/", handlers.AddressesHandler)
				rt.Post("/", handlers.CreateAddressHandler)
				rt.Put("/{name}", handlers.EditAddressHandler)
				rt.Delete("/{name}", handlers.DeleteAddressHandler)
			})

			protected.Route("/policies", func(rt chi.Router) {
				rt.Get("/", handlers.PoliciesHandler)
				rt.Post("/", handlers.CreatePolicyHandler)
				rt.Post("/{id}/move", handlers.MovePolicyHandler)
				rt.Put("/{id}", handlers.EditPolicyHandler)
				rt.Delete("/{id}", handlers.DeletePolicyHandler)
			})

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
