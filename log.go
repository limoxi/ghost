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
	Panic = log.Panicln
	Panicf = log.Panicf
)

type panicHook struct{
	
}

func (this *panicHook) Fire(entry *log.Entry) error {
	return nil
}

func (this *panicHook) Levels() []log.Level {
	return []log.Level{log.PanicLevel}
}

func init(){
	log.SetFormatter(&log.JSONFormatter{})
	log.AddHook(&panicHook{})
}