package main

import (
	router "github.com/shulganew/hb.git/internal/api"
	"github.com/shulganew/hb.git/internal/api/oapi"
	"github.com/shulganew/hb.git/internal/app"
	"github.com/shulganew/hb.git/internal/config"
	"github.com/shulganew/hb.git/internal/services"
	"go.uber.org/zap"
)

func main() {

	// Get application config.
	conf := config.InitConfig()

	// Init application logging.
	app.InitLog(conf)

	// Root app context.
	ctx, cancel := app.InitContext()

	// Error channel.
	componentsErrs := make(chan error, 1)

	// Init Repo
	stor, err := app.InitStore(ctx, conf)
	if err != nil {
		zap.S().Fatalln(err)
	}

	swagger, err := oapi.GetSwagger()
	if err != nil {
		zap.S().Fatalln(err)
	}
	// Clear out the servers array in the swagger spec, that skips validating
	// that server names match. We don't know how this thing will be run.
	swagger.Servers = nil

	// Create router.
	rt := router.Route(conf, swagger)

	hs := services.NewHappy(ctx, stor, conf)

	// We now register our GophKeeper above as the handler for the interface.
	oapi.HandlerFromMux(hs, rt)

	// Start web server.
	restDone := app.StartAPI(ctx, conf, componentsErrs, rt)

	botDone := app.StartBot(ctx, conf, componentsErrs)

	// Graceful shutdown.
	app.Graceful(ctx, cancel, componentsErrs)

	// Waiting http server shuting down.
	<-restDone
	// Telega shutdown.
	<-botDone

	zap.S().Infoln("App done.")
}
