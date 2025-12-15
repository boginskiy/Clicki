package server

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"log"
	"net/http"

	"github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/logg"
	mv "github.com/boginskiy/Clicki/internal/middleware"
	"github.com/boginskiy/Clicki/internal/router"
)

var SERT = "private.pem"
var PRIVATE = "cert.pem"
var LONG = 4096

type ServS struct {
	Cfg         config.Config
	Logg        logg.Logger
	PrivateKey  *rsa.PrivateKey
	CertName    string
	PrivateName string
}

func NewServS(config config.Config, logger logg.Logger) *ServS {
	tmp := &ServS{
		Cfg:  config,
		Logg: logger,
	}
	tmp.Settings(SERT, PRIVATE)
	return tmp
}

func (h *ServS) Settings(sertName, privateName string) {
	cert := NewX509C()                                    // создаём шаблон сертификата
	privateKey, err := rsa.GenerateKey(rand.Reader, LONG) // privateKey приватный RSA-ключ длиной 4096 бит
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
	h.PrivateName = SaveFilePem(privateName, privateKeyPEM.Bytes())
	h.CertName = SaveFilePem(sertName, certPEM.Bytes())
}

func (h *ServS) Run(router router.Router, mdlwere mv.Middleware) {
	h.Logg.RaiseFatal(
		http.ListenAndServeTLS(h.Cfg.GetSrvAddr(), h.CertName, h.PrivateName, router.Run(mdlwere)),
		"server has not started", nil)
}
