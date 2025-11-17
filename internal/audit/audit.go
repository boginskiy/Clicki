package audit

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"

	conf "github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/auther"
	"github.com/boginskiy/Clicki/internal/logg"
)

type Audit struct {
	Publisher Publisher
	Kwargs    conf.VarGetter
	Logger    logg.Logger
	Rex       *regexp.Regexp
}

func NewAudit(kwargs conf.VarGetter, logger logg.Logger, publisher Publisher) *Audit {
	// Компилируем регулярное выражение. Далее будет проверка URL
	rex, err := regexp.Compile(pattern)
	if err != nil {
		logger.RaiseFatal(err, "Audit>NewAudit>Compile", nil)
	}

	return &Audit{
		Publisher: publisher,
		Kwargs:    kwargs,
		Logger:    logger,
		Rex:       rex,
	}
}

func (a *Audit) takeUserID(req *http.Request) int {
	UserID, ok := req.Context().Value(auther.CtxUserID).(int)
	if !ok {
		return 0
	}
	return UserID
}

func (a *Audit) readJSONBody(req *http.Request, field string) (string, error) {
	tmpResult := map[string]any{}

	bodyByte, err := io.ReadAll(req.Body)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(bodyByte, &tmpResult)
	if err != nil {
		return "", err
	}

	if v, ok := tmpResult[field]; ok {
		if url, ok := v.(string); ok {
			return url, nil
		}
	}
	return "", ErrReadJSONBody
}

func (a *Audit) TakeOriginURL(req *http.Request) (string, error) {
	switch req.Header.Get("Content-Type") {
	case "text/plain":
		bodyByte, err := io.ReadAll(req.Body)
		return string(bodyByte), err

	case "application/json":
		return a.readJSONBody(req, "url")

	default:
		return "unknown Content-Type", nil
	}
}

func (a *Audit) isItRightURL(req *http.Request) bool {
	if req.URL.String() == "/" && req.Method == "POST" ||
		req.URL.String() == "/api/shorten" && req.Method == "POST" ||
		a.Rex.MatchString(req.URL.String()) && req.Method == "GET" {
		return true
	}
	return false
}

func (a *Audit) NeedAudit(r *http.Request) bool {
	// 1. Проверка входящего URL со списком, подлежащим аудированию
	// 2. Проверка наличия подписчиков. Если их нет, аудирование не работает
	return a.isItRightURL(r) && a.Publisher.CheckSubscribers()
}

func (a *Audit) NoticeCreateLink(req *http.Request) {
	if req.Method != "POST" {
		return
	}
	// Берем оригинальный URL
	originURL, err := a.TakeOriginURL(req)
	if err != nil {
		a.Logger.RaiseError(err, "Audit>CreateLink>TakeOriginURL", nil)
		return
	}
	// Собираем событие аудита
	event := NewEvent("shorten", a.takeUserID(req), originURL)

	// Отправка события подписчикам
	a.Publisher.Send(event)
}

func (a *Audit) NoticeFollowLink(req *http.Request) {
}
