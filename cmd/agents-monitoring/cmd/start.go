package cmd

import (
	"flag"
	"go.uber.org/zap"
)

func Execute() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("... Starting Agents Monitoring Service ...")

	httpServerPort := flag.String("http-server-port", "8080", "http server port")
	flag.Parse()

	logger.Info("Init Variables",
		zap.String("http-server-port", *httpServerPort),
	)
}
