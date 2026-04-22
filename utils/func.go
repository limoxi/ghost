package utils

import (
	log "github.com/sirupsen/logrus"
)

func RunInGoroutine(f func()) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorln(err)
		}
	}()

	go f()
}
