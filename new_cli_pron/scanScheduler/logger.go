package scanScheduler

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

type logger struct {
	logfilePath string
	file        *os.File
}

func newLogger(logPath string) logger {
	f, err := os.OpenFile(
		logPath,
		os.O_CREATE|os.O_APPEND|os.O_RDWR,
		os.ModePerm,
	)
	if err != nil {
		// What to do here? We can hardly log it! :)
		fmt.Println("Failed to open file")
	}
	return logger{
		logPath,
		f,
	}
}

func (l *logger) addEntry(entry string) error {
	w := bufio.NewWriter(l.file)
	const (
		layout = "Mon Jan 2 15:04:05"
	)
	augEntry := time.Now().Format(layout) + ": " + entry
	fmt.Fprintln(w, augEntry)

	err := w.Flush()

	if err != nil {
		// We failed to write to the log so let's try to recreate it
		l1 := newLogger(l.logfilePath)

		l.file = l1.file

		// And try to write the entry again...
		w := bufio.NewWriter(l.file)
		fmt.Fprintln(w, augEntry)

		err = w.Flush()
	}
	return err
}

// func (l logger) close() {
// 	l.file.Close()
// }
