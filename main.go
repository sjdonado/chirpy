package main

import (
	"chirpy/handler"
	"chirpy/internal/database"
	"chirpy/middleware"
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()

	filepathRoot := http.Dir(".")
	port := "8080"

	db, err := sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	queries := database.New(db)

	metric_middleware := middleware.NewMetricMiddleware()
	auth_middleware := middleware.NewAuthMiddleware(queries)

	metrics_handler := handler.NewMetricsHandler(queries)
	auth_handler := handler.NewAuthHandler(queries)
	users_handler := handler.NewUsersHandler(queries)
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
		r.Get("/", chirps_handler.GetAllChirps())
		r.Get("/{id}", chirps_handler.GetOneChirp())
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
