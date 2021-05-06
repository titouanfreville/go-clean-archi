package main

import (
	"time"

	"go.uber.org/fx"

	"go-clean-archi/servers/api"
	"go-clean-archi/services/fxapp"
)

const timeout = 30 * time.Second

func main() {
	app := fx.New(
		fx.NopLogger, // remove for debug

		api.Transport,

		// fx.Invoke(validation.ValidateConfig),
		fx.Invoke(api.Run),
	)

	fxapp.Start(app, timeout)
	<-app.Done()
	fxapp.Shutdown(app, timeout)
}
