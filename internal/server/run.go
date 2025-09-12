package server

import (
	"log"
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

	err := repo.GetDB().Ping()
	if err != nil {
		log.Println(">>2> Database connection is closed:", err)
	} else {
		log.Println(">>2> Database connection is active.")
	}

	// Extra
	extraFuncer := p.NewExtraFunc() // extraFuncer - дополнительные функции
	checker := v.NewChecker()       // checker - валидация данных

	// Services
	APIShortURL := s.NewAPIShortURL(repo, baseLog, checker, extraFuncer)
	ShortURL := s.NewShortURL(repo, baseLog, checker, extraFuncer)

	// writing log...
	baseLog.RaiseInfo(l.StartedServInfo, l.Fields{"port": kwargs.GetSrvAddr()})

	// Start server
	err := http.ListenAndServe(kwargs.GetSrvAddr(),
		r.Router(kwargs, midWare, APIShortURL, ShortURL))

	// writing log...
	baseLog.RaiseFatal(err, l.StartedServFatal, l.Fields{"port": kwargs.GetSrvAddr()})

}

// TODO:
// Проверка работы ,проверка работы флагов, записи в верные БД
