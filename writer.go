package main

import (
	"bufio"
	"fmt"
	"io"
)

func writeStatus(s *status, w *bufio.Writer) error {
	var (
		statusStr string = "["
		wrotename bool
	)

	if s.name != "" {
		statusStr += s.name
		wrotename = true
	}

	if s.comment != "" {
		if wrotename {
			statusStr += " - "
		}
		statusStr += s.comment
	}

	if !s.date.IsZero() {
		statusStr += " - "
		statusStr += s.date.Format("Jan _2, 2006")
	}

	statusStr += "]"
	_, err := w.WriteString(statusStr)
	if err != nil {
		return err
	}
	return nil
}

func writeTodo(t *todo, w *bufio.Writer) error {
	if t.jira != "" {
		_, err := w.WriteString(t.jira + " - ")
		if err != nil {
			return err
		}
	}

	if t.description != "" {
		_, err := w.WriteString(t.description + " ")
		if err != nil {
			return err
		}
	}

	if t.status.name != "" || t.status.comment != "" {
		err := writeStatus(&t.status, w)
		if err != nil {
			return err
		}
	}
	_, err := w.WriteString("\n")
	if err != nil {
		return err
	}

	for _, c := range t.comments {
		_, err := w.WriteString("\t" + c + "\n")
		if err != nil {
			return err
		}
	}
	if t.blankBelow {
		_, err = w.WriteString("\n")
		if err != nil {
			return err
		}
	}

	return nil
}

func writeListItem(item *listItem, w *bufio.Writer) error {
	_, err := w.WriteString(fmt.Sprintf("%d. %s ", item.number, item.description))
	if item.status.name != "" || item.status.comment != "" {
		err := writeStatus(&item.status, w)
		if err != nil {
			return err
		}
	}
	_, err = w.WriteString("\n")
	if err != nil {
		return err
	}
	return nil
}

func writeToday(t *today, w io.Writer) error {
	var wtr *bufio.Writer
	if bw, ok := w.(*bufio.Writer); ok {
		wtr = bw
	} else {
		wtr = bufio.NewWriter(w)
	}

	_, err := wtr.WriteString(startupLine + "\n")
	if err != nil {
		return err
	}
	for _, item := range t.startup {
		err = writeListItem(item, wtr)
		if err != nil {
			return err
		}
	}

	_, err = wtr.WriteString("\n" + notesLine + "\n")
	if err != nil {
		return err
	}
	for _, n := range t.notes {
		_, err = wtr.WriteString(n + "\n")
		if err != nil {
			return err
		}
	}

	_, err = wtr.WriteString("\n" + logLine + "\n")
	if err != nil {
		return err
	}
	for _, l := range t.log {
		_, err = wtr.WriteString(l + "\n")
		if err != nil {
			return err
		}
	}

	_, err = wtr.WriteString("\n" + todoLine + "\n")
	if err != nil {
		return err
	}
	for _, todo := range t.todos {
		err = writeTodo(todo, wtr)
		if err != nil {
			return err
		}
	}

	_, err = wtr.WriteString("\n")
	if err != nil {
		return err
	}

	wtr.Flush()
	return nil
}
