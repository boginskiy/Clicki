package logg

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

// Logg - custom logger above logrus.
type Logg struct {
	Log      *logrus.Logger
	Desc     *os.File
	NameFile string
}

func NewLogg(nameFile, level string) *Logg {
	// Create file.
	tmpDesc := createLogFile(nameFile)
	// Settings Logrus.
	tmpLogrus := setupLogrus(tmpDesc, LEVEL[level])

	return &Logg{
		NameFile: nameFile,
		Desc:     tmpDesc,
		Log:      tmpLogrus,
	}
}

func (e *Logg) Close() {
	e.Desc.Close()
}

func (e *Logg) RaiseInfo(msg string, dataMap Fields) {
	e.Log.WithFields(logrus.Fields(dataMap)).Info(msg)
}

func (e *Logg) RaiseWarn(msg string, dataMap Fields) {
	e.Log.WithFields(logrus.Fields(dataMap)).Warn(msg)
}

func (e *Logg) RaiseError(err error, msg string, dataMap Fields) {
	if err != nil {
		fmt.Fprintln(os.Stdout, msg)

		if dataMap != nil {
			e.Log.WithFields(logrus.Fields(dataMap)).Error(msg)
		} else {
			e.Log.WithFields(logrus.Fields(Fields{"error": err.Error()})).Error(msg)
		}
	}
}

func (e *Logg) RaiseFatal(err error, msg string, dataMap Fields) {
	if err != nil {
		fmt.Fprintln(os.Stdout, msg)

		if dataMap != nil {
			e.Log.WithFields(logrus.Fields(dataMap)).Fatal(msg)
		} else {
			e.Log.WithFields(logrus.Fields(Fields{"fatal": err.Error()})).Fatal(msg)
		}
	}
}

func (e *Logg) RaisePanic(err error, msg string, dataMap Fields) {
	if err != nil {
		fmt.Fprintln(os.Stdout, msg)
		e.Log.WithFields(logrus.Fields(dataMap)).Panic(msg)
	}
}
