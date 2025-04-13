package guestbook

import (
	"fmt"
	"html/template"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/microcosm-cc/bluemonday"
)

var guestbookFile string = "guestbook.txt"

type Entry struct {
	Timestamp       string
	TimestampPretty string
	Name            string
	Message         template.HTML
}

func renderMarkdown(message string) template.HTML {
	md := markdown.ToHTML([]byte(message), nil, nil)

	// convert \n to <br>, but only inside <p>...</p>
	withBreaks := addLineBreaks(string(md))

	// clean HTML (XSS protection)
	p := bluemonday.UGCPolicy()
	safe := p.Sanitize(withBreaks)

	return template.HTML(safe)
}

func addLineBreaks(input string) string {
	// RegEx search <p>...</p> and replace \n durch <br>
	re := regexp.MustCompile(`(?s)<p>(.*?)</p>`)
	return re.ReplaceAllStringFunc(input, func(p string) string {
		// extracting content
		content := strings.TrimSuffix(strings.TrimPrefix(p, "<p>"), "</p>")
		// replace \n to <br>
		content = strings.ReplaceAll(content, "\n", "<br>\n")
		// wrap it up again
		return "<p>" + content + "</p>"
	})
}

func Save(name string, message string) error {
	// trim spaces
	name = strings.TrimSpace(name)
	message = strings.TrimSpace(message)

	// if name or message empty then ignore
	if name == "" || message == "" {
		return nil
	}

	entry := fmt.Sprintf(
		"---ENTRY---\n%s\n%s\n%s\n",
		time.Now().Format(time.RFC3339),
		name,
		message,
	)

	// open file
	f, err := os.OpenFile(guestbookFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(entry)
	return err
}

func Show() ([]Entry, error) {
	// read file
	data, err := os.ReadFile(guestbookFile)
	if err != nil {
		return nil, err
	}

	blocks := strings.Split(string(data), "---ENTRY---")
	var entries []Entry

	for _, block := range blocks {
		block = strings.TrimSpace(block)
		if block == "" {
			continue
		}

		lines := strings.SplitN(block, "\n", 3)
		if len(lines) < 3 {
			continue
		}

		timestampStr := lines[0]
		name := lines[1]
		message := lines[2]

		timestamp, err := time.Parse(time.RFC3339, timestampStr)
		if err != nil {
			continue
		}

		entry := Entry{
			Timestamp:       timestampStr,
			TimestampPretty: timestamp.Format(time.RFC822),
			Name:            name,
			Message:         renderMarkdown(message),
		}

		entries = append(entries, entry)
	}

	return entries, nil
}
