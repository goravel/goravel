package gin

import "github.com/gin-gonic/gin"

type StreamWriter struct {
	instance *gin.Context
}

func NewStreamWriter(instance *gin.Context) *StreamWriter {
	return &StreamWriter{instance}
}

func (w *StreamWriter) Flush() error {
	w.instance.Writer.Flush()
	return nil
}

func (w *StreamWriter) Write(data []byte) (int, error) {
	return w.instance.Writer.Write(data)
}

func (w *StreamWriter) WriteString(s string) (int, error) {
	return w.instance.Writer.WriteString(s)
}
