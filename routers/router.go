package routers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/maltegrosse/meteocrop-jp-api/controllers"
	"github.com/maltegrosse/meteocrop-jp-api/log"
	l "github.com/treastech/logger"
	"net/http"
)

func Routes() *chi.Mux {

	router := chi.NewRouter()
	cors_ := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "User", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	router.Use(
		render.SetContentType(render.ContentTypeJSON),
		l.Logger(log.Logger),
		middleware.RealIP,
		middleware.Recoverer,
		cors_.Handler,
	)
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(""))
	})
	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})
	router.Get("/weather", controllers.WeatherObservations)

	return router
}
