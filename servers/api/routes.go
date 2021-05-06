package api

import (
	"github.com/go-chi/chi"
)

func (transport *API) initRoutes(r chi.Router) {
	r.Get("/uptime", transport.endpoints.Uptime)

	r.NotFound(transport.endpoints.NotFound)
}
