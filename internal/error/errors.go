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
	return fmt.Sprintf("[ERROR]:%s>%v|%s %d",
		p.Message, p.Err, p.File, p.Line)
}

func (p *ErrPlace) Unwrap() error {
	return p.Err
}

// ErrWrap
type ErrWrap struct {
	Message string
	Err     error
}

func NewErrWrap(mess string, err error) *ErrWrap {
	return &ErrWrap{
		Message: mess,
		Err:     err,
	}
}

func (w *ErrWrap) Error() string {
	return fmt.Sprintf("[ERROR]:%s>%v", w.Message, w.Err)
}

func (w *ErrWrap) Unwrap() error {
	return w.Err
}
