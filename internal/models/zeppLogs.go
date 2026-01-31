package models

import "time"

type logLevel uint8

const (
	Error logLevel = iota
	Warning
	Info
	Debug
)

type ZeppLogs struct {
	LogLevel    logLevel  `json:"log_level,omitempty"`
	Description string    `json:"description,omitempty"`
	TimeStamp   time.Time `json:"timestamp"`
}

func (ll logLevel) String() string {
	switch ll {
	case Error:
		return "ERROR"
	case Warning:
		return "WARNING"
	case Info:
		return "INFO"
	case Debug:
		return "DEBUG"
	default:
		return "UNKNOWN"
	}
}
