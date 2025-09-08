package server

import (
	"net/http"

	c "github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/db"
	"github.com/boginskiy/Clicki/internal/db2"
	l "github.com/boginskiy/Clicki/internal/logger"
	m "github.com/boginskiy/Clicki/internal/middleware"
	r "github.com/boginskiy/Clicki/internal/router"
)

func Run(kwargs c.VarGetter, logger l.Logger, db2 db2.DBConnecter) {
	// Info Logger
	infoLog := l.NewLogg(kwargs.GetLogFile(), "INFO")
	defer infoLog.CloseDesc()

	// Middleware
	midWare := m.NewMiddleware(infoLog)

	// Db
	writerFile, err := db.NewFileWorking(kwargs.GetPathToStore())
	if err != nil {
		logger.RaiseError("Run", l.Fields{"error": err.Error()})
	}
	db := db.NewDBStore(writerFile)
	defer writerFile.Close()

	// writing log...
	infoLog.RaiseInfo(l.StartedServInfo,
		l.Fields{"port": kwargs.GetSrvAddr()})

	// Start server
	err = http.ListenAndServe(kwargs.GetSrvAddr(), r.Router(kwargs, logger, midWare, db, db2))

	// writing log...
	if err != nil {
		logger.RaiseFatal(l.StartedServFatal,
			l.Fields{"port": kwargs.GetSrvAddr()})
	}
}
