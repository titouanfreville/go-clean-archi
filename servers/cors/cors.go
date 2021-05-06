package cors

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

func InitGoCHI(conf Config, router chi.Router) {
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   conf.AllowedOrigins,
		AllowedMethods:   conf.AllowedMethods,
		AllowedHeaders:   conf.AllowedHeaders,
		ExposedHeaders:   conf.ExposedHeaders,
		AllowCredentials: conf.AllowCredentials,
		MaxAge:           int(conf.MaxAge.Seconds()),
	}).Handler)
}
