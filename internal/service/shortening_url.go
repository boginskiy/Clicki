package service

import (
	"errors"
	"net/http"
	"strings"

	"github.com/boginskiy/Clicki/internal/db"
	p "github.com/boginskiy/Clicki/internal/preparation"
	v "github.com/boginskiy/Clicki/internal/validation"
	"github.com/boginskiy/Clicki/pkg"
)

const LONG = 8

type ShortenerURL interface {
	Execute(request *http.Request) error
	EncryptionLongURL() string
	GetImitationPath() string
	GetOriginURL() string
}

type ShorteningURL struct {
	imitationPath string
	originURL     string
	ExtraFuncer   p.ExtraFuncer
	DB            db.Storage
	Checker       v.Checker
}

func NewShorteningURL(db db.Storage, checker v.Checker, extraFuncer p.ExtraFuncer) *ShorteningURL {
	return &ShorteningURL{
		ExtraFuncer: extraFuncer,
		Checker:     checker,
		DB:          db,
	}
}

func (s *ShorteningURL) EncryptionLongURL() (imitationPath string) {
	for {
		// Вызов шифратора
		imitationPath = pkg.Scramble(LONG)
		// Проверка на уникальность
		if _, err := s.DB.GetValue(imitationPath); err != nil {
			break
		}
	}
	return imitationPath
}

// Getter. Слабое место для дальнейших расширений приложения
func (s *ShorteningURL) GetImitationPath() string {
	return s.imitationPath
}

// Getter. Слабое место для дальнейших расширений приложения
func (s *ShorteningURL) GetOriginURL() string {
	return s.originURL
}

func (s *ShorteningURL) executePostReq(req *http.Request) error {
	// Вынимаем тело запроса
	originURL, err := s.ExtraFuncer.TakeAllBodyFromReq(req)
	if err != nil {
		return err
	}

	// Валидируем URL. Проверка регуляркой, что строка является доменом сайта
	if !s.Checker.CheckUpURL(originURL) || originURL == "" {
		return errors.New("data not available or invalid")
	}

	// Генерируем ключ
	s.imitationPath = s.EncryptionLongURL()
	// Кладем в DB данные
	s.DB.PutValue(s.imitationPath, originURL)

	return nil
}

func (s *ShorteningURL) executeGetReq(req *http.Request) error {
	// Достаем параметр id                         \\
	tmpPath := strings.TrimLeft(req.URL.Path, "/") // Вариант для прохождения inittests
	// tmpPath := chi.URLParam(req, "id")         // Вариант из под коробки

	// Достаем origin URL
	tmpURL, err := s.DB.GetValue(tmpPath)
	if err != nil {
		return errors.New("data is not available")
	}
	s.originURL = tmpURL
	return nil
}

func (s *ShorteningURL) Execute(req *http.Request) error {
	switch req.Method {
	case "POST":
		return s.executePostReq(req)
	case "GET":
		return s.executeGetReq(req)
	default:
		return errors.New("request is not available")
	}
}
