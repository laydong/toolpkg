package httpx

import (
	"bytes"
	"github.com/gin-gonic/gin"
)

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func (r responseBodyWriter) WriteString(s string) (n int, err error) {
	r.body.WriteString(s)
	return r.ResponseWriter.WriteString(s)
}
