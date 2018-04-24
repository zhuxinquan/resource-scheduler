package slog

import (
	"github.com/cihub/seelog"
)

var (
	httpLogger seelog.LoggerInterface
)

func InitSeelog() {
	logger, err := seelog.LoggerFromConfigAsFile("conf/seelog.xml")

	if err != nil {
		seelog.Critical("err parsing config log file", err)
		return
	}
	seelog.ReplaceLogger(logger)
}
