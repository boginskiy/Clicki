package error

import (
	"errors"
	"fmt"
	"runtime"
)

var (
	ErrUniqueData   = errors.New("attempt to overwrite unique data") // Ошибка при попытке перезаписать поле с уникальными данными
	ErrPingDataBase = errors.New("bad database ping")                // Ошибка подключения к БД
)

// ErrPlace
type ErrPlace struct {
	Message string
	File    string
	Line    int
	Err     error
}

func NewErrPlace(mess string, err error) *ErrPlace {
	// Получаем информацию о месте вызова
	_, file, line, _ := runtime.Caller(1)

	return &ErrPlace{
		Message: mess,
		File:    file,
		Line:    line,
		Err:     err,
	}
}

func (p *ErrPlace) Error() string {
	return fmt.Sprintf("%s at %s:%d > %v", p.Message, p.File, p.Line, p.Err)
}

func (p *ErrPlace) Unwrap() error {
	return p.Err
}
