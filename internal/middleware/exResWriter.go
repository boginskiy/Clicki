package middleware

import "net/http"

// Extension standart function of http.ResponseWriter.
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
	// Write response with original http.ResponseWriter.
	size, err := r.ResponseWriter.Write(b)
	// Take size.
	r.size += size
	return size, err
}

func (r *ExResWriter) WriteHeader(statusCode int) {
	// Write code status with original http.ResponseWriter.
	r.ResponseWriter.WriteHeader(statusCode)
	// Take code status.
	r.status = statusCode
}
