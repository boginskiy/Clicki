package app

import (
	"context"

	"github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/cmd/server"
	"github.com/boginskiy/Clicki/internal/audit"
	"github.com/boginskiy/Clicki/internal/auth"
	"github.com/boginskiy/Clicki/internal/handler"
	"github.com/boginskiy/Clicki/internal/layers"
	"github.com/boginskiy/Clicki/internal/logg"
	mv "github.com/boginskiy/Clicki/internal/middleware"
	prep "github.com/boginskiy/Clicki/internal/preparation"
	"github.com/boginskiy/Clicki/internal/protocol"
	"github.com/boginskiy/Clicki/internal/router"
	"github.com/boginskiy/Clicki/internal/service"
	"github.com/boginskiy/Clicki/internal/validation"
)

type App struct {
	Cfg  config.Config
	Logg logg.Logger
}

func NewApp(config config.Config, logger logg.Logger) *App {
	return &App{
		Cfg:  config,
		Logg: logger,
	}
}

func (a *App) close() {
	a.Logg.Close()
}

func (a *App) Start() {
	// DB & Repo.
	setupLayers := layers.NewLayers(a.Cfg, a.Logg)
	database := setupLayers.NewLayerDB()
	repository := setupLayers.NewLayerRepo(database)

	// Loggers.
	infraLogg := logg.NewLogg(a.Cfg.GetLogFile(), "INFO")
	authLogg := logg.NewLogg("auth.log", "ERROR")

	// Audit.
	sub1 := audit.NewFileReceiver(infraLogg, a.Cfg.GetAuditFile(), 1)
	sub2 := audit.NewServerReceiver(infraLogg, a.Cfg.GetAuditURL(), 2)
	publisher := audit.NewPublish(sub1, sub2)

	// Middleware & Auth.
	auther := auth.NewAuth(a.Cfg, authLogg, repository)
	middleware := mv.NewMdlwere(a.Cfg, infraLogg, auther)

	checker := validation.NewChecker()  // Checker for validation.
	funcer := prep.NewFunctions(a.Logg) // Funcer for extra main function.

	// Context.
	ctx, cancel := context.WithCancel(context.Background())

	// Services.
	APIURLServ := service.NewAPIURLServ(a.Cfg, a.Logg, repository, checker, funcer, publisher)
	URLServ := service.NewURLServ(a.Cfg, a.Logg, repository, checker, funcer, publisher)
	APIDelServ := service.NewAPIDelServ(ctx, a.Cfg, a.Logg, repository)

	// Handlers.
	protHTTP := protocol.NewProtocolHTTP(funcer)                                     // Test
	APIURLHdler := handler.NewAPIURLHandlers(APIURLServ, APIDelServ, protHTTP) // Test

	URLHdler := handler.NewURLHandlers(URLServ)
	PprofHdler := handler.NewPprofHandlers()

	// Router.
	router := router.NewRoute(URLHdler, APIURLHdler, PprofHdler)

	// gRPC.
	autherGRPC := auth.NewAuthGRPC(a.Cfg, authLogg, repository)
	interceptor := mv.NewIntercept(a.Cfg, infraLogg, autherGRPC)
	protGRPC := protocol.NewProtocolGRPC(funcer)
	shortenerService := handler.NewShortenerService(APIURLServ, protGRPC)

	// Servers.
	server.RunGRPC(a.Cfg, a.Logg, shortenerService, interceptor)
	server.RunHTTP(a.Cfg, a.Logg, router, middleware)

	defer setupLayers.Close()
	defer infraLogg.Close()
	defer authLogg.Close()
	defer sub1.Close()
	defer sub2.Close()
	defer a.close()
	defer cancel()
}
