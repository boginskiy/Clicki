package server

import (
	"net/http"

	c "github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/db"
	l "github.com/boginskiy/Clicki/internal/logger"
	m "github.com/boginskiy/Clicki/internal/middleware"
	r "github.com/boginskiy/Clicki/internal/router"
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
	db := db.NewDBStore()

	// writing log...
	infoLog.RaiseInfo(l.StartedServInfo,
		l.Fields{"port": kwargs.GetSrvAddr()},
	)

	// Start server
	err := http.ListenAndServe(kwargs.GetSrvAddr(),
		r.Router(kwargs, midWare, db, fatalLog))

	// writing log...
	if err != nil {
		fatalLog.RaiseFatal(l.StartedServFatal,
			l.Fields{"port": kwargs.GetSrvAddr()},
		)
	}
}
