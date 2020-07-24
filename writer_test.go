package today

import (
	"bufio"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var testStartup []*ListItem = []*ListItem{
	&ListItem{number: 1, Description: "Catch up on slack"},
	&ListItem{number: 2, Description: "Check the calendar"},
	&ListItem{number: 3, Description: "Read the inbox"},
	&ListItem{number: 4, Description: "look at JIRAPROJECT"},
}

func TestWriter(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		today := &Today{Startup: testStartup}
		var b strings.Builder
		w := bufio.NewWriter(&b)
		today.Write(w)
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
		today := &Today{
			Startup: testStartup,
			Notes: []string{
				"foop boop doop This is a note.",
				"note important_facts ",
				"note testing",
				"note deploy",
				"note custom_build",
			},
			Log: []string{
				"8:30 Starting work",
				"9:00 Standup",
				"9:15 Starting on Jira task",
			},
			Tasks: TaskList{
				Tasks: []*Task{
					&Task{
						Name:        "SOMEJIRA-123",
						Description: "description of a todo task",
						Status: Status{
							Name:    "IN PROGRESS",
							Comment: "waiting for customer",
							Date:    time.Date(2020, 6, 7, 0, 0, 0, 0, time.UTC),
						},
						Comments: []string{
							"* Some note",
							"* Some other note",
						},
					},
					&Task{
						Description: "Some other random task",
					},
				},
			},
		}
		var b strings.Builder
		w := bufio.NewWriter(&b)
		today.Write(w)
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
	today := &Today{
		Startup: testStartup,
		Notes: []string{
			"foop boop doop This is a note.",
			"note important_facts ",
			"note testing",
			"note deploy",
			"note custom_build",
		},
		Log: []string{
			"8:30 Starting work",
			"9:00 Standup",
			"9:15 Starting on Jira task",
		},
		Tasks: TaskList{
			Tasks: []*Task{
				&Task{
					Name:        "SOMEJIRA-123",
					Description: "description of a todo task",
					Status: Status{
						Name:    "IN PROGRESS",
						Comment: "waiting for customer",
						Date:    time.Date(2020, 6, 7, 0, 0, 0, 0, time.Local),
					},
					Comments: []string{
						"* Some note",
						"* Some other note",
					},
				},
				&Task{
					Description: "Some other random task",
				},
			},
		},
	}
	var b strings.Builder
	w := bufio.NewWriter(&b)
	today.Write(w)

	r := strings.NewReader(b.String())
	newToday, err := Parse(r)

	assert.NoError(t, err)
	if !assert.NotNil(t, newToday) {
		return
	}

	assert.Equal(t, today, newToday)
}
