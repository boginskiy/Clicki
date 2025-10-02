package service

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	conf "github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/logg"
	mod "github.com/boginskiy/Clicki/internal/model"
	prep "github.com/boginskiy/Clicki/internal/preparation"
	rep "github.com/boginskiy/Clicki/internal/repository"
	valid "github.com/boginskiy/Clicki/internal/validation"
)

type APIShortURL struct {
	Repo        rep.Repository
	ExtraFuncer prep.ExtraFuncer
	Kwargs      conf.VarGetter
	Checker     valid.Checker
	Logger      logg.Logger
	Core        CoreServicer
	delMessChan chan rep.DelMessage
}

func NewAPIShortURL(
	kwargs conf.VarGetter, logger logg.Logger, repo rep.Repository,
	checker valid.Checker, extraFuncer prep.ExtraFuncer) *APIShortURL {

	instance := &APIShortURL{
		Core:        NewCoreService(kwargs, logger, repo),
		delMessChan: make(chan rep.DelMessage, 8),
		ExtraFuncer: extraFuncer,
		Checker:     checker,
		Logger:      logger,
		Kwargs:      kwargs,
		Repo:        repo,
	}

	// Запуск фонового удаления данных
	go instance.destroyMessages()

	return instance
}

func (s *APIShortURL) ReadURL(req *http.Request) ([]byte, error) {
	return EmptyByteSlice, nil
}

func (s *APIShortURL) CheckDB(req *http.Request) ([]byte, error) {
	return EmptyByteSlice, nil
}

func (s *APIShortURL) GetHeader() string {
	return "application/json"
}

func (s *APIShortURL) CreateURL(req *http.Request) ([]byte, error) {
	// Deserialization Body
	bodyJSON := mod.NewURLJson()
	err := s.ExtraFuncer.Deserialization(req, bodyJSON)

	if err != nil {
		s.Logger.RaiseFatal(err, DeserializFatal, nil)
		return EmptyByteSlice, err
	}

	// Валидируем URL. Проверка регуляркой, что строка является доменом сайта
	if !s.Checker.CheckUpURL(bodyJSON.URL) || bodyJSON.URL == "" {
		s.Logger.RaiseInfo("APIShortURL.Create>CheckUpURL",
			logg.Fields{"error": ErrDataNotValid.Error()})
		return EmptyByteSlice, ErrDataNotValid
	}

	userID := s.Core.takeUserIDFromCtx(req)                 // Тащим идентификатор пользователя
	correlationID := s.Core.encrypOriginURL()               // Уникальный идентификатор
	shortURL := s.Kwargs.GetBaseURL() + "/" + correlationID // Создаем новый сокращенный URL

	preRecord := mod.NewURLTb(0, correlationID, bodyJSON.URL, shortURL, userID) // Создаем черновую запись
	record, err := s.Repo.CreateRecord(context.TODO(), preRecord)               // Кладем в DB данные

	if err != nil && record == nil {
		s.Logger.RaiseError(err, "APIShortURL.Create>Repo.Create", nil)
		return EmptyByteSlice, err
	}

	// Definition type
	var resJSON *mod.ResultJSON
	switch r := record.(type) {
	case *mod.URLTb:
		resJSON = mod.NewResultJSON(bodyJSON, r.ShortURL)
	case string:
		resJSON = mod.NewResultJSON(bodyJSON, r)
	default:
		s.Logger.RaiseError(err, "APIShortURL.Create>switch", nil)
		return EmptyByteSlice, err
	}

	// Serialization
	result, err2 := s.ExtraFuncer.Serialization(resJSON)

	if err2 != nil {
		s.Logger.RaiseError(err2, "APIShortURL.Create>NewResultJSON", nil)
		return EmptyByteSlice, err2
	}
	return result, err
}

func (s *APIShortURL) CreateSetURL(req *http.Request) ([]byte, error) {
	// Создаем декодер
	decoder := json.NewDecoder(req.Body)

	// Проверка, что пришло то, что надо
	token, _ := decoder.Token()
	if token != json.Delim('[') {
		s.Logger.RaiseFatal(ErrDataNotValid, "ShortURL>SetBatch>Token",
			logg.Fields{"fatal": ErrDataNotValid.Error()})
		return EmptyByteSlice, ErrDataNotValid
	}

	// Достаем идентификатор пользователя
	userID := s.Core.takeUserIDFromCtx(req)

	// Разбор тела запроса
	respURLSet := make([]mod.ResURLSet, 0, 10)

	for decoder.More() {
		var rURL mod.ReqURLSet
		err := decoder.Decode(&rURL)

		if err != nil {
			s.Logger.RaiseFatal(err, "ShortURL>SetBatch>Decode", nil)
			return EmptyByteSlice, err
		}

		shortURL := s.Kwargs.GetBaseURL() + "/" + rURL.CorrelationID
		// Сбор множества URL
		respURLSet = append(respURLSet, mod.NewResURLSet(
			rURL.CorrelationID, rURL.OriginalURL, shortURL, userID))
	}

	// Сохранение в БД
	err := s.Repo.CreateRecords(context.TODO(), respURLSet)
	if err != nil {
		s.Logger.RaiseError(err, "APIShortURL>SetBatch>CreateSet", nil)
		return EmptyByteSlice, err
	}

	// Сериализуем
	result, err := json.Marshal(respURLSet)
	s.Logger.RaiseFatal(err, "ShortURL>SetBatch>Marshal", nil)
	return result, nil
}

func (s *APIShortURL) ReadSetUserURL(req *http.Request) ([]byte, error) {
	// Достаем идентификатор пользователя
	userID := s.Core.takeUserIDFromCtx(req)

	dataSet, err := s.Repo.ReadRecords(context.TODO(), userID)
	if err != nil {
		s.Logger.RaiseError(err, "APIShortURL.ReadSetUserURL>ReadSet", nil)
		return EmptyByteSlice, err
	}

	// Однозначно определяем наличие записей у пользователя
	records, ok := dataSet.([]mod.ResUserURLSet)
	if !ok {
		s.Logger.RaiseError(ErrDataNotValid, "APIShortURL.ReadSetUserURL>Type?", nil)
		return EmptyByteSlice, ErrDataNotValid
	}
	if len(records) == 0 {
		return EmptyByteSlice, nil
	}

	// Serialization
	result, err := s.ExtraFuncer.Serialization(records)
	if err != nil {
		s.Logger.RaiseError(err, "APIShortURL.ReadSetUserURL>Serialization", nil)
		return EmptyByteSlice, err
	}
	return result, err
}

// Producer
func (s *APIShortURL) DeleteSetUserURL(req *http.Request) ([]byte, error) {
	// Принимаем список идентификаторов URLs
	dataByte, err := io.ReadAll(req.Body)
	if err != nil {
		return EmptyByteSlice, err
	}

	// Подготовка delMessage
	userID := s.Core.takeUserIDFromCtx(req)
	delMessage := rep.NewDelMessage(userID)
	err = json.Unmarshal(dataByte, &delMessage.ListCorrelID)

	if err != nil {
		return EmptyByteSlice, err
	}

	// Отправка сообщения в канал
	s.delMessChan <- *delMessage

	return EmptyByteSlice, nil
}

// Concumer
func (s *APIShortURL) destroyMessages() {
	// Каждые N-секунд перевод удаляемых данных в "Soft Delete"
	Nsec := time.Duration(s.Kwargs.GetSoftDeleteTime())
	ticker := time.NewTicker(Nsec * time.Second)

	// Каждые N-секунд перевод удаляемых данных "Hard Delete"
	Nsec = time.Duration(s.Kwargs.GetHardDeleteTime())
	ticker2 := time.NewTicker(Nsec * time.Second)

	var delMessages []rep.DelMessage
	var deletedSoft bool

	for {
		select {

		// Добавление данных на удаление
		case msg := <-s.delMessChan:
			delMessages = append(delMessages, msg)

		// Обращаемся к БД для маркировки удаляемых данных
		case <-ticker.C:
			if len(delMessages) == 0 {
				continue
			}
			err := s.Repo.MarkerRecords(context.TODO(), delMessages...)
			if err != nil {
				s.Logger.RaiseError(err, "APIShortURL>destroyMessages>MarkerRecords", nil)
				continue
			}

			delMessages = delMessages[:0]
			deletedSoft = true

		// Физическое удаление помеченных данных
		case <-ticker2.C:
			if deletedSoft {
				err := s.Repo.DeleteRecords(context.TODO())
				if err != nil {
					s.Logger.RaiseError(err, "APIShortURL>destroyMessages>MarkerRecords", nil)
				}
				deletedSoft = false
			}
		}
	}
}
