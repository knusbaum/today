package main

import (
	"bufio"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var testStartup []*listItem = []*listItem{
	&listItem{number: 1, description: "Catch up on slack"},
	&listItem{number: 2, description: "Check the calendar"},
	&listItem{number: 3, description: "Read the inbox"},
	&listItem{number: 4, description: "look at JIRAPROJECT"},
}

func TestWriter(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		today := &today{startup: testStartup}
		var b strings.Builder
		w := bufio.NewWriter(&b)
		writeToday(today, w)
		expect := `Morning Start Up:
1. Catch up on slack 
2. Check the calendar 
3. Read the inbox 
4. look at JIRAPROJECT 

Notes:

Log:

TODO:

`
		assert.Equal(t, expect, b.String())
	})

	t.Run("full", func(t *testing.T) {
		today := &today{
			startup: testStartup,
			notes: []string{
				"foop boop doop This is a note.",
				"note important_facts ",
				"note testing",
				"note deploy",
				"note custom_build",
			},
			log: []string{
				"8:30 Starting work",
				"9:00 Standup",
				"9:15 Starting on Jira task",
			},
			todos: []*todo{
				&todo{
					jira:        "SOMEJIRA-123",
					description: "description of a todo task",
					status: status{
						name:    "IN PROGRESS",
						comment: "waiting for customer",
						date:    time.Date(2020, 6, 7, 0, 0, 0, 0, time.UTC),
					},
					comments: []string{
						"* Some note",
						"* Some other note",
					},
				},
				&todo{
					description: "Some other random task",
				},
			},
		}
		var b strings.Builder
		w := bufio.NewWriter(&b)
		writeToday(today, w)
		expect := `Morning Start Up:
1. Catch up on slack 
2. Check the calendar 
3. Read the inbox 
4. look at JIRAPROJECT 

Notes:
foop boop doop This is a note.
note important_facts 
note testing
note deploy
note custom_build

Log:
8:30 Starting work
9:00 Standup
9:15 Starting on Jira task

TODO:
SOMEJIRA-123 - description of a todo task [IN PROGRESS - waiting for customer - Jun  7, 2020]
	* Some note
	* Some other note
Some other random task 

`
		assert.Equal(t, expect, b.String())
	})
}

func TestReadback(t *testing.T) {
	today := &today{
		startup: testStartup,
		notes: []string{
			"foop boop doop This is a note.",
			"note important_facts ",
			"note testing",
			"note deploy",
			"note custom_build",
		},
		log: []string{
			"8:30 Starting work",
			"9:00 Standup",
			"9:15 Starting on Jira task",
		},
		todos: []*todo{
			&todo{
				jira:        "SOMEJIRA-123",
				description: "description of a todo task",
				status: status{
					name:    "IN PROGRESS",
					comment: "waiting for customer",
					date:    time.Date(2020, 6, 7, 0, 0, 0, 0, time.Local),
				},
				comments: []string{
					"* Some note",
					"* Some other note",
				},
			},
			&todo{
				description: "Some other random task",
			},
		},
	}
	var b strings.Builder
	w := bufio.NewWriter(&b)
	writeToday(today, w)

	r := strings.NewReader(b.String())
	p := NewParser(r)
	newToday, err := p.parseToday()

	assert.NoError(t, err)
	if !assert.NotNil(t, newToday) {
		return
	}

	assert.Equal(t, today, newToday)
}
