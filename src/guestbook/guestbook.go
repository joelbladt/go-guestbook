package guestbook

import (
	"os"
	"strings"
	"time"
)

const guestbookFile = "guestbook.txt"

func Save(name string, message string) error {
	// trim spaces
	name = strings.TrimSpace(name)
	message = strings.TrimSpace(message)

	// if name or message empty then ignore
	if name == "" || message == "" {
		return nil
	}

	// concatenate string
	var newEntry string = "[" + time.Now().Format(time.RFC822) + "] " + name + ": " + message + "\n"

	// open file
	f, err := os.OpenFile(guestbookFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(newEntry)
	return err
}

func Show() {
}
