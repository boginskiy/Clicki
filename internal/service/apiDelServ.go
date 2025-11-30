package service

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	conf "github.com/boginskiy/Clicki/cmd/config"
	"github.com/boginskiy/Clicki/internal/auth"
	"github.com/boginskiy/Clicki/internal/logg"
	repo "github.com/boginskiy/Clicki/internal/repository"
	"github.com/boginskiy/Clicki/internal/repository/utils"
)

// APIDelServ - struct about delete message.
type APIDelServ struct {
	Cfg         conf.Config
	Logg        logg.Logger
	Repo        repo.Repository
	delMessChan chan utils.DelMessage
}

func NewAPIDelServ(
	ctx context.Context,
	config conf.Config,
	logger logg.Logger,
	repository repo.Repository) *APIDelServ {

	item := &APIDelServ{
		Cfg:         config,
		Logg:        logger,
		Repo:        repository,
		delMessChan: make(chan utils.DelMessage, 8),
	}
	// Запуск фонового удаления данных.
	go item.StepByStepDelMessages(ctx)
	return item
}

// DeleteSet - Producer .
func (d *APIDelServ) DeleteSet(req *http.Request) ([]byte, error) {
	// Принимаем список идентификаторов URLs.
	dataByte, err := io.ReadAll(req.Body)
	if err != nil {
		return EmptyByteSlice, err
	}

	// Подготовка delMessage.
	userID := d.takeUserIDFromCtx(req)
	delMessage := utils.NewDelMessage(int64(userID))
	err = json.Unmarshal(dataByte, &delMessage.ListCorrelID)

	if err != nil {
		return EmptyByteSlice, err
	}

	// Отправка сообщения в канал.
	d.delMessChan <- *delMessage

	return EmptyByteSlice, nil
}

// StepByStepDelMessages - Concumer .
func (d *APIDelServ) StepByStepDelMessages(ctx context.Context) {
	// Каждые N-секунд перевод удаляемых данных в "Soft Delete".
	Nsec := time.Duration(d.Cfg.GetSoftDeleteTime())
	ticker := time.NewTicker(Nsec * time.Second)

	// Каждые N-секунд перевод удаляемых данных "Hard Delete".
	Nsec = time.Duration(d.Cfg.GetHardDeleteTime())
	ticker2 := time.NewTicker(Nsec * time.Second)

	var delMessages []utils.DelMessage
	var deletedSoft bool

	for {
		select {

		// Завершение работы горутины при отключении сервиса.
		case <-ctx.Done():
			d.sendSoftDeletion(delMessages, &deletedSoft)
			return

		// Добавление данных на удаление.
		case msg := <-d.delMessChan:
			delMessages = append(delMessages, msg)

		// Обращаемся к БД для маркировки удаляемых данных.
		case <-ticker.C:
			delMessages = d.sendSoftDeletion(delMessages, &deletedSoft)

		// Физическое удаление помеченных данных.
		case <-ticker2.C:
			deletedSoft = d.sendHardDeletion(&deletedSoft)
		}
	}
}

func (d *APIDelServ) takeUserIDFromCtx(req *http.Request) int {
	UserID, ok := req.Context().Value(auth.CtxUserID).(int)
	if !ok || UserID <= 0 {
		d.Logg.RaiseError(ErrUserIDNotValid, "APIDelServ.TakeUserIDFromCtx>CtxUserID", nil)
	}
	return UserID
}

func (d *APIDelServ) sendHardDeletion(isDel *bool) bool {
	if *isDel {
		err := d.Repo.DeleteRecords(context.TODO())
		if err != nil {
			d.Logg.RaiseError(err, "DelMess>StepByStepDelMessages>sendHardDeletion", nil)
		}
	}
	return false
}

func (d *APIDelServ) sendSoftDeletion(data []utils.DelMessage, isDel *bool) []utils.DelMessage {
	if 0 < len(data) {
		err := d.Repo.MarkRecords(context.TODO(), data)

		if err != nil {
			d.Logg.RaiseError(err, "DelMess>StepByStepDelMessages>sendSoftDeletion", nil)
		} else {
			// Обнуляем очередь сообщений.
			*isDel = true
			return data[:0]
		}
	}
	return data
}
