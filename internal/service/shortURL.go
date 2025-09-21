package service

import (
	"context"
	"net/http"
	"strings"

	conf "github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/logg"
	mod "github.com/boginskiy/Clicki/internal/model"
	prep "github.com/boginskiy/Clicki/internal/preparation"
	repo "github.com/boginskiy/Clicki/internal/repository"
	valid "github.com/boginskiy/Clicki/internal/validation"
	"github.com/boginskiy/Clicki/pkg"
)

type ShortURL struct {
	ExtraFuncer prep.ExtraFuncer
	Repo        repo.Repository
	Checker     valid.Checker
	Logger      logg.Logger
	Kwargs      conf.VarGetter
}

func NewShortURL(
	kwargs conf.VarGetter, logger logg.Logger, repo repo.Repository,
	checker valid.Checker, extraFuncer prep.ExtraFuncer) *ShortURL {

	return &ShortURL{
		ExtraFuncer: extraFuncer,
		Checker:     checker,
		Logger:      logger,
		Kwargs:      kwargs,
		Repo:        repo,
	}
}

func (s *ShortURL) encryptionLongURL() (correlID string) {
	for {
		correlID = pkg.Scramble(LONG)                   // Вызов шифратора
		if s.Repo.CheckUnic(context.TODO(), correlID) { // Проверка на уникальность
			break
		}
	}
	return correlID
}

func (s *ShortURL) GetHeader() string {
	return "text/plain"
}

func (s *ShortURL) Create(req *http.Request) ([]byte, error) {
	originURL, err := s.ExtraFuncer.TakeAllBodyFromReq(req) // Вынимаем тело запроса

	if err != nil {
		s.Logger.RaiseFatal(err, "ShortURL.Create>TakeAllBodyFromReq", nil)
		return EmptyByteSlice, err
	}

	// Валидируем URL. Проверка регуляркой, что строка является доменом сайта
	if !s.Checker.CheckUpURL(originURL) || originURL == "" {
		s.Logger.RaiseError(ErrDataNotValid, "ShortURL.Create>CheckUpURL", nil)
		return EmptyByteSlice, ErrDataNotValid
	}

	correlationID := s.encryptionLongURL()                           // Уникальный идентификатор
	shortURL := s.Kwargs.GetBaseURL() + "/" + correlationID          // Новый сокращенный URL
	preRecord := mod.NewURLTb(0, correlationID, originURL, shortURL) // Создаем запись
	record, err := s.Repo.Create(context.TODO(), preRecord)          // Кладем в DB данные

	if err != nil && record == nil {
		s.Logger.RaiseError(err, "ShortURL.Create>Repo.Create", nil)
		return EmptyByteSlice, err
	}

	//
	switch r := record.(type) {
	case *mod.URLTb:
		return []byte(r.ShortURL), err
	default:
		s.Logger.RaiseError(err, "ShortURL.Create>switch", nil)
		return EmptyByteSlice, err
	}
}

func (s *ShortURL) Read(req *http.Request) ([]byte, error) {
	correlationID := strings.TrimLeft(req.URL.Path, "/")      // Достаем параметр correlationID
	record, err := s.Repo.Read(context.TODO(), correlationID) // Достаем origin URL

	if err != nil {
		s.Logger.RaiseError(err, "ShortURL.Read>DB.Read", nil)
		return EmptyByteSlice, ErrDataNotValid
	}

	switch r := record.(type) {
	case *mod.URLTb:
		return []byte(r.OriginalURL), nil
	default:
		s.Logger.RaiseError(err, "ShortURL.Read>DB.Read>switch", nil)
		return EmptyByteSlice, ErrDataNotValid
	}
}

func (s *ShortURL) CheckPing(req *http.Request) ([]byte, error) {
	_, err := s.Repo.Ping(context.TODO())
	if err != nil {
		s.Logger.RaiseFatal(err, "ShortURL.CreaCheckPingte>Ping", nil)
		return EmptyByteSlice, err
	}
	return StoreDBIsSucces, nil
}

func (s *ShortURL) SetBatch(req *http.Request) ([]byte, error) {
	return StoreDBIsSucces, nil
}
