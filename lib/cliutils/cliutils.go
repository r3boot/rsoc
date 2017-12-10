package cliutils

import "github.com/r3boot/rsoc/lib/logger"

var log *logger.Logger

func Setup(l *logger.Logger) {
	log = l
}
