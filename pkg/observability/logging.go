package observability

import (
	"encoding/json"
	"io"
	"log"
)

type StructuredLogger struct {
	l *log.Logger
}

func NewStructuredLogger(w io.Writer) StructuredLogger {
	return StructuredLogger{l: log.New(w, "", 0)}
}

func (s StructuredLogger) Info(msg string, fields map[string]any) {
	payload := map[string]any{"level": "info", "msg": msg, "fields": fields}
	buf, _ := json.Marshal(payload)
	s.l.Println(string(buf))
}
