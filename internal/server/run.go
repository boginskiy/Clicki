package server

import (
	"net/http"

	c "github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/db"
	l "github.com/boginskiy/Clicki/internal/logger"
	m "github.com/boginskiy/Clicki/internal/middleware"
	p "github.com/boginskiy/Clicki/internal/preparation"
	r "github.com/boginskiy/Clicki/internal/router"
	s "github.com/boginskiy/Clicki/internal/service"
	v "github.com/boginskiy/Clicki/internal/validation"
)

func Run(kwargs c.VarGetter) {
	// Logger
	fatalLog := l.NewLogg(kwargs.GetNameLogFatal(), "ERROR")
	infoLog := l.NewLogg(kwargs.GetNameLogInfo(), "INFO")
	defer fatalLog.CloseDesc()
	defer infoLog.CloseDesc()

	// Middleware
	midWare := m.NewMiddleware(infoLog)

	// Db
	writerFile, err := db.NewFileWorking(kwargs.GetPathToStore())
	fatalLog.RaiseError(err, "Run", l.Fields{"error": err.Error()})
	db := db.NewDBStore(writerFile)
	defer writerFile.Close()

	// Extra
	extraFuncer := p.NewExtraFunc() // extraFuncer - дополнительные функции
	checker := v.NewChecker()       // checker - валидация данных

	// Services
	APIShortURL := s.NewAPIShortURL(db, fatalLog, checker, extraFuncer)
	ShortURL := s.NewShortURL(db, fatalLog, checker, extraFuncer)

	// writing log...
	infoLog.RaiseInfo(l.StartedServInfo, l.Fields{"port": kwargs.GetSrvAddr()})

	// Start server
	err = http.ListenAndServe(kwargs.GetSrvAddr(),
		r.Router(kwargs, midWare, APIShortURL, ShortURL))

	// writing log...
	fatalLog.RaiseFatal(err, l.StartedServFatal, l.Fields{"port": kwargs.GetSrvAddr()})
}
