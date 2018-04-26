package slog

import (
	"github.com/cihub/seelog"
)

var (
	httpLogger seelog.LoggerInterface
)

func InitSeelog(path string) {
	if path == "" {
		path = "conf/seelog.xml"
	}
	httpLogger, err := seelog.LoggerFromConfigAsFile(path)
	if err != nil {
		seelog.Critical("err parsing config log file", err)
		return
	}
	seelog.ReplaceLogger(httpLogger)
}
