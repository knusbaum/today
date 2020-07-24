package today

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
		today := &Today{
			Tasks: TaskList{
				Tasks: []*Task{
					&Task{Description: "should not exist", Status: Status{Name: "DONE", Date: now}},
					&Task{Description: "task 0", Status: Status{Name: "?", Date: now}},
					&Task{Description: "task 1", Status: Status{Name: "SOMEUNKNOWNSTATUS", Date: now}},
					&Task{Description: "task 2", Status: Status{Name: "IN PROGRESS", Date: now}},
					&Task{Description: "should not exist", Status: Status{Name: "DONE", Date: now}},
				},
			},
		}
		today.Clear()
		assert.Len(t, today.Tasks.Tasks, 3)
		for i := 0; i < len(today.Tasks.Tasks); i++ {
			assert.Equal(t, fmt.Sprintf("task %d", i), today.Tasks.Tasks[i].Description)
		}
	})
	t.Run("same", func(t *testing.T) {
		now := time.Now()
		today := &Today{
			Tasks: TaskList{
				Tasks: []*Task{
					&Task{Description: "task 0", Status: Status{Name: "", Date: now}},
					&Task{Description: "task 1", Status: Status{Name: "?", Date: now}},
					&Task{Description: "task 2", Status: Status{Name: "SOMEUNKNOWNSTATUS", Date: now}},
					&Task{Description: "task 3", Status: Status{Name: "IN PROGRESS", Date: now}},
					&Task{Description: "task 4", Status: Status{Name: "STALE", Date: now}},
				},
			},
		}
		today.Clear()
		assert.Len(t, today.Tasks.Tasks, 5)
		for i := 0; i < len(today.Tasks.Tasks); i++ {
			assert.Equal(t, fmt.Sprintf("task %d", i), today.Tasks.Tasks[i].Description)
		}
	})
}

func TestSortTodos(t *testing.T) {
	t.Run("priority", func(t *testing.T) {
		now := time.Now()
		today := &Today{
			Tasks: TaskList{
				Tasks: []*Task{
					&Task{Description: "task 3", Status: Status{Name: "IN PROGRESS", Date: now}},
					&Task{Description: "task 9", Status: Status{Name: "STALE", Date: now}},
					&Task{Description: "task 7", Status: Status{Name: "WAITING", Date: now}},
					&Task{Description: "task 10", Status: Status{Name: "HOLD", Date: now.Add(48 * time.Hour)}},
					&Task{Description: "task 0", Status: Status{Date: now}},
					&Task{Description: "task 5", Status: Status{Name: "READY", Date: now}},
					&Task{Description: "task 8", Status: Status{Name: "RESPONDED", Date: now}},
					&Task{Description: "task 1", Status: Status{Name: "?", Date: now}},
					&Task{Description: "task 4", Status: Status{Name: "IN PROGRESS", Date: now}},
					&Task{Description: "task 2", Status: Status{Name: "SOMEUNKNOWNSTATUS", Date: now}},
					&Task{Description: "task 11", Status: Status{Name: "DONE", Date: now}},
					&Task{Description: "task 6", Status: Status{Name: "REVIEW", Date: now}},
				},
			},
		}
		today.Sort()
		for i := 0; i < len(today.Tasks.Tasks); i++ {
			assert.Equal(t, fmt.Sprintf("task %d", i), today.Tasks.Tasks[i].Description)
		}
	})

	t.Run("date", func(t *testing.T) {
		now := time.Now()
		today := &Today{
			Tasks: TaskList{
				Tasks: []*Task{
					&Task{Description: "task 14", Status: Status{Name: "DONE", Date: now}},
					&Task{Description: "task 13", Status: Status{Name: "HOLD", Date: now.Add(48 * time.Hour)}},
					&Task{Description: "task 12", Status: Status{Name: "STALE", Date: now}},
					&Task{Description: "task 11", Status: Status{Name: "RESPONDED", Date: now.Add(24 * time.Hour)}},
					&Task{Description: "task 10", Status: Status{Name: "WAITING", Date: now}},
					&Task{Description: "task 9", Status: Status{Name: "REVIEW", Date: now.Add(24 * time.Hour)}},
					&Task{Description: "task 8", Status: Status{Name: "REVIEW", Date: now}},
					&Task{Description: "task 7", Status: Status{Name: "READY", Date: now}},
					&Task{Description: "task 6", Status: Status{Name: "READY", Date: now.Add(-24 * time.Hour)}},
					&Task{Description: "task 5", Status: Status{Name: "IN PROGRESS", Date: now}},
					&Task{Description: "task 4", Status: Status{Name: "IN PROGRESS", Date: now.Add(-24 * time.Hour)}},
					&Task{Description: "task 3", Status: Status{Name: "IN PROGRESS", Date: now.Add(-48 * time.Hour)}},
					&Task{Description: "task 2", Status: Status{Name: "SOMEUNKNOWNSTATUS", Date: now}},
					&Task{Description: "task 1", Status: Status{Name: "?", Date: now.Add(-24 * time.Hour)}},
					&Task{Description: "task 0", Status: Status{Date: now.Add(-48 * time.Hour)}},
				},
			},
		}
		// Unlike the other tests, which depend on the sort being stable (which is by design), this set of todos has an absolute order, so
		// we can shuffle the list before we sort.
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(today.Tasks.Tasks), func(i, j int) {
			today.Tasks.Tasks[i], today.Tasks.Tasks[j] = today.Tasks.Tasks[j], today.Tasks.Tasks[i]
		})
		today.Sort()
		for i := 0; i < len(today.Tasks.Tasks); i++ {
			assert.Equal(t, fmt.Sprintf("task %d", i), today.Tasks.Tasks[i].Description)
		}
	})

	t.Run("hold", func(t *testing.T) {
		now := time.Now()
		today := &Today{
			Tasks: TaskList{
				Tasks: []*Task{
					&Task{Description: "task 4", Status: Status{Name: "IN PROGRESS", Date: now}},
					&Task{Description: "task 10", Status: Status{Name: "STALE", Date: now}},
					&Task{Description: "task 8", Status: Status{Name: "WAITING", Date: now}},
					&Task{Description: "task 0", Status: Status{Name: "HOLD", Date: now}},
					&Task{Description: "task 1", Status: Status{Date: now}},
					&Task{Description: "task 6", Status: Status{Name: "READY", Date: now}},
					&Task{Description: "task 9", Status: Status{Name: "RESPONDED", Date: now}},
					&Task{Description: "task 2", Status: Status{Name: "?", Date: now}},
					&Task{Description: "task 5", Status: Status{Name: "IN PROGRESS", Date: now}},
					&Task{Description: "task 3", Status: Status{Name: "SOMEUNKNOWNSTATUS", Date: now}},
					&Task{Description: "task 11", Status: Status{Name: "DONE", Date: now}},
					&Task{Description: "task 7", Status: Status{Name: "REVIEW", Date: now}},
				},
			},
		}
		today.Sort()
		for i := 0; i < len(today.Tasks.Tasks); i++ {
			assert.Equal(t, fmt.Sprintf("task %d", i), today.Tasks.Tasks[i].Description)
		}
	})

	t.Run("waiting", func(t *testing.T) {
		now := time.Now()
		today := &Today{
			Tasks: TaskList{
				Tasks: []*Task{
					// 0, 1, 2 have same priority as 3, 4, 5 but always come first because their date is earlier.
					&Task{Description: "task 10", Status: Status{Name: "HOLD", Date: now.Add(48 * time.Hour)}},
					&Task{Description: "task 6", Status: Status{Name: "IN PROGRESS", Date: now}},
					&Task{Description: "task 7", Status: Status{Name: "IN PROGRESS", Date: now}},
					&Task{Description: "task 3", Status: Status{Date: now}},
					&Task{Description: "task 4", Status: Status{Name: "?", Date: now}},
					&Task{Description: "task 5", Status: Status{Name: "SOMEUNKNOWNSTATUS", Date: now}},
					&Task{Description: "task 0", Status: Status{Name: "REVIEW", Date: now.Add(-24 * time.Hour)}},
					&Task{Description: "task 9", Status: Status{Name: "STALE", Date: now}},
					&Task{Description: "task 1", Status: Status{Name: "WAITING", Date: now.Add(-24 * time.Hour)}},
					&Task{Description: "task 11", Status: Status{Name: "DONE", Date: now}},
					&Task{Description: "task 8", Status: Status{Name: "READY", Date: now}},
					&Task{Description: "task 2", Status: Status{Name: "RESPONDED", Date: now.Add(-24 * time.Hour)}},
				},
			},
		}
		today.Sort()
		for i := 0; i < len(today.Tasks.Tasks); i++ {
			assert.Equal(t, fmt.Sprintf("task %d", i), today.Tasks.Tasks[i].Description)
		}
	})

	t.Run("stale", func(t *testing.T) {
		now := time.Now()
		today := &Today{
			Tasks: TaskList{
				Tasks: []*Task{
					&Task{Description: "task 0", Status: Status{Name: "STALE", Date: now.Add(-24 * 7 * time.Hour)}},
					&Task{Description: "task 7", Status: Status{Name: "REVIEW", Date: now}},
					&Task{Description: "task 1", Status: Status{Date: now}},
					&Task{Description: "task 11", Status: Status{Name: "DONE", Date: now}},
					&Task{Description: "task 4", Status: Status{Name: "IN PROGRESS", Date: now}},
					&Task{Description: "task 5", Status: Status{Name: "IN PROGRESS", Date: now}},
					&Task{Description: "task 8", Status: Status{Name: "WAITING", Date: now}},
					&Task{Description: "task 2", Status: Status{Name: "?", Date: now}},
					&Task{Description: "task 10", Status: Status{Name: "HOLD", Date: now.Add(48 * time.Hour)}},
					&Task{Description: "task 3", Status: Status{Name: "SOMEUNKNOWNSTATUS", Date: now}},
					&Task{Description: "task 6", Status: Status{Name: "READY", Date: now}},
					&Task{Description: "task 9", Status: Status{Name: "RESPONDED", Date: now}},
				},
			},
		}
		today.Sort()
		for i := 0; i < len(today.Tasks.Tasks); i++ {
			assert.Equal(t, fmt.Sprintf("task %d", i), today.Tasks.Tasks[i].Description)
		}
	})
}

func TestUpdateList(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		today := &Today{
			Startup: []*ListItem{
				&ListItem{number: 0, Description: "item 0"},
				&ListItem{number: 0, Description: "item 1"},
				&ListItem{number: 0, Description: "item 2"},
				&ListItem{number: 0, Description: "item 3"},
				&ListItem{number: 0, Description: "item 4"},
				&ListItem{number: 0, Description: "item 5"},
			},
		}
		today.Update()
		for i := 0; i < len(today.Startup); i++ {
			assert.Equal(t, i+1, today.Startup[i].number)
			assert.Equal(t, fmt.Sprintf("item %d", i), today.Startup[i].Description)
		}
	})

	t.Run("addItem", func(t *testing.T) {
		today := &Today{
			Startup: []*ListItem{
				&ListItem{number: 1, Description: "item 0"},
				&ListItem{number: 2, Description: "item 1"},
				&ListItem{number: 0, Description: "item 2"},
				&ListItem{number: 3, Description: "item 3"},
				&ListItem{number: 4, Description: "item 4"},
				&ListItem{number: 5, Description: "item 5"},
			},
		}
		today.Update()
		for i := 0; i < len(today.Startup); i++ {
			assert.Equal(t, i+1, today.Startup[i].number)
			assert.Equal(t, fmt.Sprintf("item %d", i), today.Startup[i].Description)
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
	td, err := Parse(r)
	if !assert.NoError(t, err) {
		return
	}
	td.Update()
	td.Sort()
	var b strings.Builder
	w := bufio.NewWriter(&b)
	td.Write(w)
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
[0-9]+:[0-9]{2} - Moved TASK-1 \(Another Task\) to  IN PROGRESS

TODO:
TASK-0 - Some Task \[\? - [A-Za-z]{3} [0-9]+, [0-9]{4}\]
TASK-2 - Something else \[\? - [A-Za-z]{3} [0-9]+, [0-9]{4}\]
TASK-1 - Another Task \[IN PROGRESS - [A-Za-z]{3} [0-9]+, [0-9]{4}\]
`

	assert.Regexp(t, expected, result)

}
