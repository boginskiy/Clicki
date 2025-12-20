package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/logg"
	"golang.org/x/crypto/acme/autocert"
)

var SHDWTIME = time.Duration(10)
var DIRCACHE = "cache-dir"
var HOST = "mysite.ru"

type Server interface {
	Run()
}

type Serv struct {
	Cfg  config.Config
	Logg logg.Logger
	S    *http.Server
	done chan struct{}
}

func NewServ(config config.Config, logger logg.Logger, handler http.Handler) *Serv {
	tmpServ := &Serv{
		Cfg:  config,
		Logg: logger,
		S: &http.Server{
			Addr:    config.GetSrvAddr(),
			Handler: handler,
		},
		done: make(chan struct{}),
	}

	tmpServ.WorkingWithShutdown()

	return tmpServ
}

func (s *Serv) StartWithAutocert(nameSite string, handler http.Handler) {
	s.Logg.RaiseFatal(http.Serve(autocert.NewListener(nameSite), handler),
		"server has not started", nil)
}

// SettingsManager конструируем менеджер TLS-сертификатов.
func (s *Serv) SettingsManager(dirCache string, hosts ...string) *autocert.Manager {
	return &autocert.Manager{
		Cache:      autocert.DirCache(dirCache),      // директория для хранения сертификатов.
		Prompt:     autocert.AcceptTOS,               // функция, принимающая Terms of Service издателя сертификатов.
		HostPolicy: autocert.HostWhitelist(hosts...), // перечень доменов, для которых будут поддерживаться сертификаты.
	}
}

// SettingsServerTLS конструируем сервер с поддержкой TLS
func (s *Serv) SettingsServerTLS(manager *autocert.Manager, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:    s.Cfg.GetSrvAddr(), // ":443"
		Handler: handler,
		// для TLS-конфигурации используем менеджер сертификатов
		TLSConfig: manager.TLSConfig(),
	}
}

// StartWithCustomAutocert. Cервер запускать с правами администратора и на хосте, имеющем «белый» IP-адрес.
func (s *Serv) StartWithCustomAutocert(handler http.Handler) {
	autocertManager := s.SettingsManager(DIRCACHE, HOST)
	server := s.SettingsServerTLS(autocertManager, handler)

	s.Logg.RaiseFatal(
		server.ListenAndServeTLS("", ""),
		"server has not started", nil)
}

func (s *Serv) WorkingWithShutdown() {
	//  Registration interruption.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-ctx.Done() // Recive signal
		shdownCtx, cancel := context.WithTimeout(context.Background(), SHDWTIME*time.Second)
		defer cancel()

		if err := s.S.Shutdown(shdownCtx); err != nil {
			s.Logg.RaiseFatal(err, "http server has bad Shutdown:", nil)
		}
		close(s.done)
		defer stop()
	}()
}

func (s *Serv) Run() {
	// if need turn on HTTPS.
	// if s.Cfg.GetEnableHTTPS() == "1" {
	// 	s.StartWithAutocert(HOST, s.S.Handler)
	// 	s.StartWithCustomAutocert(s.S.Handler)
	// } else {}

	// Start server.
	if err := s.S.ListenAndServe(); err != http.ErrServerClosed {
		s.Logg.RaiseFatal(err, "http server has not started", nil)
	}

	// Waiting the end of Shutdown.
	<-s.done
	fmt.Fprint(os.Stdout, "\nServer has been successfully stopped\n")
}
