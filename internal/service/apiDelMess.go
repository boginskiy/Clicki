package service

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	repo "github.com/boginskiy/Clicki/internal/repository"
)

// DelMess - struct about delete message.
type DelMess struct {
	Core        *CoreService
	Repo        repo.Repository
	delMessChan chan repo.DelMessage
}

func NewDelMess(ctx context.Context, core *CoreService, repository repo.Repository) *DelMess {
	item := &DelMess{
		delMessChan: make(chan repo.DelMessage, 8),
		Core:        core,
		Repo:        repository,
	}
	// Запуск фонового удаления данных.
	go item.StepByStepDelMessages(ctx)
	return item
}

// DeleteSetUserURL - Producer .
func (d *DelMess) DeleteSetUserURL(req *http.Request) ([]byte, error) {
	// Принимаем список идентификаторов URLs.
	dataByte, err := io.ReadAll(req.Body)
	if err != nil {
		return EmptyByteSlice, err
	}

	// Подготовка delMessage.
	userID := d.Core.TakeUserIDFromCtx(req)
	delMessage := repo.NewDelMessage(int64(userID))
	err = json.Unmarshal(dataByte, &delMessage.ListCorrelID)

	if err != nil {
		return EmptyByteSlice, err
	}

	// Отправка сообщения в канал.
	d.delMessChan <- *delMessage

	return EmptyByteSlice, nil
}

func (d *DelMess) sendSoftDeletion(data []repo.DelMessage, isDel *bool) []repo.DelMessage {
	if 0 < len(data) {
		err := d.Repo.MarkerRecords(context.TODO(), data...)

		if err != nil {
			d.Core.Logg.RaiseError(err, "DelMess>StepByStepDelMessages>sendSoftDeletion", nil)
		} else {
			// Обнуляем очередь сообщений.
			*isDel = true
			return data[:0]
		}
	}
	return data
}

func (d *DelMess) sendHardDeletion(isDel *bool) bool {
	if *isDel {
		err := d.Repo.DeleteRecords(context.TODO())
		if err != nil {
			d.Core.Logg.RaiseError(err, "DelMess>StepByStepDelMessages>sendHardDeletion", nil)
		}
	}
	return false
}

// StepByStepDelMessages - Concumer .
func (d *DelMess) StepByStepDelMessages(ctx context.Context) {
	// Каждые N-секунд перевод удаляемых данных в "Soft Delete".
	Nsec := time.Duration(d.Core.Cfg.GetSoftDeleteTime())
	ticker := time.NewTicker(Nsec * time.Second)

	// Каждые N-секунд перевод удаляемых данных "Hard Delete".
	Nsec = time.Duration(d.Core.Cfg.GetHardDeleteTime())
	ticker2 := time.NewTicker(Nsec * time.Second)

	var delMessages []repo.DelMessage
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
