package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestClear(t *testing.T) {
	t.Run("done", func(t *testing.T) {
		now := time.Now()
		today := &today{
			todos: []*todo{
				&todo{description: "should not exist", status: status{name: "DONE", date: now}},
				&todo{description: "task 0", status: status{name: "?", date: now}},
				&todo{description: "task 1", status: status{name: "SOMEUNKNOWNSTATUS", date: now}},
				&todo{description: "task 2", status: status{name: "IN PROGRESS", date: now}},
				&todo{description: "should not exist", status: status{name: "DONE", date: now}},
			},
		}
		today.Clear()
		assert.Len(t, today.todos, 3)
		for i := 0; i < len(today.todos); i++ {
			assert.Equal(t, fmt.Sprintf("task %d", i), today.todos[i].description)
		}
	})
	t.Run("same", func(t *testing.T) {
		now := time.Now()
		today := &today{
			todos: []*todo{
				&todo{description: "task 0", status: status{name: "", date: now}},
				&todo{description: "task 1", status: status{name: "?", date: now}},
				&todo{description: "task 2", status: status{name: "SOMEUNKNOWNSTATUS", date: now}},
				&todo{description: "task 3", status: status{name: "IN PROGRESS", date: now}},
				&todo{description: "task 4", status: status{name: "STALE", date: now}},
			},
		}
		today.Clear()
		assert.Len(t, today.todos, 5)
		for i := 0; i < len(today.todos); i++ {
			assert.Equal(t, fmt.Sprintf("task %d", i), today.todos[i].description)
		}
	})
}

func TestSortTodos(t *testing.T) {
	t.Run("priority", func(t *testing.T) {
		now := time.Now()
		today := &today{
			todos: []*todo{
				&todo{description: "task 3", status: status{name: "IN PROGRESS", date: now}},
				&todo{description: "task 9", status: status{name: "STALE", date: now}},
				&todo{description: "task 7", status: status{name: "WAITING", date: now}},
				&todo{description: "task 10", status: status{name: "HOLD", date: now.Add(48 * time.Hour)}},
				&todo{description: "task 0", status: status{date: now}},
				&todo{description: "task 5", status: status{name: "READY", date: now}},
				&todo{description: "task 8", status: status{name: "RESPONDED", date: now}},
				&todo{description: "task 1", status: status{name: "?", date: now}},
				&todo{description: "task 4", status: status{name: "IN PROGRESS", date: now}},
				&todo{description: "task 2", status: status{name: "SOMEUNKNOWNSTATUS", date: now}},
				&todo{description: "task 11", status: status{name: "DONE", date: now}},
				&todo{description: "task 6", status: status{name: "REVIEW", date: now}},
			},
		}
		today.Sort()
		for i := 0; i < len(today.todos); i++ {
			assert.Equal(t, fmt.Sprintf("task %d", i), today.todos[i].description)
		}
	})

	t.Run("date", func(t *testing.T) {
		now := time.Now()
		today := &today{
			todos: []*todo{
				&todo{description: "task 14", status: status{name: "DONE", date: now}},
				&todo{description: "task 13", status: status{name: "HOLD", date: now.Add(48 * time.Hour)}},
				&todo{description: "task 12", status: status{name: "STALE", date: now}},
				&todo{description: "task 11", status: status{name: "RESPONDED", date: now.Add(24 * time.Hour)}},
				&todo{description: "task 10", status: status{name: "WAITING", date: now}},
				&todo{description: "task 9", status: status{name: "REVIEW", date: now.Add(24 * time.Hour)}},
				&todo{description: "task 8", status: status{name: "REVIEW", date: now}},
				&todo{description: "task 7", status: status{name: "READY", date: now}},
				&todo{description: "task 6", status: status{name: "READY", date: now.Add(-24 * time.Hour)}},
				&todo{description: "task 5", status: status{name: "IN PROGRESS", date: now}},
				&todo{description: "task 4", status: status{name: "IN PROGRESS", date: now.Add(-24 * time.Hour)}},
				&todo{description: "task 3", status: status{name: "IN PROGRESS", date: now.Add(-48 * time.Hour)}},
				&todo{description: "task 2", status: status{name: "SOMEUNKNOWNSTATUS", date: now}},
				&todo{description: "task 1", status: status{name: "?", date: now.Add(-24 * time.Hour)}},
				&todo{description: "task 0", status: status{date: now.Add(-48 * time.Hour)}},
			},
		}
		// Unlike the other tests, which depend on the sort being stable (which is by design), this set of todos has an absolute order, so
		// we can shuffle the list before we sort.
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(today.todos), func(i, j int) { today.todos[i], today.todos[j] = today.todos[j], today.todos[i] })
		today.Sort()
		for i := 0; i < len(today.todos); i++ {
			assert.Equal(t, fmt.Sprintf("task %d", i), today.todos[i].description)
		}
	})

	t.Run("hold", func(t *testing.T) {
		now := time.Now()
		today := &today{
			todos: []*todo{
				&todo{description: "task 4", status: status{name: "IN PROGRESS", date: now}},
				&todo{description: "task 10", status: status{name: "STALE", date: now}},
				&todo{description: "task 8", status: status{name: "WAITING", date: now}},
				&todo{description: "task 0", status: status{name: "HOLD", date: now}},
				&todo{description: "task 1", status: status{date: now}},
				&todo{description: "task 6", status: status{name: "READY", date: now}},
				&todo{description: "task 9", status: status{name: "RESPONDED", date: now}},
				&todo{description: "task 2", status: status{name: "?", date: now}},
				&todo{description: "task 5", status: status{name: "IN PROGRESS", date: now}},
				&todo{description: "task 3", status: status{name: "SOMEUNKNOWNSTATUS", date: now}},
				&todo{description: "task 11", status: status{name: "DONE", date: now}},
				&todo{description: "task 7", status: status{name: "REVIEW", date: now}},
			},
		}
		today.Sort()
		for i := 0; i < len(today.todos); i++ {
			assert.Equal(t, fmt.Sprintf("task %d", i), today.todos[i].description)
		}
	})

	t.Run("waiting", func(t *testing.T) {
		now := time.Now()
		today := &today{
			todos: []*todo{
				// 0, 1, 2 have same priority as 3, 4, 5 but always come first because their date is earlier.
				&todo{description: "task 10", status: status{name: "HOLD", date: now.Add(48 * time.Hour)}},
				&todo{description: "task 6", status: status{name: "IN PROGRESS", date: now}},
				&todo{description: "task 7", status: status{name: "IN PROGRESS", date: now}},
				&todo{description: "task 3", status: status{date: now}},
				&todo{description: "task 4", status: status{name: "?", date: now}},
				&todo{description: "task 5", status: status{name: "SOMEUNKNOWNSTATUS", date: now}},
				&todo{description: "task 0", status: status{name: "REVIEW", date: now.Add(-24 * time.Hour)}},
				&todo{description: "task 9", status: status{name: "STALE", date: now}},
				&todo{description: "task 1", status: status{name: "WAITING", date: now.Add(-24 * time.Hour)}},
				&todo{description: "task 11", status: status{name: "DONE", date: now}},
				&todo{description: "task 8", status: status{name: "READY", date: now}},
				&todo{description: "task 2", status: status{name: "RESPONDED", date: now.Add(-24 * time.Hour)}},
			},
		}
		today.Sort()
		for i := 0; i < len(today.todos); i++ {
			assert.Equal(t, fmt.Sprintf("task %d", i), today.todos[i].description)
		}
	})

	t.Run("stale", func(t *testing.T) {
		now := time.Now()
		today := &today{
			todos: []*todo{
				&todo{description: "task 0", status: status{name: "STALE", date: now.Add(-24 * 7 * time.Hour)}},
				&todo{description: "task 7", status: status{name: "REVIEW", date: now}},
				&todo{description: "task 1", status: status{date: now}},
				&todo{description: "task 11", status: status{name: "DONE", date: now}},
				&todo{description: "task 4", status: status{name: "IN PROGRESS", date: now}},
				&todo{description: "task 5", status: status{name: "IN PROGRESS", date: now}},
				&todo{description: "task 8", status: status{name: "WAITING", date: now}},
				&todo{description: "task 2", status: status{name: "?", date: now}},
				&todo{description: "task 10", status: status{name: "HOLD", date: now.Add(48 * time.Hour)}},
				&todo{description: "task 3", status: status{name: "SOMEUNKNOWNSTATUS", date: now}},
				&todo{description: "task 6", status: status{name: "READY", date: now}},
				&todo{description: "task 9", status: status{name: "RESPONDED", date: now}},
			},
		}
		today.Sort()
		for i := 0; i < len(today.todos); i++ {
			assert.Equal(t, fmt.Sprintf("task %d", i), today.todos[i].description)
		}
	})
}

func TestUpdateList(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		today := &today{
			startup: []*listItem{
				&listItem{number: 0, description: "item 0"},
				&listItem{number: 0, description: "item 1"},
				&listItem{number: 0, description: "item 2"},
				&listItem{number: 0, description: "item 3"},
				&listItem{number: 0, description: "item 4"},
				&listItem{number: 0, description: "item 5"},
			},
		}
		today.Update()
		for i := 0; i < len(today.startup); i++ {
			assert.Equal(t, i+1, today.startup[i].number)
			assert.Equal(t, fmt.Sprintf("item %d", i), today.startup[i].description)
		}
	})

	t.Run("addItem", func(t *testing.T) {
		today := &today{
			startup: []*listItem{
				&listItem{number: 1, description: "item 0"},
				&listItem{number: 2, description: "item 1"},
				&listItem{number: 0, description: "item 2"},
				&listItem{number: 3, description: "item 3"},
				&listItem{number: 4, description: "item 4"},
				&listItem{number: 5, description: "item 5"},
			},
		}
		today.Update()
		for i := 0; i < len(today.startup); i++ {
			assert.Equal(t, i+1, today.startup[i].number)
			assert.Equal(t, fmt.Sprintf("item %d", i), today.startup[i].description)
		}
	})
}

func TestTransformation(t *testing.T) {
	todayText := `Morning Start Up:

Do something

Do another thing

One more thing.

Notes:

Some note

Another Note

One More Note

Log:

TODO:
Some Task
Another Task [IN PROGRESS]



Something else
`

	r := strings.NewReader(todayText)
	p := NewParser(r)
	td, err := p.parseToday()
	if !assert.NoError(t, err) {
		return
	}
	td.Update()
	td.Sort()
	var b strings.Builder
	w := bufio.NewWriter(&b)
	writeToday(td, w)
	result := b.String()

	expected := `Morning Start Up:
1. Do something 
2. Do another thing 
3. One more thing. 

Notes:
Some note
Another Note
One More Note

Log:
[0-9]+:[0-9]{2} - Moved TODO-1 \(Another Task\) to  IN PROGRESS

TODO:
TODO-0 - Some Task \[\? - [A-Za-z]{3} [0-9]+, [0-9]{4}\]
TODO-2 - Something else \[\? - [A-Za-z]{3} [0-9]+, [0-9]{4}\]
TODO-1 - Another Task \[IN PROGRESS - [A-Za-z]{3} [0-9]+, [0-9]{4}\]
`

	assert.Regexp(t, expected, result)

}
