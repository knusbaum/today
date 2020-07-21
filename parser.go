package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type parser struct {
	rdr   *bufio.Reader
	peekp bool
	peek  string
}

const (
	StartupLine = "Morning Start Up:"
	NotesLine   = "Notes:"
	LogLine     = "Log:"
	TODOLine    = "TODO:"
)

func (p *parser) PeekLine() (string, error) {
	if p.peekp {
		//log.Printf("PEEK LINE1: [%s]", p.peek)
		return p.peek, nil
	}
	peek, err := p.rdr.ReadString('\n')
	if strings.HasSuffix(peek, "\n") {
		peek = peek[:len(peek)-1]
	}
	p.peek = peek
	p.peekp = true
	//log.Printf("PEEK LINE2: [%s]", p.peek)
	return p.peek, err
}

func (p *parser) NextLine() (string, error) {
	if p.peekp {
		r := p.peek
		p.peek = ""
		p.peekp = false
		//log.Printf("NEXT LINE1: [%s]", r)
		return r, nil
	}
	str, err := p.rdr.ReadString('\n')
	if strings.HasSuffix(str, "\n") {
		str = str[0 : len(str)-1]
	}
	//log.Printf("NEXT LINE2: [%s]", str)
	return str, err
}

func (p *parser) UngetLine(l string) {
	if p.peekp {
		panic("Cannot unget more than one line.")
	}
	p.peek = l
	p.peekp = true
}

func matchLine(line, match string) bool {
	ret := strings.HasPrefix(strings.TrimSpace(line), match)
	//log.Printf("[matchLine] [%s] == [%s] -> %t", strings.TrimSpace(line), match, ret)
	return ret
}

func parseStatus(s string) status {
	re := regexp.MustCompile(`(([A-Z-? ]*?)([[:space:]]+-[[:space:]]+|$))?(.*?)([[:space:]]+-[[:space:]]+(.*?))?$`)
	matches := re.FindStringSubmatch(s)

	name := matches[2]
	part2 := matches[4]
	part3 := matches[6]

	if part3 == "" {
		if date, err := time.ParseInLocation("Jan _2, 2006", part2, time.Local); err == nil {
			// part 2 is a date.
			return status{
				name: name,
				date: date,
			}
		}
		return status{
			name:    name,
			comment: part2 + matches[5],
		}
	}

	date, err := time.ParseInLocation("Jan _2, 2006", matches[6], time.Local)
	if err != nil {
		return status{
			name:    name,
			comment: part2 + matches[5],
		}
	}
	return status{
		name:    name,
		comment: part2,
		date:    date,
	}
}

func (p *parser) parseTodo() *todo {
	var t todo
	//log.Printf("[parseTodo] next")
	l, err := p.NextLine()
	if err != nil {
		log.Printf("Warning: Unexpected error while parsing todo: %s\n", err)
		return nil
	}
	l = strings.TrimSpace(l)
	if l == "" {
		return nil
	}

	re := regexp.MustCompile(`^(([A-Z]+-[0-9]+)[[:space:]]+-)?(.*?)(\[([^][]*)\])?$`)
	matches := re.FindStringSubmatch(l)

	jira := strings.TrimSpace(matches[2])
	description := strings.TrimSpace(matches[3])
	status := parseStatus(strings.TrimSpace(matches[5]))

	//log.Printf("JIRA: %s", jira)
	//log.Printf("DESCR: %s", description)
	//log.Printf("Status: %#v", status)

	t.jira = jira
	t.description = description
	t.status = status
	//log.Printf("[parseTodo] reading comments")
	//defer log.Printf("[parseTodo] done")
	for l, err := p.NextLine(); err == nil; l, err = p.NextLine() {
		if strings.HasPrefix(l, "\t") {
			t.comments = append(t.comments, strings.TrimSpace(l))
		} else if strings.TrimSpace(l) == "" {
			t.blankBelow = true
			continue
		} else {
			p.UngetLine(l)
			return &t
		}
	}
	return &t
}

func (p *parser) parseListItem() *listItem {
	l, err := p.NextLine()
	if err != nil {
		log.Printf("Warning: Unexpected error while parsing todo: %s\n", err)
		return nil
	}
	l = strings.TrimSpace(l)
	if l == "" {
		return nil
	}
	re := regexp.MustCompile(`^(([0-9]+)\.)?[[:space:]]*(.*?)(\[([^][]*)\])?$`)
	matches := re.FindStringSubmatch(l)

	itemNumber, err := strconv.Atoi(matches[2])
	if err != nil {
		panic(err) // TODO: Probably don't want to panic.
	}
	comment := strings.TrimSpace(matches[3])
	status := parseStatus(strings.TrimSpace(matches[5]))

	return &listItem{number: itemNumber, description: comment, status: status}
}

func (p *parser) parseList(nextSection string) []*listItem {
	var items []*listItem
	for {
		//log.Printf("[parseTodos] PEEKING for next section")
		l, err := p.PeekLine()
		if err != nil || matchLine(l, nextSection) {
			//log.Printf("[parseTodos] found next section %s", nextSection)
			return items
		}
		//log.Printf("[parseTodos] parsing todo.")
		if item := p.parseListItem(); item != nil {
			items = append(items, item)
		}
		//log.Printf("[parseTodos] done parsing todo.")
	}
}

func (p *parser) parseTodos(nextSection string) []*todo {
	var todos []*todo
	for {
		//log.Printf("[parseTodos] PEEKING for next section")
		l, err := p.PeekLine()
		if err != nil || matchLine(l, nextSection) {
			//log.Printf("[parseTodos] found next section %s", nextSection)
			if len(todos) > 0 {
				todos[len(todos)-1].blankBelow = false // last todo never gets blank line.
			}
			return todos
		}
		//log.Printf("[parseTodos] parsing todo.")
		if t := p.parseTodo(); t != nil {
			todos = append(todos, t)
		}
		//log.Printf("[parseTodos] done parsing todo.")
	}
}

func (p *parser) parseLines(nextSection string) []string {
	var lines []string
	for {
		l, err := p.PeekLine()
		if err != nil || matchLine(l, nextSection) {
			//log.Printf("[parseLines] found next section %s", nextSection)
			return lines
		}
		p.NextLine()
		if strings.TrimSpace(l) != "" {
			// Trim off the newline.
			lines = append(lines, l)
		}
	}
}

func (p *parser) parseStartup() ([]*listItem, error) {
	for {
		l, err := p.NextLine()
		if err != nil {
			return nil, err
		}
		//log.Printf("[parseStartup] looking for section %s", StartupLine)
		if matchLine(l, StartupLine) {
			//log.Printf("[parseStartup] found section %s. parsing todos.", StartupLine)
			l := p.parseList(NotesLine)
			//log.Printf("[parseStartup] todos: %#v", todos)
			return l, nil
		}
	}
}

func (p *parser) parseNotes() ([]string, error) {
	for {
		l, err := p.NextLine()
		if err != nil {
			return nil, err
		}
		//log.Printf("[parseNotes] looking for section %s", NotesLine)
		if matchLine(l, NotesLine) {
			//log.Printf("[parseNotes] found section %s. parsing lines.", NotesLine)
			return p.parseLines(LogLine), nil
		}
	}
}

func (p *parser) parseLog() ([]string, error) {
	for {
		l, err := p.NextLine()
		if err != nil {
			return nil, err
		}
		//log.Printf("[parseLog] looking for section %s", LogLine)
		if matchLine(l, LogLine) {
			//log.Printf("[parseLog] found section %s. parsing lines.", LogLine)
			return p.parseLines(TODOLine), nil
		}
	}
}

func (p *parser) parseTODO() ([]*todo, error) {
	for {
		l, err := p.NextLine()
		if err != nil {
			return nil, err
		}
		//log.Printf("[parseNotes] looking for section %s", TODOLine)
		if matchLine(l, TODOLine) {
			//log.Printf("[parseTODO] found section %s. parsing todos.", TODOLine)
			return p.parseTodos("END"), nil
		}
	}
}

func NewParser(r io.Reader) *parser {
	var rdr *bufio.Reader
	if br, ok := r.(*bufio.Reader); ok {
		rdr = br
	} else {
		rdr = bufio.NewReader(r)
	}
	return &parser{rdr: rdr}
}

func (p *parser) parseToday() (*today, error) {
	var t today

	startup, err := p.parseStartup()
	if err != nil {
		return nil, fmt.Errorf("Failed to parse startup: %s", err)
	}
	t.startup = startup

	notes, err := p.parseNotes()
	if err != nil {
		return nil, fmt.Errorf("Failed to parse notes: %s", err)
	}
	t.notes = notes

	log, err := p.parseLog()
	if err != nil {
		return nil, fmt.Errorf("Failed to parse log: %s", err)
	}
	t.log = log

	todos, err := p.parseTODO()
	if err != nil {
		return nil, fmt.Errorf("Failed to parse todo: %s", err)
	}
	t.todos = todos

	for _, todo := range t.todos {
		if strings.HasPrefix(todo.jira, "TODO-") {
			taskid, err := strconv.Atoi(strings.TrimPrefix(todo.jira, "TODO-"))
			if err == nil && taskid >= t.nextTodo {
				t.nextTodo = taskid + 1
			}
		}
	}

	return &t, nil
}
