package server

import (
	"context"
	"net/http"

	conf "github.com/boginskiy/Clicki/cmd/config"

	"github.com/boginskiy/Clicki/internal/audit"
	"github.com/boginskiy/Clicki/internal/auth"
	"github.com/boginskiy/Clicki/internal/logg"
	mv "github.com/boginskiy/Clicki/internal/middleware"
	prep "github.com/boginskiy/Clicki/internal/preparation"
	repo "github.com/boginskiy/Clicki/internal/repository"
	"github.com/boginskiy/Clicki/internal/router"
	"github.com/boginskiy/Clicki/internal/service"
	"github.com/boginskiy/Clicki/internal/validation"
)

func Run(config conf.Config, baseLog logg.Logger, repo repo.Repository) {
	// Loggers
	infraLogg := logg.NewLogg(config.GetLogFile(), "INFO")
	authLogg := logg.NewLogg("reg.log", "ERROR")

	// Audit
	sub1 := audit.NewFileReceiver(infraLogg, config.GetAuditFile(), 1)
	sub2 := audit.NewServerReceiver(infraLogg, config.GetAuditURL(), 2)
	publisher := audit.NewPublish(sub1, sub2)

	// Middleware & Registr
	// audit := audit.NewAudit(config, baseLog, publisher)
	auther := auth.NewAuth(config, authLogg, repo)
	midWare := mv.NewMiddleware(infraLogg, auther)

	// Extra
	extraFuncer := prep.NewExtraFunc()
	checker := validation.NewChecker()

	// Ctx
	ctx, cancel := context.WithCancel(context.Background())

	// Services
	CoreServ := service.NewCoreService(config, baseLog, repo, publisher)
	APIShortURL := service.NewAPIShortURL(CoreServ, repo, checker, extraFuncer)
	ShortURL := service.NewShortURL(CoreServ, repo, checker, extraFuncer)
	APIDelMess := service.NewDelMess(ctx, CoreServ, repo)

	// writing log...
	baseLog.RaiseInfo(logg.StartedServInfo, logg.Fields{"port": config.GetSrvAddr()})

	// Start server
	err := http.ListenAndServe(
		config.GetSrvAddr(), router.Router(midWare, APIShortURL, ShortURL, APIDelMess))

	// writing log...
	baseLog.RaiseFatal(err, logg.StartedServFatal, logg.Fields{"port": config.GetSrvAddr()})

	// defer
	defer infraLogg.Close()
	defer authLogg.Close()
	defer sub1.Clouse()
	defer sub2.Clouse()
	defer cancel()
}
