package today

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseMorning(t *testing.T) {
	r := strings.NewReader(`Morning Start Up:
1. Catch up on slack 
2. Check the calendar 



Read the inbox [DONE - something - Jan 5, 2020]



look at JIRAPROJECT 


Notes:
foop boop doop

Log:

TODO:

`)
	p := NewParser(r)
	today, err := p.Parse()
	assert.NoError(t, err)
	if !assert.NotNil(t, today) {
		return
	}

	if !assert.Len(t, today.startup, 4) {
		return
	}
	assert.Len(t, today.startup, 4)
	assert.Equal(t, 1, today.startup[0].number)
	assert.Equal(t, "Catch up on slack", today.startup[0].description)
	assert.Equal(t, 2, today.startup[1].number)
	assert.Equal(t, "Check the calendar", today.startup[1].description)
	assert.Equal(t, 0, today.startup[2].number)
	assert.Equal(t, "Read the inbox", today.startup[2].description)
	assert.Equal(t, "DONE", today.startup[2].status.name)
	assert.Equal(t, "something", today.startup[2].status.comment)
	assert.Equal(t, time.Date(2020, 1, 5, 0, 0, 0, 0, time.Local), today.startup[2].status.date)
	assert.Equal(t, 0, today.startup[3].number)
	assert.Equal(t, "look at JIRAPROJECT", today.startup[3].description)
}

func TestParseNotes(t *testing.T) {
	r := strings.NewReader(`Morning Start Up:
1. Catch up on slack
2. Check the calendar 


Notes:
foop boop doop This is a note.


note important_facts 
note testing

note deploy

note custom_build

Log:
8:30 Starting work 

TODO:

`)
	p := NewParser(r)
	today, err := p.Parse()
	assert.NoError(t, err)
	if !assert.NotNil(t, today) {
		return
	}

	if !assert.Len(t, today.notes, 5) {
		return
	}
}

func TestParseLog(t *testing.T) {
	r := strings.NewReader(`Morning Start Up:
1. Catch up on slack
2. Check the calendar 

Notes:
foop boop doop This is a note.

Log:
8:30 Starting work 
9:00 Standup
9:15 Starting on Jira task

TODO:
`)
	p := NewParser(r)
	today, err := p.Parse()
	assert.NoError(t, err)
	if !assert.NotNil(t, today) {
		return
	}

	if !assert.Len(t, today.log, 3) {
		return
	}
}

func TestParseTODO(t *testing.T) {
	r := strings.NewReader(`Morning Start Up:
1. Catch up on slack
2. Check the calendar 

Notes:
foop boop doop This is a note.

Log:
8:30 Starting work 

TODO:
SOMEJIRA-1234 - Do something for this jira. [IN-PROGRESS]
	* Some note.
	* Some other note.
TODO-240 - Another Task [WAITING - waiting to hear from client multi-hyphen-word - Jan 5, 2020]
Yet another task.
`)
	p := NewParser(r)
	today, err := p.Parse()
	assert.NoError(t, err)
	if !assert.NotNil(t, today) {
		return
	}

	if !assert.Len(t, today.todos, 3) {
		return
	}

	assert.Equal(t, "SOMEJIRA-1234", today.todos[0].jira)
	assert.Equal(t, "Do something for this jira.", today.todos[0].description)
	assert.Equal(t, "IN-PROGRESS", today.todos[0].status.name)
	if !assert.Len(t, today.todos[0].comments, 2) {
		return
	}
	assert.Equal(t, "* Some note.", today.todos[0].comments[0])
	assert.Equal(t, "* Some other note.", today.todos[0].comments[1])

	assert.Equal(t, "TODO-240", today.todos[1].jira)
	assert.Equal(t, "Another Task", today.todos[1].description)
	assert.Equal(t, "WAITING", today.todos[1].status.name)
	assert.Equal(t, "waiting to hear from client multi-hyphen-word", today.todos[1].status.comment)
	assert.Equal(t, time.Date(2020, 1, 5, 0, 0, 0, 0, time.Local), today.todos[1].status.date)
	assert.Equal(t, 241, today.nextTodo)

	assert.Equal(t, "", today.todos[2].jira)
	assert.Equal(t, "Yet another task.", today.todos[2].description)
	assert.Equal(t, "", today.todos[2].status.name)
	assert.Equal(t, "", today.todos[2].status.comment)
	assert.True(t, today.todos[2].status.date.IsZero())
}

func TestParseStatus(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		s := parseStatus("basic status")
		assert.Equal(t, "basic status", s.comment)
	})
	t.Run("status", func(t *testing.T) {
		s := parseStatus("IN PROGRESS")
		assert.Equal(t, "IN PROGRESS", s.name)
	})
	t.Run("basic+status", func(t *testing.T) {
		s := parseStatus("IN PROGRESS - basic status")
		assert.Equal(t, "IN PROGRESS", s.name)
		assert.Equal(t, "basic status", s.comment)
	})
	t.Run("full", func(t *testing.T) {
		s := parseStatus("WAITING FOR CUSTOMER - waiting to hear from client multi-hyphen-word - Jan 5, 2020")
		assert.Equal(t, "WAITING FOR CUSTOMER", s.name)
		assert.Equal(t, "waiting to hear from client multi-hyphen-word", s.comment)
		assert.Equal(t, time.Date(2020, 1, 5, 0, 0, 0, 0, time.Local), s.date)
	})
	t.Run("status+date", func(t *testing.T) {
		s := parseStatus("WAITING FOR CUSTOMER - Jan 5, 2020")
		assert.Equal(t, "WAITING FOR CUSTOMER", s.name)
		assert.Equal(t, "", s.comment)
		assert.Equal(t, time.Date(2020, 1, 5, 0, 0, 0, 0, time.Local), s.date)
	})
	t.Run("bad-date", func(t *testing.T) {
		s := parseStatus("WAITING FOR CUSTOMER - waiting to hear from client multi-hyphen-word - Jan 335, 2020")
		assert.Equal(t, "WAITING FOR CUSTOMER", s.name)
		assert.Equal(t, "waiting to hear from client multi-hyphen-word - Jan 335, 2020", s.comment)
	})
	t.Run("comment-brackets", func(t *testing.T) {
		p := NewParser(strings.NewReader("JIRAPROJECT-123 - [Client X] - Can't frobnicate the blips  [STALE - Jun 10, 2020]\n"))
		todo := p.parseTodo()
		assert.Equal(t, "JIRAPROJECT-123", todo.jira)
		assert.Equal(t, "[Client X] - Can't frobnicate the blips", todo.description)
		assert.Equal(t, "STALE", todo.status.name)
	})
}

func TestParseLines(t *testing.T) {
	r := strings.NewReader(`line
another line

	yet another line
line again

one more line

END
more lines
`)
	p := NewParser(r)
	lines := p.parseLines("END")
	assert.Equal(t,
		[]string{
			"line",
			"another line",
			"\tyet another line",
			"line again",
			"one more line",
		},
		lines,
	)
}
