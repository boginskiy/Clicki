package service

import (
	"context"
	"net/http"
	"strings"

	mod "github.com/boginskiy/Clicki/internal/model"
	prep "github.com/boginskiy/Clicki/internal/preparation"
	repo "github.com/boginskiy/Clicki/internal/repository"
	valid "github.com/boginskiy/Clicki/internal/validation"
)

type ShortURL struct {
	ExtraFuncer prep.ExtraFuncer
	Repo        repo.Repository
	Checker     valid.Checker
	Core        *CoreService
}

func NewShortURL(
	core *CoreService,
	repo repo.Repository,
	checker valid.Checker,
	extraFuncer prep.ExtraFuncer) *ShortURL {

	return &ShortURL{
		ExtraFuncer: extraFuncer,
		Core:        core,
		Checker:     checker,
		Repo:        repo,
	}
}

func (s *ShortURL) CreateSetURL(req *http.Request) ([]byte, error) {
	return StoreDBIsSucces, nil
}

func (s *ShortURL) ReadSetUserURL(req *http.Request) ([]byte, error) {
	return StoreDBIsSucces, nil
}

func (s *ShortURL) GetHeader() string {
	return "text/plain"
}

func (s *ShortURL) CreateURL(req *http.Request) ([]byte, error) {
	originURL, err := s.ExtraFuncer.TakeAllBodyFromReq(req) // Вынимаем тело запроса

	if err != nil {
		s.Core.Logg.RaiseFatal(err, "ShortURL.CreateURL>TakeAllBodyFromReq", nil)
		return EmptyByteSlice, err
	}

	// Валидируем URL. Проверка регуляркой, что строка является доменом сайта
	if !s.Checker.CheckUpURL(originURL) || originURL == "" {
		s.Core.Logg.RaiseError(ErrDataNotValid, "ShortURL.CreateURL>CheckUpURL", nil)
		return EmptyByteSlice, ErrDataNotValid
	}

	userID := s.Core.TakeUserIDFromCtx(req) // Тащим идентификатор пользователя

	correlationID := s.Core.EncrypOriginURL()                    // Уникальный идентификатор
	shortURL := s.Core.Kwargs.GetBaseURL() + "/" + correlationID // Новый сокращенный URL

	preRecord := mod.NewURLTb(0, correlationID, originURL, shortURL, userID) // Создаем запись
	record, err := s.Repo.CreateRecord(context.TODO(), preRecord)            // Кладем в DB данные

	if err != nil && record == nil {
		s.Core.Logg.RaiseError(err, "ShortURL.CreateURL>Repo.Create", nil)
		return EmptyByteSlice, err
	}

	// Definition type
	switch r := record.(type) {
	case *mod.URLTb:
		return []byte(r.ShortURL), err
	default:
		s.Core.Logg.RaiseError(err, "ShortURL.Create>switch", nil)
		return EmptyByteSlice, err
	}
}

func (s *ShortURL) ReadURL(req *http.Request) ([]byte, error) {
	correlationID := strings.TrimLeft(req.URL.Path, "/")            // Достаем параметр correlationID
	record, err := s.Repo.ReadRecord(context.TODO(), correlationID) // Достаем origin URL

	if err != nil {
		s.Core.Logg.RaiseError(err, "ShortURL.Read>DB.Read", nil)
		return EmptyByteSlice, ErrDataNotValid
	}

	// Definition type
	switch r := record.(type) {
	case *mod.URLTb:

		// Если фла==true, запись стоит в очереди на удаление
		if r.DeletedFlag {
			return EmptyByteSlice, ErrReadRecord
		}
		return []byte(r.OriginalURL), nil

	default:
		s.Core.Logg.RaiseError(err, "ShortURL.Read>DB.Read>switch", nil)
		return EmptyByteSlice, ErrDataNotValid
	}
}

func (s *ShortURL) CheckDB(req *http.Request) ([]byte, error) {
	_, err := s.Repo.PingDB(context.TODO())
	if err != nil {
		s.Core.Logg.RaiseFatal(err, "ShortURL.CreaCheckPingte>Ping", nil)
		return EmptyByteSlice, err
	}
	return StoreDBIsSucces, nil
}
