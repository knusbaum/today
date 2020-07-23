# today

[![Semaphore Unit Tests](https://knusbaum.semaphoreci.com/badges/today.svg)](https://knusbaum.semaphoreci.com/branches/29f4d916-1283-46fd-a311-81f74182a4c2)

Today is a simple program for tracking daily tasks using a plain-text task
list. This was developed to help me specifically with my workflow, so is fairly
rigid in its structure.


## Today Files
The `today` program operates on "today files". A today file contains several
sections to help track tasks.

A today file contains 4 sections:
1. Morning Start Up
2. Notes
3. Log
4. TODO

Each section behaves slightly differently. 

### Morning Start Up
This section contains a list of items that are meant to be undertaken daily,
beginning the day. For instance, a morning start up may contain:
```
1. Catch up on slack
2. Check the calendar
3. Read the inbox
4. look at Jira 
```

Each item in `Morning Start Up` can be marked with a [`Status`](#status) which
will be cleared during [Generation](#generation) for the following day. This
helps you ensure you are completing your daily tasks, and identify tasks that
are getting left behind.

### Notes
Notes is a simple sequence of lines that carries over from day to day.
Whitespace is eliminated (which may be changed in the future) but no other
transformations apply. This is where I dump random notes, shell commands, etc.
that I use frequently or want to remember.

The whitespace elimination is to try to force notes to be short (single lines).
For longer notes, I add the filename of a note file instead.

### Log
Log is a sequence of lines similar to [Notes](#notes), but this one is cleared
during [Generation](#generation). Every time a [`Status`](#status) is applied
to a task in the [TODO](#todo) section, a line is added to the Log section with
a timestamp describing what status was applied.

For example, when `today` notices a task with a new status (e.g. `TASK-123 - Do
something important [DONE - Finished up]`), It Makes an entry in the log like
so:
```
4:48 - Moved TASK-123 (Do something important) to DONE (Finished up)
```

I also add my own timestamped entries to the Log when I want to record some
important event.

### TODO
TODO is the most complicated section. It is a sequence of tasks that have
[`Status`](#status)es and optional comments.


#### Task Structure
In its most basic form, a task is a single line of text describing the task:
```
Do something important
```

When `today` reads such a task, it gives it a name like `TASK-1`, and an empty
[`Status`](#status) with the current day's date:

```
TASK-1 - Do something important [? - Jul 22, 2020]
```

You can choose to provide your own name, followed by a hyphen, as long as it
matches the regex `[A-Z]+-[0-9]+`. This allows you to use jira ids as task
names. This may change in the future, as jira is by no means the only task
tracking system, just the one I regularly use.

#### Comments
A task may have any number of comments beneath it. A comment is a line that
begins with a tab (`\t`) character. Blank lines are allowed between tasks and
their comments, but will be eliminated by `today`.
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

#### Sorting
`today` sorts tasks by their [`Status`](#status) and date. The goal is to
always have a list of tasks sorted by priority. Tasks are sorted by
[`Status`](#status) name in the following order:

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

Tasks with the same [`Status`](#status) are sorted by date, oldest first.

Tasks with no status or unknown status are first, the idea being they should be
given one of the existing, known statuses.

`"IN PROGRESS"` and variations of that follow. In progress tasks should always be
at the top of my list.

`"READY"` tasks are next. When I'm through working on an `"IN PROGRESS"` task, I
pick up a new `"READY"` task.

Then are `"REVIEW"`, `"WAITING"`, `"RESPONDED"`, and `"STALE"` tasks. These are
tasks which cannot be actively worked on, but which I am tracking. The sorting
rules push these to the top periodically to make sure I don't forget about them
(See the sorting exceptions below)

`"HOLD"` tasks are next. They are a notable exception to the normal status format
since tasks in `"HOLD"` status have a date in the future. Once that date arrives,
`"HOLD"` tasks are moved to the top of the list. This allows me to put off a task
until a specific date.

`"DONE"` tasks are last. They should include the final status for the task, and
are put at the bottom to keep a record of how and when a task was completed.
They are cleared during [Generation](#generation)


##### Sorting Exceptions
Some useful exceptions to the ordering are:
* Tasks marked `"HOLD"` with today's date or a date in the past are sent to the
  top. This allows me to mark a task as HOLD, at which point it will be sent to
  the bottom of the list, until a particular date when it will rise to the top of
  the list to get my attention.
* Tasks marked `"WAITING"`, `"REVIEW"` or `"RESPONDED"` with yesterday's date or
  older will be sent to the top. This allows me to mark a task with any of these
  statuses and have it sent to the bottom of the list, and pop back to the top
  the next day so I can check them again. I use this for PR's which need to be
  reviewed over-night and other tasks I need to wait until the following day to
  check up on.
* Tasks marked `"STALE"` with a date 7 days or more in the past will be sent to
  the top. I use this to periodically check up on slow-moving or stale tasks.

### Status
Each item in `Morning Start Up` and `TODO` has a `Status`, which appears as the
last element of the line. In its most basic form, a `Status` is just a status
name inside square brackets. A status name is a series of capitalized letters.
A basic status might be `[IN PROGRESS]`.

Each status also has a comment field, and a date. Comments are optional, but
`today` will apply a date to any `Status` that doesn't already have one.

Comments come after the status name, separated by space and a hyphen, and can
be any string of text: `[IN PROGRESS - Working on pr #12]`.

Dates follow the status name and the comment if one is present, again separated
by space and hyphen. Dates must follow the format (given in Go's
[`time.Format`](https://golang.org/pkg/time/#Time.Format) notation) `Jan _2,
2006` or they will be considered part of the comment. Both of these are valid
statuses with dates:
```
[IN PROGRESS - Working on pr #12 - Jan 16, 2020]
[READY - Jan 14, 2020]
```

## Use 
Today operates in the directory `~/today`, or the directory specified with the
`-d` option.

When `today` is run, it looks for a file named with the current day's date in
the operating directory. If it does not find one, it attempts to generate one
from the most recent previous file, determined by the date in the file name.
(See [Generation](#generation))


If doing [Generation](#generation), or when passed the `-c` flag, `today` will
clear `"DONE"` tasks and `Morning Start Up` statuses.

By default, `today` will update the statuses of the TODO section, and then sort
the tasks.

With the `-i` flag, `today` will read from stdin and write to stdout rather
than looking in any directory.



### Generation
Generation is simply the process of using a previous day's today file to
generate a today file for the current day. With no flags, `today` will first
look for a today file for the current day. If one is not present, it will
attempt to generate one from the most recent previous today file. If none is
present, it will generate an empty today file for the current day.

Currently, when generating a today file from an existing previous today file,
`today` clears out tasks with `"DONE"` status, clears the `Log` section, and
removes statuses from the `Morning Start Up` section.
