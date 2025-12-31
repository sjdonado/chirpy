package main

import (
	"chirpy/handler"
	"chirpy/internal/api"
	"chirpy/internal/database"
	"chirpy/middleware"
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	_ "github.com/lib/pq"
)

func main() {
	filepathRoot := http.Dir(".")
	port := "8080"

	apiConfig := api.NewConfig()

	db, err := sql.Open("postgres", apiConfig.DBURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	queries := database.New(db)

	metric_middleware := middleware.NewMetricMiddleware()
	auth_middleware := middleware.NewAuthMiddleware(apiConfig, queries)

	metrics_handler := handler.NewMetricsHandler(apiConfig, queries)
	auth_handler := handler.NewAuthHandler(apiConfig, queries)
	users_handler := handler.NewUsersHandler(apiConfig, queries)
	chirps_handler := handler.NewChirpsHandler(queries)

	r := chi.NewRouter()

	r.Handle("/app/", metric_middleware.FileServerHits(http.StripPrefix("/app", http.FileServer(filepathRoot))))

	r.Route("/admin", func(r chi.Router) {
		r.Use(metric_middleware.FileServerHits)
		r.Get("/metrics", metrics_handler.GetMetrics())
		r.Post("/reset", metrics_handler.ResetMetrics())
	})

	r.Get("/api/healthz", handler.GetHealthz)

	r.Post("/api/login", auth_handler.Login())
	r.Post("/api/refresh", auth_handler.RefreshToken())
	r.Post("/api/revoke", auth_handler.RevokeToken())

	r.Route("/api/users", func(r chi.Router) {
		r.Post("/", users_handler.CreateUser())
		r.With(auth_middleware.Authenticated).Put("/", users_handler.UpdateUser())
	})

	r.Route("/api/chirps", func(r chi.Router) {
		r.Get("/", chirps_handler.FilterChirps())
		r.Get("/{id}", chirps_handler.GetChirp())
		r.With(auth_middleware.Authenticated).Post("/", chirps_handler.CreateChirp())
		r.With(auth_middleware.Authenticated).Delete("/{id}", chirps_handler.DeleteChrip())
	})

	r.Post("/api/polka/webhooks", users_handler.UpgradeUser())

	s := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(s.ListenAndServe())
}
