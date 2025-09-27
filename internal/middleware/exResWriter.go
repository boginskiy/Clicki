package middleware

import "net/http"

// Extension standart function of http.ResponseWriter
type ExResWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func NewExResWriter(w http.ResponseWriter) *ExResWriter {
	return &ExResWriter{
		ResponseWriter: w,
	}
}

func (r *ExResWriter) Write(b []byte) (int, error) {
	// записываем ответ, используя оригинальный http.ResponseWriter
	size, err := r.ResponseWriter.Write(b)
	r.size += size // захватываем размер
	return size, err
}

func (r *ExResWriter) WriteHeader(statusCode int) {
	// записываем код статуса, используя оригинальный http.ResponseWriter
	r.ResponseWriter.WriteHeader(statusCode)
	r.status = statusCode // захватываем код статуса
}
