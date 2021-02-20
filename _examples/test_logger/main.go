package main

import (
	"github.com/cherry-game/cherry/logger"
	"go.uber.org/zap"
)

func main() {

	zap.NewProduction()

	config := &cherryLogger.Config{
		Level:           "debug",
		StackLevel:      "error",
		EnableWriteFile: false,
		EnableConsole:   true,
		FilePath:        "",
		MaxSize:         0,
		MaxAge:          0,
		MaxBackups:      0,
		Compress:        false,
		TimeFormat:      "",
		PrintCaller:     false,
	}

	logger := cherryLogger.NewConfigLogger(config)

	logger.Info("111111111111111111111111111111")
	logger.Debugf("aaaaaaaaaaaaaa %s", "aaaaa args.......")
	logger.Infow("failed to fetch URL.", "url", "http://example.com")
	logger.Infow("failed to fetch URL.",
		"url", "http://example.com",
		"name", "url name",
	)
	logger.Warnw("failed to fetch URL.",
		"url", "http://example.com",
		"name", "url name",
	)
	logger.Errorw("failed to fetch URL.",
		"url", "http://example.com",
		"name", "url name",
	)
	logger.Fatal("fatal fatal fatal fatal fatal")

}
