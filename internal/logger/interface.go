package logger

type Logger interface {
	RaiseInfo(string, Fields)
	RaiseWarn(string, Fields)
	RaiseError(string, Fields)
	RaiseFatal(string, Fields)
	RaisePanic(string, Fields)
	CloseDesc()
}
