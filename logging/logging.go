package logging

import (
	"os"
	"time"
)

func Log(msg string) error {
	f, err := os.OpenFile("app.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	line := time.Now().Format("2006-01-02 15:04:05") + " " + msg + "\n"

	_, err = f.WriteString(line)
	return err
}
