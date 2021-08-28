package ghost

import (
	log "github.com/sirupsen/logrus"
)

var (
	Debug = log.Debugln
	Debugf = log.Debugf
	Info = log.Infoln
	Infof = log.Infof
	Warn = log.Warningln
	Warnf = log.Warningf
	Error = log.Errorln
	Errorf = log.Errorf
)

func init(){
	log.SetFormatter(&log.JSONFormatter{})
}