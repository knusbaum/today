package today

import (
	"bufio"
	"fmt"
	"io"
)

func writeStatus(s *Status, w *bufio.Writer) error {
	var (
		statusStr string = "["
		wrotename bool
	)

	if s.Name != "" {
		statusStr += s.Name
		wrotename = true
	}

	if s.Comment != "" {
		if wrotename {
			statusStr += " - "
		}
		statusStr += s.Comment
	}

	if !s.Date.IsZero() {
		statusStr += " - "
		statusStr += s.Date.Format("Jan _2, 2006")
	}

	statusStr += "]"
	_, err := w.WriteString(statusStr)
	if err != nil {
		return err
	}
	return nil
}

func writeTodo(t *Task, w *bufio.Writer) error {
	if t.Name != "" {
		_, err := w.WriteString(t.Name + " - ")
		if err != nil {
			return err
		}
	}

	if t.Description != "" {
		_, err := w.WriteString(t.Description + " ")
		if err != nil {
			return err
		}
	}

	if t.Status.Name != "" || t.Status.Comment != "" {
		err := writeStatus(&t.Status, w)
		if err != nil {
			return err
		}
	}
	_, err := w.WriteString("\n")
	if err != nil {
		return err
	}

	for _, c := range t.Comments {
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

func writeListItem(item *ListItem, w *bufio.Writer) error {
	_, err := w.WriteString(fmt.Sprintf("%d. %s ", item.number, item.Description))
	if err != nil {
		return err
	}
	if item.Status.Name != "" || item.Status.Comment != "" {
		err := writeStatus(&item.Status, w)
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

func (t *TaskList) Write(w *bufio.Writer) error {
	for _, todo := range t.Tasks {
		err := writeTodo(todo, w)
		if err != nil {
			return err
		}
	}
	return nil
}

// Write writes a Today out to writer w in the normal form
func (t *Today) Write(w io.Writer) error {
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
	for _, item := range t.Startup {
		err = writeListItem(item, wtr)
		if err != nil {
			return err
		}
	}

	_, err = wtr.WriteString("\n" + notesLine + "\n")
	if err != nil {
		return err
	}
	for _, n := range t.Notes {
		_, err = wtr.WriteString(n + "\n")
		if err != nil {
			return err
		}
	}

	_, err = wtr.WriteString(logLine + "\n")
	if err != nil {
		return err
	}
	for _, l := range t.Log {
		_, err = wtr.WriteString(l + "\n")
		if err != nil {
			return err
		}
	}

	_, err = wtr.WriteString("\n" + todoLine + "\n")
	if err != nil {
		return err
	}
	err = t.Tasks.Write(wtr)
	if err != nil {
		return err
	}

	_, err = wtr.WriteString("\n")
	if err != nil {
		return err
	}

	wtr.Flush()
	return nil
}
