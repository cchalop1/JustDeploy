package domain

import "time"

type Logs struct {
	Date    string `json:"date"`
	Message string `json:"message"`
}

type LogCollector struct {
	logs []Logs
}

func (l *LogCollector) AddLog(message string) {
	l.logs = append(l.logs, Logs{
		Date:    time.Now().Format(time.RFC3339),
		Message: message,
	})
}

func (l *LogCollector) AddListLogs(logs []Logs) {
	l.logs = append(l.logs, logs...)
}

func (l *LogCollector) AddLogWithTime(message string, date time.Time) {
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
