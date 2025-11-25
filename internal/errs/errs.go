package errs

import (
	"errors"
	"fmt"
	"runtime"
)

// Errors.
var (
	// ErrUniqueData - if try not unic data.
	ErrUniqueData = errors.New("attempt to overwrite unique data")
	// ErrPingDataBase - if db not available.
	ErrPingDataBase = errors.New("bad database ping")
)

// ErrPlace - error with place when raise err.
type ErrPlace struct {
	Message string
	File    string
	Line    int
	Err     error
}

func NewErrPlace(mess string, err error) *ErrPlace {
	// Get info about place of call.
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

// ErrWrap - error with wrapp another err.
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
