package guestbook

import (
	"html/template"
	"os"
	"strings"
	"time"

	"github.com/gomarkdown/markdown"
)

const guestbookFile = "guestbook.txt"

type Entry struct {
	Timestamp       string
	TimestampPretty string
	Name            string
	Message         template.HTML
}

func SanitizeAndFormat(message string) template.HTML {
	message = template.HTMLEscapeString(message)
	res := markdown.ToHTML([]byte(message), nil, nil)
	return template.HTML(res)
}

func Save(name string, message string) error {
	// trim spaces
	name = strings.TrimSpace(name)
	message = strings.TrimSpace(message)

	// if name or message empty then ignore
	if name == "" || message == "" {
		return nil
	}

	// concatenate string
	var newEntry string = "[" + time.Now().Format(time.RFC3339) + "] " + name + ": " + message + "\n"

	// open file
	f, err := os.OpenFile(guestbookFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(newEntry)
	return err
}

func Show() ([]Entry, error) {
	// read file
	data, err := os.ReadFile(guestbookFile)
	if err != nil {
		return nil, err
	}

	// split all guestbook entries
	var lines []string = strings.Split(string(data), "\n")
	var entries []Entry

	for _, line := range lines {
		// if current line is empty then skip
		if strings.TrimSpace(line) == "" {
			continue
		}

		// split current line in different parts
		var parts []string = strings.SplitN(line, "] ", 2)
		if len(parts) != 2 {
			continue
		}

		var timestamp string = strings.TrimPrefix(parts[0], "[")
		var rest []string = strings.SplitN(parts[1], ": ", 2)

		// make timestamp user readable
		parsedTime, err := time.Parse(time.RFC3339, timestamp)
		if err != nil {
			continue
		}
		prettyTime := parsedTime.Format(time.RFC822)

		var name string = rest[0]
		var message template.HTML = SanitizeAndFormat(rest[1])

		entries = append(entries, Entry{
			Timestamp:       timestamp,
			TimestampPretty: prettyTime,
			Name:            name,
			Message:         message,
		})
	}

	return entries, nil
}
