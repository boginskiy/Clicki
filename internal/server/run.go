package server

import (
	"context"
	"net/http"

	conf "github.com/boginskiy/Clicki/cmd/config"

	"github.com/boginskiy/Clicki/internal/audit"
	auth "github.com/boginskiy/Clicki/internal/auther"
	"github.com/boginskiy/Clicki/internal/logg"
	midw "github.com/boginskiy/Clicki/internal/middleware"
	prep "github.com/boginskiy/Clicki/internal/preparation"
	repo "github.com/boginskiy/Clicki/internal/repository"
	route "github.com/boginskiy/Clicki/internal/router"
	srv "github.com/boginskiy/Clicki/internal/service"
	valid "github.com/boginskiy/Clicki/internal/validation"
)

func Run(kwargs conf.VarGetter, baseLog logg.Logger, repo repo.Repository) {
	// Loggers
	infraLogg := logg.NewLogg(kwargs.GetLogFile(), "INFO")
	authLogg := logg.NewLogg("LogRegister.log", "ERROR")

	// Audit
	sub1 := audit.NewFileReceiver(infraLogg, kwargs.GetAuditFile(), 1)
	sub2 := audit.NewServerReceiver(infraLogg, kwargs.GetAuditURL(), 2)
	publisher := audit.NewPublish(sub1, sub2)

	// Middleware & Registr
	// audit := audit.NewAudit(kwargs, baseLog, publisher)
	auther := auth.NewAuth(kwargs, authLogg, repo)
	midWare := midw.NewMiddleware(infraLogg, auther)

	// Extra
	extraFuncer := prep.NewExtraFunc()
	checker := valid.NewChecker()

	// Ctx
	ctx, cancel := context.WithCancel(context.Background())

	// Services
	CoreServ := srv.NewCoreService(kwargs, baseLog, repo, publisher)
	APIShortURL := srv.NewAPIShortURL(CoreServ, repo, checker, extraFuncer)
	ShortURL := srv.NewShortURL(CoreServ, repo, checker, extraFuncer)
	APIDelMess := srv.NewDelMess(ctx, CoreServ, repo)

	// writing log...
	baseLog.RaiseInfo(logg.StartedServInfo, logg.Fields{"port": kwargs.GetSrvAddr()})

	// Start server
	err := http.ListenAndServe(
		kwargs.GetSrvAddr(), route.Router(midWare, APIShortURL, ShortURL, APIDelMess))

	// writing log...
	baseLog.RaiseFatal(err, logg.StartedServFatal, logg.Fields{"port": kwargs.GetSrvAddr()})

	// defer
	defer infraLogg.CloseDesc()
	defer authLogg.CloseDesc()
	defer sub1.Clouse()
	defer sub2.Clouse()
	defer cancel()
}
