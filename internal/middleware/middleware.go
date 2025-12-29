package middleware

import (
	"context"
	"errors"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/auth"
	"github.com/boginskiy/Clicki/internal/gzip"
	"github.com/boginskiy/Clicki/internal/logg"
)

// MvFunc - type func of HandlerFunc.
type MvFunc func(http.HandlerFunc) http.HandlerFunc

// ErrTrustedSubnet -.
var ErrTrustedSubnet error = errors.New(`{"error":"ip is not trusted"}`)

// Mdlwere - struct with function of Middleware.
type Mdlwere struct {
	Cfg  config.Config
	Auth auth.Auther
	Logg logg.Logger
}

func NewMdlwere(config config.Config, logger logg.Logger, auther auth.Auther) *Mdlwere {
	return &Mdlwere{Cfg: config, Logg: logger, Auth: auther}
}

func (m *Mdlwere) Conveyor(next http.HandlerFunc) http.HandlerFunc {
	// TODO: m.WithAudit .
	for _, middleware := range []MvFunc{m.WithAuth, m.WithLogg, m.WithGzip} {
		next = middleware(next)
	}
	return next
}

func (m *Mdlwere) WithLogg(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		uri := r.RequestURI
		method := r.Method

		// Extension standart ResponseWriter.
		extW := NewExResWriter(w)
		next(extW, r)

		duration := time.Since(start)

		m.Logg.RaiseInfo(
			logg.DataReqResInfo,
			map[string]any{
				"uri":      uri,
				"method":   method,
				"duration": duration,
				"status":   extW.status,
				"size":     extW.size,
			})
	}
}

func (m *Mdlwere) WithGzip(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpW := w

		// Checking Encoding and Type.
		acceptEncoding := r.Header.Get("Accept-Encoding")
		acceptContent := r.Header.Get("Content-Type")

		jsonGzip := strings.Contains(acceptContent, "application/json")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		htmlGzip := strings.Contains(acceptContent, "text/html")

		if supportsGzip && (jsonGzip || htmlGzip) {
			// Wrapp http.ResponseWriter with new gzip.
			compW := gzip.NewCompressWriter(w)
			tmpW = compW
			defer compW.Close()
		}

		// Check about user sent compressed data.
		contentEncoding := r.Header.Get("Content-Encoding")
		sendGzip := strings.Contains(contentEncoding, "gzip")

		if sendGzip {
			decompR, err := gzip.NewCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			// Change body of request on new.
			r.Body = decompR
		}
		// Transfer of control.
		next(tmpW, r)
	}
}

func (m *Mdlwere) WithAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, UserID, err := m.Auth.Authentication(r)

		// Errors with "user not found" and "create token".
		if err == auth.ErrUserNotFound || err == auth.ErrCreateToken {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Erorrs validation token.
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.SetCookie(w, cookie)
		ctx := context.WithValue(r.Context(), auth.CtxUserID, UserID)
		next(w, r.WithContext(ctx))
	}
}

func (m *Mdlwere) WithTrustedSubnet(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		_, cidr, _ := net.ParseCIDR(m.Cfg.GetTrustedSubnet()) // TrustedSubnet.
		userIP := net.ParseIP(r.Header.Get("X-Real-IP"))      // userIP.

		// Check of IP belonging to a subnet.
		if userIP != nil && cidr != nil && cidr.Contains(userIP) {
			next(w, r)
		} else {
			http.Error(w, ErrTrustedSubnet.Error(), http.StatusForbidden)
			return
		}
	}
}
