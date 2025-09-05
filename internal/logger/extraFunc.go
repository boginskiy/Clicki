package logger

import (
	"io"
	"log"
	"os"

	"github.com/sirupsen/logrus"
)

func createLogFile(name string) *os.File {
	file, err := os.OpenFile(name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	return file
}

func setupLogrus(w io.Writer, level logrus.Level) *logrus.Logger {
	l := logrus.New() // Create Lorgus
	l.SetOutput(w)    // Add io.Writer
	l.SetLevel(level) // Add logrus.Level
	return l
}
