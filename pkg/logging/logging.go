package logging

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"runtime"
)

type Logger struct {
	*logrus.Entry
}

func Init(isProd bool) (*Logger, error) {
	logsLevel := "debug"
	if isProd {
		logsLevel = "info"
	}

	logrusLogsLevel, err := logrus.ParseLevel(logsLevel)
	if err != nil {
		return nil, errors.New("cannot parse log level")
	}

	l := logrus.New()
	l.SetReportCaller(true)
	l.Formatter = &logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return f.Function, fmt.Sprintf("%s:%d", filename, f.Line)
		},
		DisableColors: false,
		FullTimestamp: true,
	}

	l.SetOutput(os.Stdout) // Send all logs to stdout

	l.SetLevel(logrusLogsLevel)

	return &Logger{
		logrus.NewEntry(l),
	}, nil
}
