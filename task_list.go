package today

import (
	"fmt"
	"sort"
	"time"
)

// TaskList represents a list of Tasks.
type TaskList struct {
	Tasks      []*Task
	nextTaskID int
}

// A Task is a structure representing a task.
// A Task in a today file may look like this:
//   TASK-1 - Description [STATUS NAME - Status.Comment - Status.Date]
//   	Comment 1
//   	Comment 2
// Name must be an all-caps sequence followed by a hyphen and a number. In other words, it must
// match the regex: `[A-Z]+-[0-9]+`
//
// In its most basic form, a task is a single line of text describing the task:
//   Do something important
//
// When such a task is parsed, it is given a name like "TASK-1", and an empty
// Status with the current day's date:
//   TASK-1 - Do something important [? - Jul 22, 2020]
//
// You can choose to provide your own name, followed by a hyphen, as long as it
// matches the regex `[A-Z]+-[0-9]+`. This allows you to use jira ids as task
// names. This may change in the future, as jira is by no means the only task
// tracking system, just the one I regularly use.
//
// A Task may have any number of comments beneath it. A comment is a line that
// begins with a tab (`\t`) character. Blank lines are allowed between tasks and
// their comments, but will be eliminated when writing a TaskList:
//   Do something important
//   	step 1
//
//   	step 2
//
//   	step 3
// will become:
//   TASK-1 - Do something important [? - Jul 22, 2020]
//   	step 1
//   	step 2
//   	step 3
type Task struct {
	Name        string
	Description string
	Status      Status
	Comments    []string
	blankBelow  bool
}

// Update adds dates and statuses to any todos without them. If log is not nil, it will add entries
// to the log whenever it adds a date to a task's status.
func (t *TaskList) Update(log *Lines) {
	for _, todo := range t.Tasks {
		if todo.Name == "" {
			todo.Name = fmt.Sprintf("TASK-%d", t.nextTaskID)
			t.nextTaskID++
		}
		if todo.Status.Date.IsZero() {
			todo.Status.Date = time.Now()
			timestr := time.Now().Format("3:04")
			if log != nil && !todo.Status.isUnknown() {
				if todo.Status.Comment != "" {
					log.Add(fmt.Sprintf("%s - Moved %s (%s) to %s (%s)", timestr, todo.Name, todo.Description, todo.Status.Name, todo.Status.Comment))
				} else {
					log.Add(fmt.Sprintf("%s - Moved %s (%s) to  %s", timestr, todo.Name, todo.Description, todo.Status.Name))
				}
			}
		}
		if todo.Status.Name == "" {
			todo.Status.Name = "?"
		}
	}
}

// Sort sorts a TaskList according to the statuses of the tasks. The goal is to always have a
// TaskList sorted by priority. Tasks are sorted by Status name in the following order:
//   <blank>
//   "?"
//   <other statuses not in this list>
//   "IN PROGRESS"
//   "IN-PROGRESS"
//   "INPROGRESS"
//   "READY"
//   "REVIEW"
//   "WAITING"
//   "RESPONDED"
//  "STALE"
//   "HOLD"
//   "DONE"
//
// Tasks with the same Status are sorted by date, oldest first.
//
// Tasks with no status or unknown status are first, the idea being they should be given one of the
// existing, known statuses.
//
// "IN PROGRESS" and variations of that follow. In progress tasks should always be at the top of my
// list.
//
// "READY" tasks are next. When I'm through working on an "IN PROGRESS" task, I pick up a new
// "READY" task.
//
// Then are "REVIEW", "WAITING", "RESPONDED", and "STALE" tasks. These are tasks which cannot be
// actively worked on, but which I am tracking. The sorting rules push these to the top
// periodically to make sure I don't forget about them (See the sorting exceptions below)
//
// "HOLD" tasks are next. They are a notable exception to the normal status format since tasks in
// "HOLD" status have a date in the future. Once that date arrives, "HOLD" tasks are moved to the
// top of the list. This allows me to put off a task until a specific date.
//
// "DONE" tasks are last. They should include the final status for the task, and are put at the
// bottom to keep a record of how and when a task was completed. They are cleared by Clear()
func (t *TaskList) Sort() {
	if len(t.Tasks) == 0 {
		return
	}
	sort.Stable(byPriority(t.Tasks))
}

// Clear removes all items with Status.Name == "DONE" from the TaskList
func (t *TaskList) Clear() {
	k := 0
	for i := 0; i < len(t.Tasks); {
		if t.Tasks[i].Status.Name != "DONE" {
			t.Tasks[k] = t.Tasks[i]
			k++
		}
		i++
	}
	t.Tasks = t.Tasks[:k]
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

func (s *Status) isUnknown() bool {
	return s.Name == "" || s.Name == "?"
}

func priority(s *Status) int {
	// Tasks marked "HOLD" should be held until the specified date.
	if s.Name == "HOLD" && time.Now().After(s.Date) {
		return 0
	}

	// Tasks waiting or in review should be checked again after a day.
	if (s.Name == "WAITING" || s.Name == "REVIEW" || s.Name == "RESPONDED") &&
		time.Now().After(s.Date.Add(24*time.Hour)) {
		return 0
	}

	// Check up on stale tasks once per week
	if s.Name == "STALE" && time.Now().After(s.Date.Add(24*7*time.Hour)) {
		return 0
	}

	v, ok := priorityOrder[s.Name]
	if !ok {
		return priorityOrder["OTHER"]
	}
	return v
}

type byPriority []*Task

func (a byPriority) Len() int      { return len(a) }
func (a byPriority) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byPriority) Less(i, j int) bool {
	pi := priority(&a[i].Status)
	pj := priority(&a[j].Status)
	if pi == pj {
		return a[i].Status.Date.Before(a[j].Status.Date)
	}
	return pi < pj
}
