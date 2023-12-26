package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

// lg Общая переменная для логирования будет доступна всему коду
// Не лучшее решение, но самое простое
var lg *zerolog.Logger

func Init(level string) {

	// Пока используем консольный вывод
	//zl := zerolog.New(os.Stderr).With().Timestamp().Logger()
	//file, err := os.Create("log" + time.Now().String() + ".log")
	//if err != nil {
	//	panic(fmt.Errorf("error due create log file %w", err))
	//}
	//zl := log.Output(file)
	zl := log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	switch level {
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)

	default:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		zl.Debug().Msg("passed wrong error level")
	}
	zl.Debug().Str("log level", level).Msg("log initialized")

	lg = &zl
}

func Fatal() *zerolog.Event {
	return lg.Fatal()
}
func Error() *zerolog.Event {
	return lg.Error()
}
func Info() *zerolog.Event {
	return lg.Info()
}
func Debug() *zerolog.Event {
	return lg.Debug()
}
