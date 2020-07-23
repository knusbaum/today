package today

import (
	"fmt"
	"sort"
	"time"
)

type Today struct {
	startup  []*listItem
	notes    []string
	log      []string
	todos    []*todo
	nextTodo int
}

type status struct {
	name    string
	comment string
	date    time.Time
}

type todo struct {
	jira        string
	description string
	status      status
	comments    []string
	blankBelow  bool
}

type listItem struct {
	number      int
	description string
	status      status
}

// updateStartup makes sure that every startup item is numbered according
// to it's place in the startup list.
func (t *Today) updateStartup() {
	for i, item := range t.startup {
		item.number = i + 1
	}
}

// updateTodos adds dates and statuses to any todos without them.
func (t *Today) updateTodos() {
	for _, todo := range t.todos {
		if todo.jira == "" {
			todo.jira = fmt.Sprintf("TODO-%d", t.nextTodo)
			t.nextTodo++
		}
		if todo.status.date.IsZero() {
			todo.status.date = time.Now()
			timestr := time.Now().Format("3:04")
			if !todo.status.isUnknown() {
				if todo.status.comment != "" {
					t.log = append(t.log, fmt.Sprintf("%s - Moved %s (%s) to %s (%s)", timestr, todo.jira, todo.description, todo.status.name, todo.status.comment))
				} else {
					t.log = append(t.log, fmt.Sprintf("%s - Moved %s (%s) to  %s", timestr, todo.jira, todo.description, todo.status.name))
				}
			}
		}
		if todo.status.name == "" {
			todo.status.name = "?"
		}
	}
}

func (t *Today) Update() {
	t.updateStartup()
	t.updateTodos()
}

var priorityOrder map[string]int = map[string]int{
	"":            0,
	"?":           0,
	"OTHER":       0, // Other (not in this list)
	"IN PROGRESS": 1,
	"IN-PROGRESS": 1,
	"INPROGRESS":  1,
	"READY":       2,
	"REVIEW":      4,
	"WAITING":     5,
	"RESPONDED":   5,
	"STALE":       6,
	"HOLD":        7,
	"DONE":        8,
}

func (s *status) isUnknown() bool {
	return s.name == "" || s.name == "?"
}

func priority(s *status) int {
	// Tasks marked "HOLD" should be held until the specified date.
	if s.name == "HOLD" && time.Now().After(s.date) {
		return 0
	}

	// Tasks waiting or in review should be checked again after a day.
	if (s.name == "WAITING" || s.name == "REVIEW" || s.name == "RESPONDED") &&
		time.Now().After(s.date.Add(24*time.Hour)) {
		return 0
	}

	// Check up on stale tasks once per week
	if s.name == "STALE" && time.Now().After(s.date.Add(24*7*time.Hour)) {
		return 0
	}

	v, ok := priorityOrder[s.name]
	if !ok {
		return priorityOrder["OTHER"]
	}
	return v
}

type byPriority []*todo

func (a byPriority) Len() int      { return len(a) }
func (a byPriority) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byPriority) Less(i, j int) bool {
	pi := priority(&a[i].status)
	pj := priority(&a[j].status)
	if pi == pj {
		return a[i].status.date.Before(a[j].status.date)
	}
	return pi < pj
}

func (t *Today) Sort() {
	if len(t.todos) == 0 {
		return
	}
	sort.Stable(byPriority(t.todos))
}

func (t *Today) Clear() {
	k := 0
	for i := 0; i < len(t.todos); {
		if t.todos[i].status.name != "DONE" {
			t.todos[k] = t.todos[i]
			k++
		}
		i++
	}
	t.todos = t.todos[:k]

	// Eliminate status from startup items
	for _, item := range t.startup {
		item.status = status{}
	}
	t.log = make([]string, 0)
}
