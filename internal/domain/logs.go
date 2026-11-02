package domain

import (
	"strings"
	"time"
)

type Logs struct {
	Date    string `json:"date"`
	Message string `json:"message"`
}

type LogCollector struct {
	logs []Logs
}

func (l *LogCollector) AddLog(message string) {
	if len(l.logs) >= 5000 {
		l.logs = l.logs[1:]
	}
	l.logs = append(l.logs, Logs{
		Date:    time.Now().Format(time.RFC3339),
		Message: message,
	})
}

func (l *LogCollector) AddListLogs(logs []Logs) {
	l.logs = append(l.logs, logs...)
}

func (l *LogCollector) AddLogWithTime(message string, date time.Time) {
	if len(l.logs) >= 5000 {
		l.logs = l.logs[1:]
	}
	l.logs = append(l.logs, Logs{
		Date:    time.Now().Format(time.RFC3339),
		Message: message,
	})
}

func (l *LogCollector) GetLogs() []Logs {
	return l.logs
}

func (l *LogCollector) Count() int {
	return len(l.logs)
}

func (l *LogCollector) Redact(secret string) {
	if secret == "" {
		return
	}
	for i := range l.logs {
		l.logs[i].Message = strings.ReplaceAll(l.logs[i].Message, secret, "[REDACTED]")
	}
}
