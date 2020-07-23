package today

import (
	"time"
)

// Today represents a today file. A today file contains several
// sections to help track tasks.
//
// A today file contains 4 sections in this order:
//   1. Startup ("Morning Start Up")
//   2. Notes ("Notes")
//   3. Log ("Log")
//   4. Tasks ("TODO")
//
// Each section behaves slightly differently.
//
// The "Startup" section contains a list of items that are meant to be undertaken daily,
// beginning the day. For instance, a morning start up may contain:
//   1. Catch up on slack
//   2. Check the calendar
//   3. Read the inbox
//   4. look at issue tracker
//
// The "Notes" section is a simple sequence of lines that carries over from day to day. It is not
// emptied by a call to Clear(). Whitespace is eliminated (which may be changed in the future) but
// no other transformations apply. This is where I dump random notes, shell commands, etc. that I
// use frequently or want to remember. The whitespace elimination is to try to force notes to be
// short (single lines). For longer notes, I add the filename of a note file instead.
//
// The "Log" section is a sequence of lines similar to "Notes", but this one is cleared by a call
// to Clear(). Every time Update() applies a Status to a task in the "Tasks" section, a line is
// added to the "Log" section with a timestamp and a description of the applied status. For
// example, when Update() notices a task with a new status without a date (e.g. "TASK-123 - Do
// something important [DONE - Finished up]"), It Makes an entry in the log like so:
//   4:48 - Moved TASK-123 (Do something important) to DONE (Finished up)
// I also add my own timestamped entries to the Log when I want to record some important event.
//
// The "Tasks" section is the most complicated section. It is a sequence of tasks that have
// Statuses and optional comments. See TaskList for details.
type Today struct {
	Startup List
	Notes   Lines
	Log     Lines
	Tasks   TaskList
}

// Status describes the current status of a ListItem or Task. A status has a Name, which should be
// an string of capital letters and spaces. Each status also has a Comment field, and a Date.
//
// In its most basic text form, a Status is just a status name inside square brackets. A status
// name is a series of capitalized letters. A basic status might be "[IN PROGRESS]".
//
// Comments come after the status name, separated by space and a hyphen, and can be any string of
// text: "[IN PROGRESS - Working on pr #12]".
//
// Dates follow the status name and the comment if one is present, again separated by space and
// hyphen. Dates must follow the format (given in Go's time.Format
// (https://golang.org/pkg/time/#Time.Format) notation) "Jan _2, 2006" or they will be considered
// part of the comment. Both of these are valid statuses with dates:
//   [IN PROGRESS - Working on pr #12 - Jan 16, 2020]
//   [READY - Jan 14, 2020]
type Status struct {
	Name    string
	Comment string
	Date    time.Time
}

// ListItem represents one line of a List. Lists are ordered and have a Description and optional Status.
type ListItem struct {
	number      int
	Description string
	Status      Status
}

// Lines is a Section containing an arbitrary sequence of lines. When parsing, whitespace is
// eliminated (which may be changed in the future) but no other transformations apply. The
// whitespace elimination is to try to force entries in these sections to be short (single lines).
//
// Lines is used for both the Notes and Log sections, although each of those sections has slightly
// different behavior.
type Lines []string

// Add adds a line to a Lines. s should not contain newline characters.
func (l *Lines) Add(s string) {
	*l = append(*l, s)
}

// List is a Section containing an ordered list of items. Each ListItem represents one line of a
// list. ListItems are given numbers from 1 to len(list). A ListItem's number must be the first
// thing on the line. It is some number of digits followed by a period.
type List []*ListItem

// Update makes sure that every startup item is numbered according to it's place in the startup
// list. It will assign numbers to mis-numbered and un-numbered items. For instance, a list
// containing:
//   3. Check the calendar
//   Catch up on slack
//   Read the inbox
//   look at ticket tracker
// will be re-numbered to:
//   1. Check the calendar
//   2. Catch up on slack
//   3. Read the inbox
//   4. look at ticket tracker
func (l List) Update() {
	for i, item := range l {
		item.number = i + 1
	}
}

// Update makes sure items in Startup are numbered correctly, and applies statuses to un-statused
// items in the Tasks section. (See TaskList.Update and List.Update)
func (t *Today) Update() {
	t.Startup.Update()
	t.Tasks.Update(&t.Log)
}

// Sort sorts the Tasks section (See Tasks.Sort)
func (t *Today) Sort() {
	t.Tasks.Sort()
}

// Clear clears statuses from the Startup section, and eliminates "DONE" tasks from the Tasks
// section. (See Tasks.Clear)
func (t *Today) Clear() {
	t.Tasks.Clear()

	// Eliminate status from startup items
	for _, item := range t.Startup {
		item.Status = Status{}
	}
	t.Log = make([]string, 0)
}
