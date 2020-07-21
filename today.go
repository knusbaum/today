package main

import (
	"fmt"
	"sort"
	"time"
)

type today struct {
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
func (t *today) updateStartup() {
	for i, item := range t.startup {
		item.number = i + 1
	}
}

func (t *today) Update() {
	t.updateStartup()

	for _, todo := range t.todos {
		if todo.jira == "" {
			todo.jira = fmt.Sprintf("TODO-%d", t.nextTodo)
			t.nextTodo++
		}
		if todo.status.date.IsZero() {
			todo.status.date = time.Now()
			timestr := time.Now().Format("3:04")
			if todo.status.comment != "" {
				t.log = append(t.log, fmt.Sprintf("%s - Moved %s (%s) to %s (%s)", timestr, todo.jira, todo.description, todo.status.name, todo.status.comment))
			} else {
				t.log = append(t.log, fmt.Sprintf("%s - Moved %s (%s) to  %s", timestr, todo.jira, todo.description, todo.status.name))
			}
		}
		if todo.status.name == "" {
			todo.status.name = "?"
		}
	}
}

type Priority []*todo

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

func (a Priority) Len() int      { return len(a) }
func (a Priority) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a Priority) Less(i, j int) bool {
	pi := priority(&a[i].status)
	pj := priority(&a[j].status)
	if pi == pj {
		return a[i].status.date.Before(a[j].status.date)
	} else {
		return pi < pj
	}
}

func (t *today) Sort() {
	if len(t.todos) == 0 {
		return
	}
	sort.Stable(Priority(t.todos))
}

// TODO: Test Clear function
func (t *today) Clear() {
	k := 0
	for i := 0; i < (len(t.todos) - 1); {
		if t.todos[i].status.name != "DONE" {
			t.todos[k] = t.todos[i]
			k++
		}
		i++
	}
	t.todos = t.todos[:k]
	for _, item := range t.startup {
		item.status = status{}
	}
	t.log = make([]string, 0)
}
