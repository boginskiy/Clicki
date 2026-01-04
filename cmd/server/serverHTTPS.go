package server

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/logg"
)

type ServS struct {
	Cfg  config.Config
	Logg logg.Logger
	S    *http.Server
	done chan struct{}

	PrivateKey  *rsa.PrivateKey
	CertName    string
	PrivateName string
}

func NewServS(config config.Config, logger logg.Logger, handler http.Handler) *ServS {
	tmpServ := &ServS{
		Cfg:  config,
		Logg: logger,
		S: &http.Server{
			Addr:    config.GetSrvAddr(),
			Handler: handler,
		},
		done: make(chan struct{}),
	}

	tmpServ.WorkingWithShutdown()
	tmpServ.Settings(SERT, PRIVATE)
	return tmpServ
}

func (s *ServS) Settings(sertName, privateName string) {
	cert := NewX509C()                                       // создаём шаблон сертификата
	privateKey, err := rsa.GenerateKey(rand.Reader, LONGKEY) // privateKey приватный RSA-ключ длиной 4096 бит
	if err != nil {
		log.Fatal(err)
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &privateKey.PublicKey, privateKey)
	if err != nil {
		log.Fatal(err)
	}

	certPEM := NewEncodeCertPEM(certBytes)              // кодируем сертификат в формате PEM
	privateKeyPEM := NewEncodePrivateKeyPEM(privateKey) // кодируем ключ в формате PEM

	// Сохраняем сертификат и приватный ключ в файлы ~/cert.pem и ~/private.pem
	s.PrivateName = SaveFilePem(privateName, privateKeyPEM.Bytes())
	s.CertName = SaveFilePem(sertName, certPEM.Bytes())
}

func (s *ServS) WorkingWithShutdown() {
	//  Registration interruption.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-ctx.Done() // Recive signal
		shdownCtx, cancel := context.WithTimeout(context.Background(), SHDWTIME*time.Second)
		defer cancel()

		if err := s.S.Shutdown(shdownCtx); err != nil {
			s.Logg.RaiseFatal(err, "https server has bad Shutdown:", nil)
		}
		close(s.done)
		defer stop()
	}()
}

func (s *ServS) Run() {
	// Start server.
	if err := s.S.ListenAndServeTLS(s.CertName, s.PrivateName); err != http.ErrServerClosed {
		s.Logg.RaiseFatal(err, "https server has not started", nil)
	}

	// Waiting the end of Shutdown.
	<-s.done
	fmt.Fprint(os.Stdout, "\nServer has been successfully stopped\n")
}
