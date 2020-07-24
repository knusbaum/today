//go:generate sh -c "godoc2md -template readme.tmpl github.com/knusbaum/today > README.md"

// Package today provides tools for manipulating a plain-text task list. This was developed to help
// me specifically with my workflow. The package manages "today files" which contain several
// sections. Each section has slightly different structure and different rules about how it is
// updated, sorted, etc. The today package is primarily intended to be used by the today program in
// package github.com/knusbaum/today/today.
package today
