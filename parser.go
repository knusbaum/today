package today

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

// Parser parses a today file
type Parser struct {
	rdr   *bufio.Reader
	peekp bool
	peek  string
}

const (
	startupLine = "Morning Start Up:"
	notesLine   = "Notes:"
	logLine     = "Log:"
	todoLine    = "TODO:"
)

func (p *Parser) peekLine() (string, error) {
	if p.peekp {
		return p.peek, nil
	}
	peek, err := p.rdr.ReadString('\n')
	if strings.HasSuffix(peek, "\n") {
		peek = peek[:len(peek)-1]
	}
	p.peek = peek
	p.peekp = true
	return p.peek, err
}

func (p *Parser) nextLine() (string, error) {
	if p.peekp {
		r := p.peek
		p.peek = ""
		p.peekp = false
		return r, nil
	}
	str, err := p.rdr.ReadString('\n')
	if strings.HasSuffix(str, "\n") {
		str = str[0 : len(str)-1]
	}
	return str, err
}

func (p *Parser) ungetLine(l string) {
	if p.peekp {
		panic("Cannot unget more than one line.")
	}
	p.peek = l
	p.peekp = true
}

func matchLine(line, match string) bool {
	ret := strings.HasPrefix(strings.TrimSpace(line), match)
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

func (p *Parser) parseTodo() *todo {
	var t todo
	l, err := p.nextLine()
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
	t.jira = strings.TrimSpace(matches[2])
	t.description = strings.TrimSpace(matches[3])
	t.status = parseStatus(strings.TrimSpace(matches[5]))
	for l, err := p.nextLine(); err == nil; l, err = p.nextLine() {
		if strings.HasPrefix(l, "\t") {
			t.comments = append(t.comments, strings.TrimSpace(l))
		} else if strings.TrimSpace(l) == "" {
			t.blankBelow = true
			continue
		} else {
			p.ungetLine(l)
			return &t
		}
	}
	return &t
}

func (p *Parser) parseListItem() *listItem {
	l, err := p.nextLine()
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

	var itemNumber int
	if matches[2] != "" {
		itemNumber, err = strconv.Atoi(matches[2])
		if err != nil {
			log.Printf("Failed to parse list item number: %s", err)
		}
	}
	comment := strings.TrimSpace(matches[3])
	status := parseStatus(strings.TrimSpace(matches[5]))

	return &listItem{number: itemNumber, description: comment, status: status}
}

func (p *Parser) parseList(nextSection string) []*listItem {
	var items []*listItem
	for {
		l, err := p.peekLine()
		if err != nil || matchLine(l, nextSection) {
			return items
		}
		if item := p.parseListItem(); item != nil {
			items = append(items, item)
		}
	}
}

func (p *Parser) parseTodos(nextSection string) []*todo {
	var todos []*todo
	for {
		l, err := p.peekLine()
		if err != nil || matchLine(l, nextSection) {
			if len(todos) > 0 {
				todos[len(todos)-1].blankBelow = false // last todo never gets blank line.
			}
			return todos
		}
		if t := p.parseTodo(); t != nil {
			todos = append(todos, t)
		}
	}
}

func (p *Parser) parseLines(nextSection string) []string {
	var lines []string
	for {
		l, err := p.peekLine()
		if err != nil || matchLine(l, nextSection) {
			return lines
		}
		p.nextLine()
		if strings.TrimSpace(l) != "" {
			lines = append(lines, l)
		}
	}
}

func (p *Parser) parseStartup() ([]*listItem, error) {
	for {
		l, err := p.nextLine()
		if err != nil {
			return nil, err
		}
		if matchLine(l, startupLine) {
			l := p.parseList(notesLine)
			return l, nil
		}
	}
}

func (p *Parser) parseNotes() ([]string, error) {
	for {
		l, err := p.nextLine()
		if err != nil {
			return nil, err
		}
		if matchLine(l, notesLine) {
			return p.parseLines(logLine), nil
		}
	}
}

func (p *Parser) parseLog() ([]string, error) {
	for {
		l, err := p.nextLine()
		if err != nil {
			return nil, err
		}
		if matchLine(l, logLine) {
			return p.parseLines(todoLine), nil
		}
	}
}

func (p *Parser) parseTODO() ([]*todo, error) {
	for {
		l, err := p.nextLine()
		if err != nil {
			return nil, err
		}
		if matchLine(l, todoLine) {
			return p.parseTodos("END"), nil
		}
	}
}

// NewParser will create a new
func NewParser(r io.Reader) *Parser {
	var rdr *bufio.Reader
	if br, ok := r.(*bufio.Reader); ok {
		rdr = br
	} else {
		rdr = bufio.NewReader(r)
	}
	return &Parser{rdr: rdr}
}

func (p *Parser) Parse() (*Today, error) {
	var t Today

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
