package logg

import "github.com/sirupsen/logrus"

const (
	DataReqResInfo   = "[INFO]:The data request and response"
	StartedServInfo  = "[INFO]:The server has started"
	StartedServFatal = "[FATAL]:The server has not started"
)

type Fields map[string]any

var LEVEL = map[string]logrus.Level{
	"DEBUG": logrus.DebugLevel,
	"INFO":  logrus.InfoLevel,
	"WARN":  logrus.WarnLevel,
	"ERROR": logrus.ErrorLevel,
	"FATAL": logrus.FatalLevel,
	"PANIC": logrus.PanicLevel,
}
