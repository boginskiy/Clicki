package middleware

import "net/http"

// Extension standart function of http.ResponseWriter
type ExtRespWrtr struct {
	http.ResponseWriter
	status int
	size   int
}

func NewExtRespWrtr(w http.ResponseWriter) *ExtRespWrtr {
	return &ExtRespWrtr{
		ResponseWriter: w,
	}
}

func (r *ExtRespWrtr) Write(b []byte) (int, error) {
	// записываем ответ, используя оригинальный http.ResponseWriter
	size, err := r.ResponseWriter.Write(b)
	r.size += size // захватываем размер
	return size, err
}

func (r *ExtRespWrtr) WriteHeader(statusCode int) {
	// записываем код статуса, используя оригинальный http.ResponseWriter
	r.ResponseWriter.WriteHeader(statusCode)
	r.status = statusCode // захватываем код статуса
}
