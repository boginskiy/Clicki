package server

import (
	"net/http"

	c "github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/db"
	"github.com/boginskiy/Clicki/internal/db2"
	l "github.com/boginskiy/Clicki/internal/logger"
	m "github.com/boginskiy/Clicki/internal/middleware"
	p "github.com/boginskiy/Clicki/internal/preparation"
	r "github.com/boginskiy/Clicki/internal/router"
	s "github.com/boginskiy/Clicki/internal/service"
	v "github.com/boginskiy/Clicki/internal/validation"
)

func Run(kwargs c.VarGetter, baseLog l.Logger, db2 db2.DBConnecter) {
	// Info Logger
	infoLog := l.NewLogg(kwargs.GetLogFile(), "INFO")
	defer infoLog.CloseDesc()

	// Middleware
	midWare := m.NewMiddleware(infoLog)

	// Db
	writerFile, err := db.NewFileWorking(kwargs.GetPathToStore())
	baseLog.RaiseError(err, "Run", nil)
	db := db.NewDBStore(writerFile)
	defer writerFile.Close()

	// Extra
	extraFuncer := p.NewExtraFunc() // extraFuncer - дополнительные функции
	checker := v.NewChecker()       // checker - валидация данных

	// Services
	APIShortURL := s.NewAPIShortURL(db, db2, baseLog, checker, extraFuncer)
	ShortURL := s.NewShortURL(db, db2, baseLog, checker, extraFuncer)

	// writing log...
	baseLog.RaiseInfo(l.StartedServInfo, l.Fields{"port": kwargs.GetSrvAddr()})

	// Start server
	err = http.ListenAndServe(kwargs.GetSrvAddr(),
		r.Router(kwargs, midWare, APIShortURL, ShortURL))

	// writing log...
	baseLog.RaiseFatal(err, l.StartedServFatal, l.Fields{"port": kwargs.GetSrvAddr()})

}
