package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"sort"
	"time"

	"github.com/knusbaum/today"
)

const (
	noteFormat = "note.2006.Jan.02.txt"
)

var errNoTodayFiles error = fmt.Errorf("no existing today files")

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

func todayExists(dir string) (bool, error) {
	name := path.Join(dir, time.Now().Format(noteFormat))
	_, err := os.Stat(name)
	if err != nil && !os.IsNotExist(err) {
		return false, err
	} else if os.IsNotExist(err) {
		return false, nil
	}
	return true, nil
}

type fileDate struct {
	name string
	date time.Time
}

type byDate []fileDate

func (a byDate) Len() int           { return len(a) }
func (a byDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byDate) Less(i, j int) bool { return a[i].date.After(a[j].date) }

func openMostRecent(dir string) (*os.File, error) {
	d, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	files, err := d.Readdirnames(0)
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, errNoTodayFiles
	}

	filedates := make([]fileDate, 0, len(files))
	for i := range files {
		date, err := time.Parse(noteFormat, files[i])
		if err == nil {
			filedates = append(filedates, fileDate{files[i], date})
		}
	}

	sort.Sort(byDate(filedates))

	f, err := os.Open(path.Join(dir, filedates[0].name))
	if err != nil {
		return nil, err
	}
	return f, nil
}

func openReadToday(dir string) (*os.File, error) {
	name := path.Join(dir, time.Now().Format(noteFormat))
	return os.Open(name)
}

func openWriteToday(dir string) (*os.File, error) {
	name := path.Join(dir, time.Now().Format(noteFormat))
	copyFileContents(name, path.Join(dir, ".backup"))
	return os.Create(name)
}

func generateToday(dir string) error {
	f, err := openMostRecent(dir)
	if err != nil {
		if err == errNoTodayFiles {
			out, err := openWriteToday(dir)
			if err != nil {
				return err
			}
			defer out.Close()
			var t today.Today
			return t.Write(out)
		}
		return err
	}
	defer f.Close()

	t, err := today.NewParser(f).Parse()
	if err != nil {
		log.Fatalf("Failed to parse today: %s", err)
	}
	t.Update()
	t.Sort()
	t.Clear()

	out, err := openWriteToday(dir)
	if err != nil {
		return err
	}
	defer out.Close()
	return t.Write(out)
}

func main() {
	todaydir, _ := os.UserHomeDir()
	if todaydir != "" {
		todaydir += "/today"
	}

	var dir = flag.String("d", todaydir, "The directory in which the today logs reside.")
	var pipe = flag.Bool("i", false, "Read from stdin and write to stdout rather than files in the directory specified with -d.")
	var sort = flag.Bool("s", true, "Sort the todo entries according to priority.")
	var update = flag.Bool("u", true, "Update the dates for the todo entries.")
	var clear = flag.Bool("c", false, "Clear the DONE tasks. By default, this only happens when generating the today file.")

	flag.Parse()

	var (
		in  io.Reader
		out io.Writer
	)

	if *pipe {
		in = os.Stdin
	} else {
		exists, err := todayExists(*dir)
		if err != nil {
			log.Fatalf("Failed to read todayfile: %s", err)
		}
		if exists {
			f, err := openReadToday(*dir)
			if err != nil {
				log.Fatalf("Failed to read todayfile: %s", err)
			}
			defer f.Close()
			in = f
		} else {
			err = generateToday(*dir)
			if err != nil {
				log.Fatalf("Failed to generate todayfile: %s", err)
			}
			return
		}
	}

	t, err := today.NewParser(in).Parse()
	if err != nil {
		log.Fatalf("Failed to parse today: %s", err)
	}

	if *update {
		t.Update()
	}
	if *sort {
		t.Sort()
	}
	if *clear {
		t.Clear()
	}

	if *pipe {
		out = os.Stdout
	} else {
		f, err := openWriteToday(*dir)
		if err != nil {
			log.Fatalf("Failed to write todayfile: %s", err)
		}
		defer f.Close()
		out = f
	}
	err = t.Write(out)
	if err != nil {
		log.Fatalf("Failed to write today: %s", err)
	}
}
