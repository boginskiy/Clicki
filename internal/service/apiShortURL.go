package service

import (
	"net/http"

	"github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/db"
	"github.com/boginskiy/Clicki/internal/db2"
	l "github.com/boginskiy/Clicki/internal/logger"
	m "github.com/boginskiy/Clicki/internal/model"
	p "github.com/boginskiy/Clicki/internal/preparation"
	v "github.com/boginskiy/Clicki/internal/validation"
	"github.com/boginskiy/Clicki/pkg"
)

type APIShortURL struct {
	ExtraFuncer p.ExtraFuncer
	DB          db.Storage
	DB2         db2.DBConnecter
	Checker     v.Checker
	Logger      l.Logger
}

func NewAPIShortURL(db db.Storage, db2 db2.DBConnecter,
	logger l.Logger, checker v.Checker, extraFuncer p.ExtraFuncer) *APIShortURL {

	return &APIShortURL{
		ExtraFuncer: extraFuncer,
		Checker:     checker,
		Logger:      logger,
		DB:          db,
		DB2:         db2,
	}
}

func (s *APIShortURL) encryptionLongURL() (imitationPath string) {
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

func (s *APIShortURL) Create(req *http.Request, kwargs config.VarGetter) ([]byte, error) {
	// Deserialization Body
	baseLink := m.NewBaseLink()
	err := s.ExtraFuncer.Deserialization(req, baseLink)

	if err != nil {
		s.Logger.RaiseFatal(err, DeserializFatal, nil)
		return EmptyByteSlice, err
	}

	// Валидируем URL. Проверка регуляркой, что строка является доменом сайта
	if !s.Checker.CheckUpURL(baseLink.URL) || baseLink.URL == "" {
		s.Logger.RaiseInfo("APIShortURL.Create>CheckUpURL",
			l.Fields{"error": ErrDataNotValid.Error()})
		return EmptyByteSlice, ErrDataNotValid
	}

	imitationPath := s.encryptionLongURL()     // Генерируем ключ
	s.DB.PutValue(imitationPath, baseLink.URL) // Кладем в DB данные

	// Serialization Body
	extraLink := m.NewExtraLink(baseLink, kwargs.GetBaseURL()+"/"+imitationPath)
	result, err := s.ExtraFuncer.Serialization(extraLink)

	if err != nil {
		s.Logger.RaiseError(err, "APIShortURL.Create>NewExtraLink", nil)
		return EmptyByteSlice, err
	}
	return result, nil
}

func (s *APIShortURL) Read(req *http.Request) ([]byte, error) {
	return EmptyByteSlice, nil
}

func (s *APIShortURL) CheckPing(req *http.Request) ([]byte, error) {
	return EmptyByteSlice, nil
}
