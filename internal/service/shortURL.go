package service

import (
	"net/http"
	"strings"

	"github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/db"
	"github.com/boginskiy/Clicki/internal/db2"
	l "github.com/boginskiy/Clicki/internal/logger"
	p "github.com/boginskiy/Clicki/internal/preparation"
	v "github.com/boginskiy/Clicki/internal/validation"
	"github.com/boginskiy/Clicki/pkg"
)

type ShortURL struct {
	ExtraFuncer p.ExtraFuncer
	DB          db.Storage
	DB2         db2.DBConnecter
	Checker     v.Checker
	Log         l.Logger
}

func NewShortURL(db db.Storage, db2 db2.DBConnecter,
	log l.Logger, checker v.Checker, extraFuncer p.ExtraFuncer) *ShortURL {

	return &ShortURL{
		ExtraFuncer: extraFuncer,
		Checker:     checker,
		Log:         log,
		DB:          db,
		DB2:         db2,
	}
}

func (s *ShortURL) encryptionLongURL() (imitationPath string) {
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

func (s *ShortURL) Create(req *http.Request, kwargs config.VarGetter) ([]byte, error) {
	// Вынимаем тело запроса
	originURL, err := s.ExtraFuncer.TakeAllBodyFromReq(req)

	if err != nil {
		s.Log.RaiseFatal("ShortURL.Create>TakeAllBodyFromReq",
			l.Fields{"error": err.Error()})
		return EmptyByteSlice, err
	}

	// Валидируем URL. Проверка регуляркой, что строка является доменом сайта
	if !s.Checker.CheckUpURL(originURL) || originURL == "" {
		s.Log.RaiseError("ShortURL.Create>CheckUpURL",
			l.Fields{"error": ErrDataNotValid.Error()})
		return EmptyByteSlice, ErrDataNotValid
	}

	// Генерируем ключ
	imitationPath := s.encryptionLongURL()
	// Кладем в DB данные
	s.DB.PutValue(imitationPath, originURL)

	return []byte(imitationPath), nil
}

func (s *ShortURL) Read(req *http.Request) ([]byte, error) {
	// Достаем параметр id                         \\
	tmpPath := strings.TrimLeft(req.URL.Path, "/") // Вариант для прохождения inittests
	// tmpPath := chi.URLParam(req, "id")         // Вариант из под коробки

	// Достаем origin URL
	tmpURL, err := s.DB.GetValue(tmpPath)

	if err != nil {
		s.Log.RaiseError("ShortURL.Read>GetValue",
			l.Fields{"error": ErrDataNotValid.Error()})
		return EmptyByteSlice, ErrDataNotValid
	}

	return []byte(tmpURL), nil
}

// CheckPing - check of connection DB
func (s *ShortURL) CheckPing(req *http.Request) ([]byte, error) {
	err := s.DB2.GetDB().Ping()
	if err != nil {
		return EmptyByteSlice, err
	}
	return ConnDBIsSucces, nil
}
