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

	if !assert.Len(t, today.Startup, 4) {
		return
	}
	assert.Len(t, today.Startup, 4)
	assert.Equal(t, 1, today.Startup[0].number)
	assert.Equal(t, "Catch up on slack", today.Startup[0].Description)
	assert.Equal(t, 2, today.Startup[1].number)
	assert.Equal(t, "Check the calendar", today.Startup[1].Description)
	assert.Equal(t, 0, today.Startup[2].number)
	assert.Equal(t, "Read the inbox", today.Startup[2].Description)
	assert.Equal(t, "DONE", today.Startup[2].Status.Name)
	assert.Equal(t, "something", today.Startup[2].Status.Comment)
	assert.Equal(t, time.Date(2020, 1, 5, 0, 0, 0, 0, time.Local), today.Startup[2].Status.Date)
	assert.Equal(t, 0, today.Startup[3].number)
	assert.Equal(t, "look at JIRAPROJECT", today.Startup[3].Description)
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

	if !assert.Len(t, today.Notes, 5) {
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

	if !assert.Len(t, today.Log, 3) {
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
TASK-240 - Another Task [WAITING - waiting to hear from client multi-hyphen-word - Jan 5, 2020]
Yet another task.
`)
	p := NewParser(r)
	today, err := p.Parse()
	assert.NoError(t, err)
	if !assert.NotNil(t, today) {
		return
	}

	if !assert.Len(t, today.Tasks.tasks, 3) {
		return
	}

	assert.Equal(t, "SOMEJIRA-1234", today.Tasks.tasks[0].Name)
	assert.Equal(t, "Do something for this jira.", today.Tasks.tasks[0].Description)
	assert.Equal(t, "IN-PROGRESS", today.Tasks.tasks[0].Status.Name)
	if !assert.Len(t, today.Tasks.tasks[0].Comments, 2) {
		return
	}
	assert.Equal(t, "* Some note.", today.Tasks.tasks[0].Comments[0])
	assert.Equal(t, "* Some other note.", today.Tasks.tasks[0].Comments[1])

	assert.Equal(t, "TASK-240", today.Tasks.tasks[1].Name)
	assert.Equal(t, "Another Task", today.Tasks.tasks[1].Description)
	assert.Equal(t, "WAITING", today.Tasks.tasks[1].Status.Name)
	assert.Equal(t, "waiting to hear from client multi-hyphen-word", today.Tasks.tasks[1].Status.Comment)
	assert.Equal(t, time.Date(2020, 1, 5, 0, 0, 0, 0, time.Local), today.Tasks.tasks[1].Status.Date)
	assert.Equal(t, 241, today.Tasks.nextTaskID)

	assert.Equal(t, "", today.Tasks.tasks[2].Name)
	assert.Equal(t, "Yet another task.", today.Tasks.tasks[2].Description)
	assert.Equal(t, "", today.Tasks.tasks[2].Status.Name)
	assert.Equal(t, "", today.Tasks.tasks[2].Status.Comment)
	assert.True(t, today.Tasks.tasks[2].Status.Date.IsZero())
}

func TestParseStatus(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		s := parseStatus("basic status")
		assert.Equal(t, "basic status", s.Comment)
	})
	t.Run("status", func(t *testing.T) {
		s := parseStatus("IN PROGRESS")
		assert.Equal(t, "IN PROGRESS", s.Name)
	})
	t.Run("basic+status", func(t *testing.T) {
		s := parseStatus("IN PROGRESS - basic status")
		assert.Equal(t, "IN PROGRESS", s.Name)
		assert.Equal(t, "basic status", s.Comment)
	})
	t.Run("full", func(t *testing.T) {
		s := parseStatus("WAITING FOR CUSTOMER - waiting to hear from client multi-hyphen-word - Jan 5, 2020")
		assert.Equal(t, "WAITING FOR CUSTOMER", s.Name)
		assert.Equal(t, "waiting to hear from client multi-hyphen-word", s.Comment)
		assert.Equal(t, time.Date(2020, 1, 5, 0, 0, 0, 0, time.Local), s.Date)
	})
	t.Run("status+date", func(t *testing.T) {
		s := parseStatus("WAITING FOR CUSTOMER - Jan 5, 2020")
		assert.Equal(t, "WAITING FOR CUSTOMER", s.Name)
		assert.Equal(t, "", s.Comment)
		assert.Equal(t, time.Date(2020, 1, 5, 0, 0, 0, 0, time.Local), s.Date)
	})
	t.Run("bad-date", func(t *testing.T) {
		s := parseStatus("WAITING FOR CUSTOMER - waiting to hear from client multi-hyphen-word - Jan 335, 2020")
		assert.Equal(t, "WAITING FOR CUSTOMER", s.Name)
		assert.Equal(t, "waiting to hear from client multi-hyphen-word - Jan 335, 2020", s.Comment)
	})
	t.Run("comment-brackets", func(t *testing.T) {
		p := NewParser(strings.NewReader("JIRAPROJECT-123 - [Client X] - Can't frobnicate the blips  [STALE - Jun 10, 2020]\n"))
		todo := p.parseTodo()
		assert.Equal(t, "JIRAPROJECT-123", todo.Name)
		assert.Equal(t, "[Client X] - Can't frobnicate the blips", todo.Description)
		assert.Equal(t, "STALE", todo.Status.Name)
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
