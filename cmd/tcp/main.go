package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"kv_db/config"
	"kv_db/internal/initialization"
)

var ConfigFileName = os.Getenv("CONFIG_FILE_NAME")

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	cfg := config.Config{}
	if ConfigFileName != "" {
		var err error
		cfg, err = config.Load(ConfigFileName)
		if err != nil {
			return err
		}
	}
	file, err := initialization.CreateLogFile(cfg.Logging.Output)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Println(err)
		}
	}()

	initializer, err := initialization.NewInitializer(cfg, file)
	if err != nil {
		return err
	}

	if err = initializer.Start(ctx); err != nil {
		return err
	}
	return nil
}
