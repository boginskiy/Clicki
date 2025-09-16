package server

import (
	"net/http"

	c "github.com/boginskiy/Clicki/cmd/config"

	l "github.com/boginskiy/Clicki/internal/logger"
	m "github.com/boginskiy/Clicki/internal/middleware"
	p "github.com/boginskiy/Clicki/internal/preparation"
	rp "github.com/boginskiy/Clicki/internal/repository"
	r "github.com/boginskiy/Clicki/internal/router"
	s "github.com/boginskiy/Clicki/internal/service"
	v "github.com/boginskiy/Clicki/internal/validation"
)

func Run(kwargs c.VarGetter, baseLog l.Logger, repo rp.URLRepository) {
	// Info Logger
	infoLog := l.NewLogg(kwargs.GetLogFile(), "INFO")
	defer infoLog.CloseDesc()

	// Middleware
	midWare := m.NewMiddleware(infoLog)

	// Extra
	extraFuncer := p.NewExtraFunc() // extraFuncer - дополнительные функции
	checker := v.NewChecker()       // checker - валидация данных

	// Services
	APIShortURL := s.NewAPIShortURL(kwargs, baseLog, repo, checker, extraFuncer)
	ShortURL := s.NewShortURL(kwargs, baseLog, repo, checker, extraFuncer)

	// writing log...
	baseLog.RaiseInfo(l.StartedServInfo, l.Fields{"port": kwargs.GetSrvAddr()})

	// Start server
	err := http.ListenAndServe(kwargs.GetSrvAddr(), r.Router(midWare, APIShortURL, ShortURL))

	// writing log...
	baseLog.RaiseFatal(err, l.StartedServFatal, l.Fields{"port": kwargs.GetSrvAddr()})

}
