package log

import (
	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func InitZeroLog() {
	// zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	// runLogFile, _ := os.OpenFile(
	// 	"log/app.log",
	// 	os.O_APPEND|os.O_CREATE|os.O_WRONLY,
	// 	0664,
	// )
	// multi := zerolog.MultiLevelWriter(runLogFile)
	// log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()

	lumberjackLogger := &lumberjack.Logger{
		Filename:   "log/app.log",
		MaxSize:    10,   // Max size in megabytes before log is rotated
		MaxBackups: 3,    // Max number of old log files to retain
		MaxAge:     28,   // Max number of days to retain old log files
		Compress:   true, // Compress old log files
	}

	log.Logger = zerolog.New(lumberjackLogger).With().Timestamp().Logger()
}
