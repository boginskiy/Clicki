package middleware

import "net/http"

// Extension Status and Size standart of http.ResponseWriter
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

// Extension UserID standart of http.ResponseWriter
type ExResWriter2 struct {
	http.ResponseWriter
	UserID int
}

func NewExResWriter2(w http.ResponseWriter, id int) *ExResWriter2 {
	return &ExResWriter2{
		ResponseWriter: w,
		UserID:         id,
	}
}

func (r *ExResWriter2) Write(b []byte) (int, error) {
	return r.ResponseWriter.Write(b)
}

func (r *ExResWriter2) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
}
