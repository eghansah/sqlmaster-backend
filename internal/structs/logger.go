package structs

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type LogFile struct {
	LogFileName string
	LogPath     string
}

func (l *LogFile) Write(p []byte) (n int, err error) {
	var pathError *os.PathError
	now := time.Now()
	logpath := fmt.Sprintf("%s/%s", l.LogPath, now.Format("2006/01"))
	logfile := fmt.Sprintf("%s_%s.log", l.LogFileName, strings.ToUpper(now.Format("02JAN2006")))
	fullLogPath := fmt.Sprintf("%s/%s", logpath, logfile)
	f, err := os.OpenFile(fullLogPath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if errors.As(err, &pathError) {
		if err := os.MkdirAll(logpath, 0755); err != nil {
			log.Println(err)
		}

		f, err = os.OpenFile(fullLogPath,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err)
		}
	} else if err != nil {
		log.Println(err)
	}
	defer f.Close()

	if _, err := f.Write(p); err != nil {
		return 0, err
	}

	return len(p), nil
}

func (l *LogFile) Sync() error {
	return nil
}
