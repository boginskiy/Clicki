package server

import (
	"net/http"

	conf "github.com/boginskiy/Clicki/cmd/config"

	"github.com/boginskiy/Clicki/internal/logg"
	midw "github.com/boginskiy/Clicki/internal/middleware"
	prep "github.com/boginskiy/Clicki/internal/preparation"
	repo "github.com/boginskiy/Clicki/internal/repository"
	route "github.com/boginskiy/Clicki/internal/router"
	srv "github.com/boginskiy/Clicki/internal/service"
	valid "github.com/boginskiy/Clicki/internal/validation"
)

func Run(kwargs conf.VarGetter, baseLog logg.Logger, repo repo.Repository) {
	// Info Logger
	infoLog := logg.NewLogg(kwargs.GetLogFile(), "INFO")
	defer infoLog.CloseDesc()

	// Middleware
	midWare := midw.NewMiddleware(infoLog)

	// Extra
	extraFuncer := prep.NewExtraFunc() // extraFuncer - дополнительные функции
	checker := valid.NewChecker()      // checker - валидация данных

	// Services
	APIShortURL := srv.NewAPIShortURL(kwargs, baseLog, repo, checker, extraFuncer)
	ShortURL := srv.NewShortURL(kwargs, baseLog, repo, checker, extraFuncer)

	// writing log...
	baseLog.RaiseInfo(logg.StartedServInfo, logg.Fields{"port": kwargs.GetSrvAddr()})

	// Start server
	err := http.ListenAndServe(kwargs.GetSrvAddr(), route.Router(midWare, APIShortURL, ShortURL))

	// writing log...
	baseLog.RaiseFatal(err, logg.StartedServFatal, logg.Fields{"port": kwargs.GetSrvAddr()})

}
