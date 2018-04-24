package main

import (
	"common/slog"
	"github.com/cihub/seelog"
)

func main() {
	slog.InitSeelog()
	defer seelog.Flush()
	seelog.Debug("123")
}
