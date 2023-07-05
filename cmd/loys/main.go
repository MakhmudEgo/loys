package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"dz1/internal/app"
	"dz1/internal/infrastructure/config"
)

func main() {
	// psw := "sdkjfksdf"
	// password, err := bcrypt.GenerateFromPassword([]byte(psw), 3)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// err = bcrypt.CompareHashAndPassword(password, []byte(psw))
	// log.Fatal(err)

	// d := "2023-01-01"
	//
	// parse, err := time.Parse(time.DateOnly, d)
	//
	// log.Fatal(parse.String(), err)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var configPath string

	flag.StringVar(&configPath, "config", "configs/config.yml", "server configuration file path")
	flag.Parse()

	if err := run(ctx, configPath); err != nil {
		log.Fatalln(err)
	}
}

func run(ctx context.Context, cfgPath string) error {
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("load config file:%w", err)
	}

	logger, err := configureLogger(&cfg.Log)
	if err != nil {
		return fmt.Errorf("configure logger: %w", err)
	}
	defer func() { _ = logger.Sync() }()

	if err := app.New(cfg).Run(ctx, logger); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}

	return nil
}

func configureLogger(cfg *config.Log) (*zap.Logger, error) {
	zapCfg := zap.NewProductionEncoderConfig()
	zapCfg.EncodeTime = zapcore.RFC3339TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(zapCfg)
	logFile := os.Stdout
	if cfg.File != "" {
		var err error
		logFile, err = os.OpenFile(cfg.File, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
	}

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, logFile, cfg.Level),
	)

	return zap.New(core, zap.AddCaller()), nil
}
