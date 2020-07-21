package main

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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
