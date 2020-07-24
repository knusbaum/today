# today
[![Semaphore Unit Tests](https://knusbaum.semaphoreci.com/badges/today.svg)](https://knusbaum.semaphoreci.com/branches/29f4d916-1283-46fd-a311-81f74182a4c2)

`import "github.com/knusbaum/today"`
* [Overview](#pkg-overview)
* [Index](#pkg-index)
* [Subdirectories](#pkg-subdirectories)
## <a name="pkg-overview">Overview</a>
Package today provides tools for manipulating a plain-text task list. This was developed to help
me specifically with my workflow. The package manages "today files" which contain several
sections. Each section has slightly different structure and different rules about how it is
updated, sorted, etc. The today package is primarily intended to be used by the today program in
package github.com/knusbaum/today/today.

## <a name="pkg-index">Index</a>

* [type Lines](#Lines)
  * [func (l *Lines) Add(s string)](#Lines.Add)
* [type List](#List)
  * [func (l List) Update()](#List.Update)
* [type ListItem](#ListItem)
* [type Status](#Status)
* [type Task](#Task)
* [type TaskList](#TaskList)
  * [func (t *TaskList) Clear()](#TaskList.Clear)
  * [func (t *TaskList) Sort()](#TaskList.Sort)
  * [func (t *TaskList) Update(log *Lines)](#TaskList.Update)
  * [func (t *TaskList) Write(w *bufio.Writer) error](#TaskList.Write)
* [type Today](#Today)
  * [func Parse(r io.Reader) (*Today, error)](#Parse)
  * [func (t *Today) Clear()](#Today.Clear)
  * [func (t *Today) Sort()](#Today.Sort)
  * [func (t *Today) Update()](#Today.Update)
  * [func (t *Today) Write(w io.Writer) error](#Today.Write)
#### <a name="pkg-files">Package files</a>
[doc.go](https://github.com/knusbaum/today/blob/master/doc.go) [parser.go](https://github.com/knusbaum/today/blob/master/parser.go) [task_list.go](https://github.com/knusbaum/today/blob/master/task_list.go) [today.go](https://github.com/knusbaum/today/blob/master/today.go) [writer.go](https://github.com/knusbaum/today/blob/master/writer.go) 
## <a name="Lines">type</a> [Lines](https://github.com/knusbaum/today/blob/master/today.go#L87)
```go
type Lines []string
```

Lines is a Section containing an arbitrary sequence of lines. When parsing, whitespace is
eliminated (which may be changed in the future) but no other transformations apply. The
whitespace elimination is to try to force entries in these sections to be short (single lines).

Lines is used for both the Notes and Log sections, although each of those sections has slightly
different behavior.

### <a name="Lines.Add">func</a> (\*Lines) [Add](https://github.com/knusbaum/today/blob/master/today.go#L90)
```go
func (l *Lines) Add(s string)
```

Add adds a line to a Lines. s should not contain newline characters.

## <a name="List">type</a> [List](https://github.com/knusbaum/today/blob/master/today.go#L97)
```go
type List []*ListItem
```

List is a Section containing an ordered list of items. Each ListItem represents one line of a
list. ListItems are given numbers from 1 to len(list). A ListItem's number must be the first
thing on the line. It is some number of digits followed by a period.

### <a name="List.Update">func</a> (List) [Update](https://github.com/knusbaum/today/blob/master/today.go#L111)
```go
func (l List) Update()
```

Update makes sure that every startup item is numbered according to it's place in the startup
list. It will assign numbers to mis-numbered and un-numbered items. For instance, a list
containing:

```
3. Check the calendar
Catch up on slack
Read the inbox
look at ticket tracker
```

will be re-numbered to:

```
1. Check the calendar
2. Catch up on slack
3. Read the inbox
4. look at ticket tracker
```

## <a name="ListItem">type</a> [ListItem](https://github.com/knusbaum/today/blob/master/today.go#L75)
```go
type ListItem struct {
    Description string
    Status      Status
    // contains filtered or unexported fields
}
```

ListItem represents one line of a List. Lists are ordered and have a Description and optional Status.

## <a name="Status">type</a> [Status](https://github.com/knusbaum/today/blob/master/today.go#L68)
```go
type Status struct {
    Name    string
    Comment string
    Date    time.Time
}
```

Status describes the current status of a ListItem or Task. A status has a Name, which should be
an string of capital letters and spaces. Each status also has a Comment field, and a Date.

In its most basic text form, a Status is just a status name inside square brackets. A status
name is a series of capitalized letters. A basic status might be "[IN PROGRESS]".

Comments come after the status name, separated by space and a hyphen, and can be any string of
text: "[IN PROGRESS - Working on pr #12]".

Dates follow the status name and the comment if one is present, again separated by space and
hyphen. Dates must follow the format (given in Go's time.Format
(<a href="https://golang.org/pkg/time/#Time.Format">https://golang.org/pkg/time/#Time.Format</a>) notation) "Jan _2, 2006" or they will be considered
part of the comment. Both of these are valid statuses with dates:

```
[IN PROGRESS - Working on pr #12 - Jan 16, 2020]
[READY - Jan 14, 2020]
```

## <a name="Task">type</a> [Task](https://github.com/knusbaum/today/blob/master/task_list.go#L49)
```go
type Task struct {
    Name        string
    Description string
    Status      Status
    Comments    []string
    // contains filtered or unexported fields
}
```

A Task is a structure representing a task.
A Task in a today file may look like this:

```
TASK-1 - Description [STATUS NAME - Status.Comment - Status.Date]
	Comment 1
	Comment 2
```

Name must be an all-caps sequence followed by a hyphen and a number. In other words, it must
match the regex: `[A-Z]+-[0-9]+`

In its most basic form, a task is a single line of text describing the task:

```
Do something important
```

When such a task is parsed, it is given a name like "TASK-1", and an empty
Status with the current day's date:

```
TASK-1 - Do something important [? - Jul 22, 2020]
```

You can choose to provide your own name, followed by a hyphen, as long as it
matches the regex `[A-Z]+-[0-9]+`. This allows you to use jira ids as task
names. This may change in the future, as jira is by no means the only task
tracking system, just the one I regularly use.

A Task may have any number of comments beneath it. A comment is a line that
begins with a tab (`\t`) character. Blank lines are allowed between tasks and
their comments, but will be eliminated when writing a TaskList:

```
Do something important
	step 1

	step 2

	step 3
```

will become:

```
TASK-1 - Do something important [? - Jul 22, 2020]
	step 1
	step 2
	step 3
```

## <a name="TaskList">type</a> [TaskList](https://github.com/knusbaum/today/blob/master/task_list.go#L10)
```go
type TaskList struct {
    Tasks []*Task
    // contains filtered or unexported fields
}
```

TaskList represents a list of Tasks.

### <a name="TaskList.Clear">func</a> (\*TaskList) [Clear](https://github.com/knusbaum/today/blob/master/task_list.go#L127)
```go
func (t *TaskList) Clear()
```

Clear removes all items with Status.Name == "DONE" from the TaskList

### <a name="TaskList.Sort">func</a> (\*TaskList) [Sort](https://github.com/knusbaum/today/blob/master/task_list.go#L119)
```go
func (t *TaskList) Sort()
```

Sort sorts a TaskList according to the statuses of the tasks. The goal is to always have a
TaskList sorted by priority. Tasks are sorted by Status name in the following order:

```
 <blank>
 "?"
 <other statuses not in this list>
 "IN PROGRESS"
 "IN-PROGRESS"
 "INPROGRESS"
 "READY"
 "REVIEW"
 "WAITING"
 "RESPONDED"
"STALE"
 "HOLD"
 "DONE"
```

Tasks with the same Status are sorted by date, oldest first.

Tasks with no status or unknown status are first, the idea being they should be given one of the
existing, known statuses.

"IN PROGRESS" and variations of that follow. In progress tasks should always be at the top of my
list.

"READY" tasks are next. When I'm through working on an "IN PROGRESS" task, I pick up a new
"READY" task.

Then are "REVIEW", "WAITING", "RESPONDED", and "STALE" tasks. These are tasks which cannot be
actively worked on, but which I am tracking. The sorting rules push these to the top
periodically to make sure I don't forget about them (See the sorting exceptions below)

"HOLD" tasks are next. They are a notable exception to the normal status format since tasks in
"HOLD" status have a date in the future. Once that date arrives, "HOLD" tasks are moved to the
top of the list. This allows me to put off a task until a specific date.

"DONE" tasks are last. They should include the final status for the task, and are put at the
bottom to keep a record of how and when a task was completed. They are cleared by Clear()

### <a name="TaskList.Update">func</a> (\*TaskList) [Update](https://github.com/knusbaum/today/blob/master/task_list.go#L59)
```go
func (t *TaskList) Update(log *Lines)
```

Update adds dates and statuses to any todos without them. If log is not nil, it will add entries
to the log whenever it adds a date to a task's status.

### <a name="TaskList.Write">func</a> (\*TaskList) [Write](https://github.com/knusbaum/today/blob/master/writer.go#L100)
```go
func (t *TaskList) Write(w *bufio.Writer) error
```

## <a name="Today">type</a> [Today](https://github.com/knusbaum/today/blob/master/today.go#L46)
```go
type Today struct {
    Startup List
    Notes   Lines
    Log     Lines
    Tasks   TaskList
}
```

Today represents a today file. A today file contains 4 sections in this order:

```
1. Startup ("Morning Start Up")
2. Notes ("Notes")
3. Log ("Log")
4. Tasks ("TODO")
```

Each section behaves slightly differently.

### Startup
The "Startup" section contains a list of items that are meant to be undertaken daily,
beginning the day. For instance, a morning start up may contain:

```
1. Catch up on slack
2. Check the calendar
3. Read the inbox
4. look at issue tracker
```

### Notes
The "Notes" section is a simple sequence of lines that carries over from day to day. It is not
emptied by a call to Clear(). Whitespace is eliminated (which may be changed in the future) but
no other transformations apply. This is where I dump random notes, shell commands, etc. that I
use frequently or want to remember. The whitespace elimination is to try to force notes to be
short (single lines). For longer notes, I add the filename of a note file instead.

### Log
The "Log" section is a sequence of lines similar to "Notes", but this one is cleared by a call
to Clear(). Every time Update() applies a Status to a task in the "Tasks" section, a line is
added to the "Log" section with a timestamp and a description of the applied status. For
example, when Update() notices a task with a new status without a date (e.g. "TASK-123 - Do
something important [DONE - Finished up]"), It Makes an entry in the log like so:

```
4:48 - Moved TASK-123 (Do something important) to DONE (Finished up)
```

I also add my own timestamped entries to the Log when I want to record some important event.

### Tasks
The "Tasks" section is the most complicated section. It is a sequence of tasks that have
Statuses and optional comments. See TaskList for details.

### <a name="Parse">func</a> [Parse](https://github.com/knusbaum/today/blob/master/parser.go#L301)
```go
func Parse(r io.Reader) (*Today, error)
```

Parse attempts to parse a *Today, from r. It returns an error if a *Today
could not be parsed.

### <a name="Today.Clear">func</a> (\*Today) [Clear](https://github.com/knusbaum/today/blob/master/today.go#L131)
```go
func (t *Today) Clear()
```

Clear clears statuses from the Startup section, and eliminates "DONE" tasks from the Tasks
section. (See TaskList.Clear)

### <a name="Today.Sort">func</a> (\*Today) [Sort](https://github.com/knusbaum/today/blob/master/today.go#L125)
```go
func (t *Today) Sort()
```

Sort sorts the Tasks section (See Tasks.Sort)

### <a name="Today.Update">func</a> (\*Today) [Update](https://github.com/knusbaum/today/blob/master/today.go#L119)
```go
func (t *Today) Update()
```

Update makes sure items in Startup are numbered correctly, and applies statuses to un-statused
items in the Tasks section. (See TaskList.Update and List.Update)

### <a name="Today.Write">func</a> (\*Today) [Write](https://github.com/knusbaum/today/blob/master/writer.go#L111)
```go
func (t *Today) Write(w io.Writer) error
```

Write writes a Today out to writer w in the normal form

- - -
Created: 24-Jul-2020 14:49:45 +0000
Generated by [godoc2md](http://github.com/thatgerber/godoc2md)
