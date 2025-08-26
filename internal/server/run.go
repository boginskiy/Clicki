package server

import (
	"net/http"

	c "github.com/boginskiy/Clicki/cmd/config"
	db "github.com/boginskiy/Clicki/internal/db"
	l "github.com/boginskiy/Clicki/internal/logger"
	m "github.com/boginskiy/Clicki/internal/middleware"
	p "github.com/boginskiy/Clicki/internal/preparation"
	r "github.com/boginskiy/Clicki/internal/router"
	s "github.com/boginskiy/Clicki/internal/service"
	v "github.com/boginskiy/Clicki/internal/validation"
)

func Run(kwargs c.VarGetter) {
	// Info Logger
	infoLog := l.NewLogg(kwargs.GetNameLogInfo(), "INFO")
	// Fatal Logger
	fatalLog := l.NewLogg(kwargs.GetNameLogFatal(), "FATAL")

	defer fatalLog.CloseDesc()
	defer infoLog.CloseDesc()

	// Middleware
	midWare := m.NewMiddleware(infoLog)

	// Business logic
	extraFuncer := p.NewExtraFunc() // extraFuncer - дополнительные возможности
	checker := v.NewChecker()       // checker - валидация данных
	db := db.NewDBStore()           // db - слой базы данных 'DBStore'

	// Service 'ShorteningURL'
	shortingURL := s.NewShorteningURL(db, checker, extraFuncer, fatalLog)

	// writing log...
	infoLog.RaiseInfo(
		l.StartedServInfo,
		l.Fields{"port": kwargs.GetSrvAddr()},
	)

	// Start server
	err := http.ListenAndServe(kwargs.GetSrvAddr(),
		r.Router(kwargs, midWare, shortingURL))

	// writing log...
	if err != nil {
		fatalLog.RaiseFatal(
			l.StartedServFatal,
			l.Fields{"port": kwargs.GetSrvAddr()},
		)
	}
}
