package logger

import (
	"fmt"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

type Fields map[string]any

type Logger interface {
	RaiseInfo(string, Fields)
	RaiseWarn(string, Fields)
	RaiseError(error, string, Fields)
	RaiseFatal(error, string, Fields)
	RaisePanic(error, string, Fields)
	CloseDesc()
}

var LEVEL = map[string]logrus.Level{
	"DEBUG": logrus.DebugLevel,
	"INFO":  logrus.InfoLevel,
	"WARN":  logrus.WarnLevel,
	"ERROR": logrus.ErrorLevel,
	"FATAL": logrus.FatalLevel,
	"PANIC": logrus.PanicLevel,
}

type Logg struct {
	Log      *logrus.Logger
	mu       sync.Mutex
	Desc     *os.File
	NameFile string
}

func NewLogg(nameFile, level string) *Logg {
	// Create file
	tmpDesc := createLogFile(nameFile)
	// Settings Logrus
	tmpLogrus := setupLogrus(tmpDesc, LEVEL[level])

	return &Logg{
		NameFile: nameFile,
		Desc:     tmpDesc,
		Log:      tmpLogrus,
	}
}

func (e *Logg) CloseDesc() {
	e.Desc.Close()
}

func (e *Logg) RaiseInfo(msg string, dataMap Fields) {
	e.mu.Lock()
	e.Log.WithFields(logrus.Fields(dataMap)).Info(msg)
	e.mu.Unlock()
}

func (e *Logg) RaiseWarn(msg string, dataMap Fields) {
	e.mu.Lock()
	e.Log.WithFields(logrus.Fields(dataMap)).Warn(msg)
	e.mu.Unlock()
}

func (e *Logg) RaiseError(err error, msg string, dataMap Fields) {
	if err != nil {
		e.mu.Lock()
		fmt.Fprintln(os.Stdout, msg)
		e.Log.WithFields(logrus.Fields(dataMap)).Error(msg)
		e.mu.Unlock()
	}
}

func (e *Logg) RaiseFatal(err error, msg string, dataMap Fields) {
	if err != nil {
		e.mu.Lock()
		fmt.Fprintln(os.Stdout, msg)
		e.Log.WithFields(logrus.Fields(dataMap)).Fatal(msg)
		e.mu.Unlock()
	}
}

func (e *Logg) RaisePanic(err error, msg string, dataMap Fields) {
	if err != nil {
		e.mu.Lock()
		fmt.Fprintln(os.Stdout, msg)
		e.Log.WithFields(logrus.Fields(dataMap)).Panic(msg)
		e.mu.Unlock()
	}
}
