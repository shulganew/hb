package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	router "github.com/shulganew/hb.git/internal/api"
	"github.com/shulganew/hb.git/internal/api/oapi"
	"github.com/shulganew/hb.git/internal/app"
	"github.com/shulganew/hb.git/internal/config"
	"github.com/shulganew/hb.git/internal/services"
	"github.com/shulganew/hb.git/internal/storage/mem"
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

	// Init mem storage.
	mem := mem.NewMemory()

	// BOT service.
	// Create new bot.
	b, err := tgbotapi.NewBotAPI(conf.Bot)
	if err != nil {
		panic(err)
	}

	bs := services.NewBot(ctx, stor, conf, mem)

	swagger, err := oapi.GetSwagger()
	if err != nil {
		zap.S().Fatalln(err)
	}
	// Clear out the servers array in the swagger spec, that skips validating
	// that server names match. We don't know how this thing will be run.
	swagger.Servers = nil

	// Create router.
	rt := router.Route(conf, swagger)

	// API service.
	hs := services.NewHappy(ctx, stor, conf, mem)

	// We now register our GophKeeper above as the handler for the interface.
	oapi.HandlerFromMux(hs, rt)

	// Start web server.
	restDone := app.StartAPI(ctx, conf, componentsErrs, rt)

	botDone := app.StartBot(ctx, conf, b, bs, componentsErrs)

	// Start cron.
	c := app.InitCron(ctx, b, bs)

	// Graceful shutdown.
	app.Graceful(ctx, cancel, componentsErrs)

	// Waiting http server shuting down.
	<-restDone
	// Telega shutdown.
	<-botDone
	// Stop scheduler.
	c.Stop()
	zap.S().Infoln("App done.")
}
