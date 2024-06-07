package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jmoiron/sqlx"
	"github.com/shulganew/hb.git/internal/config"
	"github.com/shulganew/hb.git/internal/storage"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLog(conf config.Config) zap.SugaredLogger {
	lvl, err := zap.ParseAtomicLevel(conf.ZapLevel)
	if err != nil {
		fmt.Println("Can't set log level: ", err, conf.ZapLevel)
		panic(err)
	}

	var op []string
	var ep []string
	if conf.ZapRunLocal {
		op = []string{"stdout"}
		ep = []string{"stderr"}
	} else {
		op = []string{conf.ZapPath}
		ep = []string{conf.ZapPath}
	}

	cfg := zap.Config{
		Encoding:         "console",
		Level:            lvl,
		OutputPaths:      op,
		ErrorOutputPaths: ep,
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			TimeKey:    "time",
			EncodeTime: zapcore.RFC3339TimeEncoder,

			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	zapLogger := zap.Must(cfg.Build())
	zapLogger.Info("logger construction succeeded")
	zap.ReplaceGlobals(zapLogger)
	defer func() {
		_ = zapLogger.Sync()
	}()

	sugar := *zapLogger.Sugar()

	defer func() {
		_ = sugar.Sync()
	}()
	return sugar
}

// Init context from graceful shutdown. Send to all function for return by
//
//	syscall.SIGINT, syscall.SIGTERM
func InitContext() (ctx context.Context, cancel context.CancelFunc) {
	ctx, cancel = signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	return
}

func InitStore(ctx context.Context, conf config.Config) (stor *storage.Repo, err error) {
	// Connection for Database.
	db, err := sqlx.Connect(config.DataBaseType, conf.DSN)
	if err != nil {
		return nil, err
	}

	// Load storage.
	stor, err = storage.NewRepo(ctx, db)
	if err != nil {
		return nil, err
	}

	zap.S().Infoln("Application init complite")
	return stor, nil

}