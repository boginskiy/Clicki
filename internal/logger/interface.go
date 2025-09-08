package logger

type Logger interface {
	RaiseInfo(string, Fields)
	RaiseWarn(string, Fields)
	RaiseError(error, string, Fields)
	RaiseFatal(error, string, Fields)
	RaisePanic(error, string, Fields)
	CloseDesc()
}
