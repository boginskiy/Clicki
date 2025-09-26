package server

import (
	"net/http"

	conf "github.com/boginskiy/Clicki/cmd/config"

	"github.com/boginskiy/Clicki/internal/logg"
	midw "github.com/boginskiy/Clicki/internal/middleware"
	prep "github.com/boginskiy/Clicki/internal/preparation"
	regstr "github.com/boginskiy/Clicki/internal/register"
	"github.com/boginskiy/Clicki/internal/repository"
	route "github.com/boginskiy/Clicki/internal/router"
	srv "github.com/boginskiy/Clicki/internal/service"
	valid "github.com/boginskiy/Clicki/internal/validation"
)

func Run(kwargs conf.VarGetter, baseLog logg.Logger, repo repository.Repository) {
	// Special Loggers for middleware, registration
	midWareLogger := logg.NewLogg(kwargs.GetLogFile(), "INFO")
	registLogger := logg.NewLogg("LogRegister.log", "ERROR")

	defer midWareLogger.CloseDesc()
	defer registLogger.CloseDesc()

	// Middleware & Registr
	register := regstr.NewRegist(kwargs, registLogger, repo)
	midWare := midw.NewMiddleware(midWareLogger, register, repo)

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
