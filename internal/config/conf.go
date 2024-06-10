package config

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
)

const AuthPrefix = "Bearer "
const TokenExp = time.Hour * 3600
const DataBaseType = "postgres"
const Cron = "0 13 * * *"
const Domain = "https://dlearn.ru/hb/login"

type Config struct {
	Bot     string // telegram bot token.
	Address string // Server address and post, ex localhost:8080.
	DSN     string // DB connection.

	ZapPath     string // Logging path.
	ZapLevel    string // Logging level.
	ZapRunLocal bool   // Run for testing on local host, not ptint logs to file.
}

func InitConfig() Config {
	config := Config{}

	// Read OS ENVs.
	b, exist := os.LookupEnv(("HB_BOT_TOKEN"))
	if !exist {
		fmt.Println("ENV HB_BOT_TOKEN not set")
	}
	config.Bot = b

	// Logging level "debug", "info", "warn", "error", "dpanic", "panic", and "fatal".
	level, exist := os.LookupEnv(("HB_ZAP_LEVEL"))
	if !exist {
		fmt.Println("ENV HB_ZAP_LEVEL not set")
	}
	config.ZapLevel = level

	// Logging path.
	lp, exist := os.LookupEnv(("HB_ZAP_PATH"))
	if !exist {
		fmt.Println("ENV HB_ZAP_PATH not set")
	}
	config.ZapPath = lp

	// Run local, no printing logs to file.
	_, exist = os.LookupEnv(("HB_ZAP_LOCAL"))
	if exist {
		config.ZapRunLocal = true
	}

	// Check and parse server URL.
	addr, exist := os.LookupEnv(("HB_ADDRESS"))
	if !exist {
		fmt.Println("ENV HB_ADDRESS not set")
	}

	// Server address.
	config.Address = addr

	// DSN postgres URI.
	dsn, exist := os.LookupEnv(("HB_DSN_URI"))
	if !exist {
		fmt.Println("ENV HB_DSN_URI not set")
	}
	config.DSN = dsn

	zap.S().Debugf("Config: %#v \n", config)
	return config
}
