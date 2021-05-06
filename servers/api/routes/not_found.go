package routes

import (
	"net/http"
)

func (ep *Endpoints) NotFound(w http.ResponseWriter, r *http.Request) {
	ep.logger.Debug("not found")
	w.Write([]byte("<html><head><title>Vous Etes Perdu ?</title></head><body><h1>Perdu sur l'Internet ?</h1><h2>Pas de panique, on va vous aider</h2><strong><pre>    * <----- vous &ecirc;tes ici</pre></strong></body></html>"))
}
