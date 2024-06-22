package main

import (
	_ "embed"
	"fmt"
	"github.com/jimu-server/config"
	_ "github.com/jimu-server/gpt"
	"github.com/jimu-server/logger"
	"github.com/jimu-server/web"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"

	"syscall"
)

// assistant  ollama 独立服务,用于独立编译嵌入 Electron 运行时候启动

func main() {

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", "0.0.0.0", config.Evn.App.Ollama.Port),
		Handler: web.Engine,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			panic(err.Error())
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-signals:
		if err := zap.L().Sync(); err != nil {
			logger.Logger.Error("sync zap log error", zap.Error(err))
		}
		if err := server.Close(); err != nil {
			logger.Logger.Error("close server error", zap.Error(err))
		}
		logger.Logger.Info("server shutdown")
	}
}
