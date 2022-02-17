package main

import (
	"context"
	"github.com/aibotsoft/ntp-service/pkg/config"
	"github.com/aibotsoft/ntp-service/pkg/logger"
	"github.com/aibotsoft/ntp-service/pkg/signals"
	"github.com/aibotsoft/ntp-service/pkg/version"
	"github.com/aibotsoft/ntp-service/services/ntp"
	"go.uber.org/zap"
)

func main() {
	cfg := config.NewConfig()
	log, err := logger.NewLogger(cfg.Zap.Level, cfg.Zap.Encoding, cfg.Zap.Caller)
	if err != nil {
		panic(err)
	}
	log.Info("start_service", zap.Any("config", cfg), zap.String("version", version.Version))
	ctx, cancel := context.WithCancel(context.Background())
	n := ntp.NewNtp(cfg, log, ctx)
	errCh := make(chan error)
	go func() {
		errCh <- n.Run()
	}()
	defer func() {
		log.Info("closing_services...")
		cancel()
		err := n.Close()
		if err != nil {
			log.Error("close_ntp_error", zap.Error(err))
		}
		_ = log.Sync()
	}()
	stopCh := signals.SetupSignalHandler()
	select {
	case err := <-errCh:
		log.Error("stop_service_by_error", zap.Error(err))
	case sig := <-stopCh:
		log.Info("stop_service_by_os_signal", zap.String("signal", sig.String()))
	}
}
